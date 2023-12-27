package core

import (
	"reflect"
	"testing"

	"github.com/ggicci/httpin/internal"
	"github.com/stretchr/testify/assert"
)

func TestStringSlicable_FromStringSlice(t *testing.T) {
	rv := reflect.New(reflect.TypeOf(MyStruct{})).Elem()
	s := rv.Addr().Interface().(*MyStruct)

	testAssignStringSlice(t, rv.FieldByName("Name"), []string{"Alice"})
	testAssignStringSlice(t, rv.FieldByName("NamePointer"), []string{"Charlie"})
	testAssignStringSlice(t, rv.FieldByName("Names"), []string{"Alice", "Bob", "Charlie"})

	testAssignStringSlice(t, rv.FieldByName("Age"), []string{"18"})
	testAssignStringSlice(t, rv.FieldByName("AgePointer"), []string{"20"})
	testAssignStringSlice(t, rv.FieldByName("Ages"), []string{"18", "20"})

	assert.Equal(t, "Alice", s.Name)
	assert.Equal(t, "Charlie", *s.NamePointer)
	assert.Equal(t, []string{"Alice", "Bob", "Charlie"}, s.Names)

	assert.Equal(t, 18, s.Age)
	assert.Equal(t, 20, *s.AgePointer)
	assert.Equal(t, []int{18, 20}, s.Ages)
}

func TestStringSlicable_ToStringSlice(t *testing.T) {
	var s = &MyStruct{
		Name:        "Alice",
		NamePointer: internal.Pointerize[string]("Charlie"),
		Names:       []string{"Alice", "Bob", "Charlie"},

		Age:        18,
		AgePointer: internal.Pointerize[int](20),
		Ages:       []int{18, 20},
	}

	rv := reflect.ValueOf(s).Elem()
	assert.Equal(t, []string{"Alice"}, testGetStringSlice(t, rv.FieldByName("Name")))
	assert.Equal(t, []string{"Charlie"}, testGetStringSlice(t, rv.FieldByName("NamePointer")))
	assert.Equal(t, []string{"Alice", "Bob", "Charlie"}, testGetStringSlice(t, rv.FieldByName("Names")))

	assert.Equal(t, []string{"18"}, testGetStringSlice(t, rv.FieldByName("Age")))
	assert.Equal(t, []string{"20"}, testGetStringSlice(t, rv.FieldByName("AgePointer")))
	assert.Equal(t, []string{"18", "20"}, testGetStringSlice(t, rv.FieldByName("Ages")))
}

func testAssignStringSlice(t *testing.T, rv reflect.Value, values []string) {
	ss, err := NewStringSlicable(rv, nil)
	assert.NoError(t, err)
	assert.NoError(t, ss.FromStringSlice(values))
}

func testGetStringSlice(t *testing.T, rv reflect.Value) []string {
	ss, err := NewStringSlicable(rv, nil)
	assert.NoError(t, err)
	values, err := ss.ToStringSlice()
	assert.NoError(t, err)
	return values
}
