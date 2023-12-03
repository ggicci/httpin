package internal

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	var s string
	sv := (*String)(&s)
	sv.FromString("Alice")
	assert.Equal(t, "Alice", s)
}

func TestInt(t *testing.T) {
	var i int
	iv := (*Int)(&i)
	iv.FromString("18")
	assert.Equal(t, 18, i)
}

func TestNewStringable_string(t *testing.T) {
	var s string = "hello"
	rvString := reflect.ValueOf(s)
	assert.Panics(t, func() {
		NewStringable(rvString)
	})

	rvStringPointer := reflect.ValueOf(&s)
	sv, err := NewStringable(rvStringPointer)
	assert.NoError(t, err)
	got, err := sv.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "hello", got)
	sv.FromString("world")
	assert.Equal(t, "world", s)
}

func TestNewStringable_Time(t *testing.T) {
	var now = time.Now()
	rvTime := reflect.ValueOf(now)
	assert.Panics(t, func() {
		NewStringable(rvTime)
	})

	rvTimePointer := reflect.ValueOf(&now)
	sv, err := NewStringable(rvTimePointer)
	assert.NoError(t, err)
	assert.NoError(t, sv.FromString("1991-11-10T08:00:00+08:00"))
	assert.Equal(t, "1991-11-10T00:00:00Z", now.Format(time.RFC3339))
}
