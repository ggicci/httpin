package patch_test

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/ggicci/httpin/patch"
	"github.com/stretchr/testify/assert"
)

func shouldBeNil(t *testing.T, err error, failMessage string) {
	if err != nil {
		t.Logf("%s, got error: %v", failMessage, err)
		t.Fail()
	}
}

func shouldResemble(t *testing.T, va, vb any, failMessage string) {
	if reflect.DeepEqual(va, vb) {
		return
	}
	t.Logf("%s, expected %#v, got %#v", failMessage, va, vb)
	t.Fail()
}

func fixedZone(offset int) *time.Location {
	if offset == 0 {
		return time.UTC
	}
	_, localOffset := time.Now().Local().Zone()
	if offset == localOffset {
		return time.Local
	}
	return time.FixedZone("", offset)
}

func testJSONMarshalling(t *testing.T, tc testcase) {
	bs, err := json.Marshal(tc.Expected)
	if err != nil {
		t.Logf("marshal failed, got error: %v", err)
		t.Fail()
	}
	if string(bs) != tc.Content {
		t.Logf("marshal failed, expected %q, got %q", tc.Content, string(bs))
		t.Fail()
	}
}

func testJSONUnmarshalling(t *testing.T, tc testcase) {
	rt := reflect.TypeOf(tc.Expected) // type: patch.Field
	rv := reflect.New(rt)             // rv: *patch.Field

	shouldBeNil(t, json.Unmarshal([]byte(tc.Content), rv.Interface()), "unmarshal failed")
	shouldResemble(t, rv.Elem().Interface(), tc.Expected, "unmarshal failed")
}

type testcase struct {
	Content  string
	Expected any
}

type GitHubProfile struct {
	Id        int64  `json:"id"`
	Login     string `json:"login"`
	AvatarUrl string `json:"avatar_url"`
}

type GenderType string

type Account struct {
	Id     int64
	Email  string
	Tags   []string
	Gender GenderType
	GitHub *GitHubProfile
}

type AccountPatch struct {
	Email  patch.Field[string]         `json:"email"`
	Tags   patch.Field[[]string]       `json:"tags"`
	Gender patch.Field[GenderType]     `json:"gender"`
	GitHub patch.Field[*GitHubProfile] `json:"github"`
}

func TestField(t *testing.T) {
	var cases = []testcase{
		{"true", patch.Field[bool]{true, true}},
		{"false", patch.Field[bool]{false, true}},
		{"2045", patch.Field[int]{2045, true}},
		{"127", patch.Field[int8]{127, true}},
		{"32767", patch.Field[int16]{32767, true}},
		{"2147483647", patch.Field[int32]{2147483647, true}},
		{"9223372036854775807", patch.Field[int64]{9223372036854775807, true}},
		{"2045", patch.Field[uint]{2045, true}},
		{"255", patch.Field[uint8]{255, true}},
		{"65535", patch.Field[uint16]{65535, true}},
		{"4294967295", patch.Field[uint32]{4294967295, true}},
		{"18446744073709551615", patch.Field[uint64]{18446744073709551615, true}},
		{"3.14", patch.Field[float32]{3.14, true}},
		{"3.14", patch.Field[float64]{3.14, true}},
		{"\"hello\"", patch.Field[string]{"hello", true}},

		// Array
		{`[true,false]`, patch.Field[[]bool]{[]bool{true, false}, true}},
		{"[1,2,3]", patch.Field[[]int]{[]int{1, 2, 3}, true}},
		{"[1,2,3]", patch.Field[[]int8]{[]int8{1, 2, 3}, true}},
		{"[1,2,3]", patch.Field[[]int16]{[]int16{1, 2, 3}, true}},
		{"[1,2,3]", patch.Field[[]int32]{[]int32{1, 2, 3}, true}},
		{"[1,2,3]", patch.Field[[]int64]{[]int64{1, 2, 3}, true}},
		{"[1,2,3]", patch.Field[[]uint]{[]uint{1, 2, 3}, true}},
		// NOTE(ggicci): []uint8 is a special case, check TestFieldUint8Array
		{"[1,2,3]", patch.Field[[]uint16]{[]uint16{1, 2, 3}, true}},
		{"[1,2,3]", patch.Field[[]uint32]{[]uint32{1, 2, 3}, true}},
		{"[1,2,3]", patch.Field[[]uint64]{[]uint64{1, 2, 3}, true}},
		{"[0.618,1,3.14]", patch.Field[[]float32]{[]float32{0.618, 1, 3.14}, true}},
		{"[0.618,1,3.14]", patch.Field[[]float64]{[]float64{0.618, 1, 3.14}, true}},
		{`["hello","world"]`, patch.Field[[]string]{[]string{"hello", "world"}, true}},

		// time.Time
		{
			`"2019-08-25T07:19:34Z"`,
			patch.Field[time.Time]{
				time.Date(2019, 8, 25, 7, 19, 34, 0, fixedZone(0)),
				true,
			},
		},
		{
			`"1991-11-10T08:00:00-07:00"`,
			patch.Field[time.Time]{
				time.Date(1991, 11, 10, 8, 0, 0, 0, fixedZone(-7*3600)),
				true,
			},
		},
		{
			`"1991-11-10T08:00:00+08:00"`,
			patch.Field[time.Time]{
				time.Date(1991, 11, 10, 8, 0, 0, 0, fixedZone(+8*3600)),
				true,
			},
		},

		// Custom structs
		{
			`{"Id":1000,"Email":"ggicci@example.com","Tags":["developer","修勾"],"Gender":"male","GitHub":{"id":3077555,"login":"ggicci","avatar_url":"https://avatars.githubusercontent.com/u/3077555?v=4"}}`,
			patch.Field[*Account]{
				&Account{
					Id:     1000,
					Email:  "ggicci@example.com",
					Tags:   []string{"developer", "修勾"},
					Gender: "male",
					GitHub: &GitHubProfile{
						Id:        3077555,
						Login:     "ggicci",
						AvatarUrl: "https://avatars.githubusercontent.com/u/3077555?v=4",
					},
				},
				true,
			},
		},
	}

	for _, c := range cases {
		testJSONMarshalling(t, c)
		testJSONUnmarshalling(t, c)
	}
}

