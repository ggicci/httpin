package codec

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/ggicci/httpin/internal"
	"github.com/ggicci/httpin/internal/testutil"
	"github.com/ggicci/httpin/patch"
	"github.com/ggicci/strconvx"
	"github.com/stretchr/testify/assert"
)

var testNS = NewNamespace()

type Point2D struct {
	X int
	Y int
}

func (p Point2D) ToString() (string, error) {
	return fmt.Sprintf("Point2D(%d,%d)", p.X, p.Y), nil
}

func (p *Point2D) FromString(s string) error {
	_, err := fmt.Sscanf(s, "Point2D(%d,%d)", &p.X, &p.Y)
	return err
}

type MyStruct struct {
	Name        string
	NamePointer *string
	Names       []string
	PatchName   patch.Field[string]
	PatchNames  patch.Field[[]string]

	Age        int
	AgePointer *int
	Ages       []int
	PatchAge   patch.Field[int]
	PatchAges  patch.Field[[]int]

	Dot        Point2D
	DotPointer *Point2D
	Dots       []Point2D
	PatchDot   patch.Field[Point2D]
	PatchDots  patch.Field[[]Point2D]
}

func TestStringCodec_FromString(t *testing.T) {
	rv := reflect.New(reflect.TypeOf(MyStruct{})).Elem()
	s := rv.Addr().Interface().(*MyStruct)

	// string
	testAssignString(t, rv.FieldByName("Name"), "Alice")
	testAssignString(t, rv.FieldByName("NamePointer"), "Charlie")
	testNewStringCodecErrUnsupported(t, rv.FieldByName("Names"))
	testAssignString(t, rv.FieldByName("PatchName"), "Bob")
	testNewStringCodecErrUnsupported(t, rv.FieldByName("PatchNames"))

	assert.Equal(t, "Alice", s.Name)
	assert.Equal(t, "Charlie", *s.NamePointer)
	assert.Equal(t, []string(nil), s.Names)
	assert.Equal(t, "Bob", s.PatchName.Value)
	assert.True(t, s.PatchName.Valid)

	// int
	testAssignString(t, rv.FieldByName("Age"), "18")
	testAssignString(t, rv.FieldByName("AgePointer"), "20")
	testNewStringCodecErrUnsupported(t, rv.FieldByName("Ages"))
	testAssignString(t, rv.FieldByName("PatchAge"), "18")
	testNewStringCodecErrUnsupported(t, rv.FieldByName("PatchAges"))

	assert.Equal(t, 18, s.Age)
	assert.Equal(t, 20, *s.AgePointer)
	assert.Equal(t, []int(nil), s.Ages)
	assert.Equal(t, 18, s.PatchAge.Value)
	assert.True(t, s.PatchAge.Valid)

	// Point2D
	testAssignString(t, rv.FieldByName("Dot"), "Point2D(1,2)")
	testAssignString(t, rv.FieldByName("DotPointer"), "Point2D(3,4)")
	testNewStringCodecErrUnsupported(t, rv.FieldByName("Dots"))
	testAssignString(t, rv.FieldByName("PatchDot"), "Point2D(5,6)")
	testNewStringCodecErrUnsupported(t, rv.FieldByName("PatchDots"))

	assert.Equal(t, Point2D{1, 2}, s.Dot)
	assert.Equal(t, &Point2D{3, 4}, s.DotPointer)
	assert.Equal(t, []Point2D(nil), s.Dots)
	assert.Equal(t, Point2D{5, 6}, s.PatchDot.Value)
	assert.True(t, s.PatchDot.Valid)
}

func TestStringCodec_ToString(t *testing.T) {
	var s = &MyStruct{
		Name:        "Alice",
		NamePointer: internal.Pointerize[string]("Charlie"),
		Names:       []string{"Alice", "Bob", "Charlie"},
		PatchName:   patch.Field[string]{Value: "Bob", Valid: true},
		PatchNames:  patch.Field[[]string]{Value: []string{"Alice", "Bob", "Charlie"}, Valid: true},

		Age:        18,
		AgePointer: internal.Pointerize[int](20),
		Ages:       []int{18, 20},
		PatchAge:   patch.Field[int]{Value: 18, Valid: true},
		PatchAges:  patch.Field[[]int]{Value: []int{18, 20}, Valid: true},

		Dot:        Point2D{1, 2},
		DotPointer: internal.Pointerize[Point2D](Point2D{3, 4}),
		Dots:       []Point2D{{1, 2}, {3, 4}},
		PatchDot:   patch.Field[Point2D]{Value: Point2D{5, 6}, Valid: true},
		PatchDots:  patch.Field[[]Point2D]{Value: []Point2D{{1, 2}, {3, 4}}, Valid: true},
	}

	rv := reflect.ValueOf(s).Elem()

	assert.Equal(t, "Alice", testGetString(t, rv.FieldByName("Name")))
	assert.Equal(t, "Charlie", testGetString(t, rv.FieldByName("NamePointer")))
	testNewStringCodecErrUnsupported(t, rv.FieldByName("Names"))
	assert.Equal(t, "Bob", testGetString(t, rv.FieldByName("PatchName")))
	testNewStringCodecErrUnsupported(t, rv.FieldByName("PatchNames"))

	assert.Equal(t, "18", testGetString(t, rv.FieldByName("Age")))
	assert.Equal(t, "20", testGetString(t, rv.FieldByName("AgePointer")))
	testNewStringCodecErrUnsupported(t, rv.FieldByName("Ages"))
	assert.Equal(t, "18", testGetString(t, rv.FieldByName("PatchAge")))
	testNewStringCodecErrUnsupported(t, rv.FieldByName("PatchAges"))

	assert.Equal(t, "Point2D(1,2)", testGetString(t, rv.FieldByName("Dot")))
	assert.Equal(t, "Point2D(3,4)", testGetString(t, rv.FieldByName("DotPointer")))
	testNewStringCodecErrUnsupported(t, rv.FieldByName("Dots"))
	assert.Equal(t, "Point2D(5,6)", testGetString(t, rv.FieldByName("PatchDot")))
	testNewStringCodecErrUnsupported(t, rv.FieldByName("PatchDots"))
}