// TestFieldUint8Array runs JSON marshalling & unmarshalling tests on type Field[[]uint8].
// Because in golang's encoding/json package, encoding uint8[] is special.
// See: https://golang.org/pkg/encoding/json/#Marshal
//
// > Array and slice values encode as JSON arrays, except that []byte encodes
// as a base64-encoded string, and a nil slice encodes as the null JSON
// value.
//
//	uint8       the set of all unsigned  8-bit integers (0 to 255)
//	byte        alias for uint8
func TestFieldUint8Array(t *testing.T) {
	var a1 patch.Field[[]uint8]
	// unmarshal
	shouldBeNil(t, json.Unmarshal([]byte("[1,2,3]"), &a1), "unmarshal Field[[]uint8] failed")
	shouldResemble(t, patch.Field[[]uint8]{[]uint8{1, 2, 3}, true}, a1, "unmarshal Field[[]uint8] failed")

	// marshal
	var a2 = patch.Field[[]uint8]{[]uint8{1, 2, 3}, true}
	out, err := json.Marshal(a2)
	shouldBeNil(t, err, "marshal Field[[]uint8] failed")
	shouldResemble(t, `"AQID"`, string(out), "marshal Field[[]uint8] failed")
}

func TestField_UnmarshalJSON_Struct(t *testing.T) {
	var testcases = []testcase{
		{
			`{"email":"ggicci.2@example.com","tags":["artist","photographer"]}`,
			AccountPatch{
				Email:  patch.Field[string]{"ggicci.2@example.com", true},
				Gender: patch.Field[GenderType]{"", false},
				Tags:   patch.Field[[]string]{[]string{"artist", "photographer"}, true},
				GitHub: patch.Field[*GitHubProfile]{nil, false},
			},
		},
		{
			`{"tags":null,"gender":"female","github":{"id":100,"login":"ggicci.2","avatar_url":null}}`,
			AccountPatch{
				Email:  patch.Field[string]{"", false},
				Gender: patch.Field[GenderType]{"female", true},
				Tags:   patch.Field[[]string]{nil, false},
				GitHub: patch.Field[*GitHubProfile]{&GitHubProfile{
					Id:        100,
					Login:     "ggicci.2",
					AvatarUrl: "",
				}, true},
			},
		},
	}

	for _, c := range testcases {
		testJSONUnmarshalling(t, c)
	}
}

func TestField_MarshalJSON_Struct(t *testing.T) {
	var testcases = []testcase{
		{
			`{"email":"hello","tags":null,"gender":null,"github":null}`,
			AccountPatch{
				Email:  patch.Field[string]{"hello", true},
				Tags:   patch.Field[[]string]{nil, false},
				Gender: patch.Field[GenderType]{"", false},
				GitHub: patch.Field[*GitHubProfile]{nil, false},
			},
		},
	}

	for _, c := range testcases {
		testJSONMarshalling(t, c)
	}
}

func TestField_ValidSentinel(t *testing.T) {
	f := patch.Field[string]{"hello", true}
	assert.True(t, f.IsValid())
	f.SetValid(false)
	assert.False(t, f.IsValid())
	f.SetValid(true)
	assert.True(t, f.IsValid())
}