func TestStringCodec4PatchField_ToString(t *testing.T) {
	var patchString = patch.Field[string]{Value: "Alice", Valid: true}
	rv := reflect.ValueOf(&patchString).Elem()
	assert.True(t, IsPatchField(rv.Type()))
	StringCodec, err := testNS.NewStringCodec4PatchField(rv, nil)
	assert.NoError(t, err)

	sv, err := StringCodec.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "Alice", sv)

	patchString.Valid = false // make it invalid
	sv, err = StringCodec.ToString()
	assert.ErrorContains(t, err, "invalid value")
	assert.Empty(t, sv, "invalid patch field should return empty string")
}

func TestStringCodec4PatchField_FromString(t *testing.T) {
	// string
	var patchString = patch.Field[string]{}

	assert.Empty(t, patchString.Value)
	assert.False(t, patchString.Valid)

	rv := reflect.ValueOf(&patchString).Elem()
	assert.True(t, IsPatchField(rv.Type()))
	StringCodec, err := testNS.NewStringCodec4PatchField(rv, nil)
	assert.NoError(t, err)
	assert.NoError(t, StringCodec.FromString("Alice"))
	assert.Equal(t, "Alice", patchString.Value)
	assert.True(t, patchString.Valid, "Valid should be set to true after a succssful FromString")

	// int
	var patchInt = patch.Field[int]{}
	rv = reflect.ValueOf(&patchInt).Elem()
	assert.True(t, IsPatchField(rv.Type()))
	StringCodec, err = testNS.NewStringCodec4PatchField(rv, nil)
	assert.NoError(t, err)
	assert.Error(t, StringCodec.FromString("Alice")) // cannot convert "Alice" to int
	assert.Zero(t, patchInt.Value, "Value should not be changed after a failed FromString")
	assert.False(t, patchInt.Valid, "Valid should not be changed after a failed FromString")

	assert.NoError(t, StringCodec.FromString("18"))
	assert.Equal(t, 18, patchInt.Value)
	assert.True(t, patchInt.Valid, "Valid should be set to true after a succssful FromString")

	assert.Error(t, StringCodec.FromString("18.0")) // cannot convert "18.0" to int
	assert.Equal(t, 18, patchInt.Value, "Value should not be changed after a failed FromString")
	assert.True(t, patchInt.Valid, "Valid should not be changed after a failed FromString")
}

func TestStringCodec_WithAdaptor(t *testing.T) {
	adapt := func(t *time.Time) (StringCodec, error) { return (*testutil.MyDate)(t), nil }
	var now = time.Now()
	rvTimePointer := reflect.ValueOf(&now)

	_, adaptor := strconvx.ToAnyAdaptor(adapt)
	coder, err := testNS.NewStringCodec(rvTimePointer, adaptor)
	assert.NoError(t, err)
	assert.NoError(t, coder.FromString("1991-11-10"))

	s, err := coder.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "1991-11-10", s)

	assert.ErrorContains(t, coder.FromString("1991-11-10T08:00:00+08:00"), "parsing time")
}

func testAssignString(t *testing.T, rv reflect.Value, value string) {
	s, err := testNS.NewStringCodec(rv, nil)
	assert.NoError(t, err)
	assert.NoError(t, s.FromString(value))
}

func testNewStringCodecErrUnsupported(t *testing.T, rv reflect.Value) {
	s, err := testNS.NewStringCodec(rv, nil)
	assert.ErrorIs(t, err, ErrUnsupportedType)
	assert.Nil(t, s)
}

func testGetString(t *testing.T, rv reflect.Value) string {
	coder, err := testNS.NewStringCodec(rv, nil)
	assert.NoError(t, err)
	s, err := coder.ToString()
	assert.NoError(t, err)
	return s
}
