package internal

import (
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

var (
	ErrUnsupportedType = errors.New("unsupported type")
	ErrTypeMismatch    = errors.New("type mismatch")

	builtinStringableAdaptors = make(map[reflect.Type]AnyStringableAdaptor)
)

func init() {
	builtinStringable[string](func(v *string) (Stringable, error) { return (*String)(v), nil })
	builtinStringable[bool](func(v *bool) (Stringable, error) { return (*Bool)(v), nil })
	builtinStringable[int](func(v *int) (Stringable, error) { return (*Int)(v), nil })
	builtinStringable[int8](func(v *int8) (Stringable, error) { return (*Int8)(v), nil })
	builtinStringable[int16](func(v *int16) (Stringable, error) { return (*Int16)(v), nil })
	builtinStringable[int32](func(v *int32) (Stringable, error) { return (*Int32)(v), nil })
	builtinStringable[int64](func(v *int64) (Stringable, error) { return (*Int64)(v), nil })
	builtinStringable[uint](func(v *uint) (Stringable, error) { return (*Uint)(v), nil })
	builtinStringable[uint8](func(v *uint8) (Stringable, error) { return (*Uint8)(v), nil })
	builtinStringable[uint16](func(v *uint16) (Stringable, error) { return (*Uint16)(v), nil })
	builtinStringable[uint32](func(v *uint32) (Stringable, error) { return (*Uint32)(v), nil })
	builtinStringable[uint64](func(v *uint64) (Stringable, error) { return (*Uint64)(v), nil })
	builtinStringable[float32](func(v *float32) (Stringable, error) { return (*Float32)(v), nil })
	builtinStringable[float64](func(v *float64) (Stringable, error) { return (*Float64)(v), nil })
	builtinStringable[complex64](func(v *complex64) (Stringable, error) { return (*Complex64)(v), nil })
	builtinStringable[complex128](func(v *complex128) (Stringable, error) { return (*Complex128)(v), nil })
	builtinStringable[time.Time](func(v *time.Time) (Stringable, error) { return (*Time)(v), nil })
	builtinStringable[[]byte](func(b *[]byte) (Stringable, error) { return (*ByteSlice)(b), nil })
}

type StringMarshaler interface {
	ToString() (string, error)
}

type StringUnmarshaler interface {
	FromString(string) error
}

type Stringable interface {
	StringMarshaler
	StringUnmarshaler
}

// NewStringable returns a Stringable from the given reflect.Value.
// We assume that the given reflect.Value is a non-nil pointer to a value.
// It will panic if the given reflect.Value is not a pointer.
func NewStringable(rv reflect.Value) (Stringable, error) {
	baseType := rv.Type().Elem()
	if adapt, ok := builtinStringableAdaptors[baseType]; ok {
		return adapt(rv.Interface())
	} else {
		return nil, UnsupportedType(baseType)
	}
}

func builtinStringable[T any](builder StringableAdaptor[T]) {
	builtinStringableAdaptors[TypeOf[T]()] = ToAnyStringableAdaptor[T](builder)
}

type String string

func (sv String) ToString() (string, error) {
	return string(sv), nil
}

func (sv *String) FromString(s string) error {
	*sv = String(s)
	return nil
}

type Bool bool

func (bv Bool) ToString() (string, error) {
	return strconv.FormatBool(bool(bv)), nil
}

func (bv *Bool) FromString(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	*bv = Bool(v)
	return nil
}

type Int int

func (iv Int) ToString() (string, error) {
	return strconv.Itoa(int(iv)), nil
}

func (iv *Int) FromString(s string) error {
	v, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*iv = Int(v)
	return nil
}

type Int8 int8

func (iv Int8) ToString() (string, error) {
	return strconv.FormatInt(int64(iv), 10), nil
}

func (iv *Int8) FromString(s string) error {
	v, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		return err
	}
	*iv = Int8(v)
	return nil
}

type Int16 int16

func (iv Int16) ToString() (string, error) {
	return strconv.FormatInt(int64(iv), 10), nil
}

func (iv *Int16) FromString(s string) error {
	v, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		return err
	}
	*iv = Int16(v)
	return nil
}

type Int32 int32

func (iv Int32) ToString() (string, error) {
	return strconv.FormatInt(int64(iv), 10), nil
}

func (iv *Int32) FromString(s string) error {
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return err
	}
	*iv = Int32(v)
	return nil
}

type Int64 int64

func (iv Int64) ToString() (string, error) {
	return strconv.FormatInt(int64(iv), 10), nil
}

func (iv *Int64) FromString(s string) error {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	*iv = Int64(v)
	return nil
}

type Uint uint

func (uv Uint) ToString() (string, error) {
	return strconv.FormatUint(uint64(uv), 10), nil
}

func (uv *Uint) FromString(s string) error {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}
	*uv = Uint(v)
	return nil
}

type Uint8 uint8

func (uv Uint8) ToString() (string, error) {
	return strconv.FormatUint(uint64(uv), 10), nil
}

func (uv *Uint8) FromString(s string) error {
	v, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return err
	}
	*uv = Uint8(v)
	return nil
}

type Uint16 uint16

func (uv Uint16) ToString() (string, error) {
	return strconv.FormatUint(uint64(uv), 10), nil
}

func (uv *Uint16) FromString(s string) error {
	v, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return err
	}
	*uv = Uint16(v)
	return nil
}

type Uint32 uint32

func (uv Uint32) ToString() (string, error) {
	return strconv.FormatUint(uint64(uv), 10), nil
}

func (uv *Uint32) FromString(s string) error {
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return err
	}
	*uv = Uint32(v)
	return nil
}

type Uint64 uint64

func (uv Uint64) ToString() (string, error) {
	return strconv.FormatUint(uint64(uv), 10), nil
}

func (uv *Uint64) FromString(s string) error {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}
	*uv = Uint64(v)
	return nil
}

type Float32 float32

func (fv Float32) ToString() (string, error) {
	return strconv.FormatFloat(float64(fv), 'f', -1, 32), nil
}

func (fv *Float32) FromString(s string) error {
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return err
	}
	*fv = Float32(v)
	return nil
}

type Float64 float64

func (fv Float64) ToString() (string, error) {
	return strconv.FormatFloat(float64(fv), 'f', -1, 64), nil
}

func (fv *Float64) FromString(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*fv = Float64(v)
	return nil
}

type Complex64 complex64

func (cv Complex64) ToString() (string, error) {
	return strconv.FormatComplex(complex128(cv), 'f', -1, 64), nil
}

func (cv *Complex64) FromString(s string) error {
	v, err := strconv.ParseComplex(s, 64)
	if err != nil {
		return err
	}
	*cv = Complex64(v)
	return nil
}

type Complex128 complex128

func (cv Complex128) ToString() (string, error) {
	return strconv.FormatComplex(complex128(cv), 'f', -1, 128), nil
}

func (cv *Complex128) FromString(s string) error {
	v, err := strconv.ParseComplex(s, 128)
	if err != nil {
		return err
	}
	*cv = Complex128(v)
	return nil
}

type Time time.Time

func (tv Time) ToString() (string, error) {
	return time.Time(tv).UTC().Format(time.RFC3339Nano), nil
}

func (tv *Time) FromString(s string) error {
	if t, err := DecodeTime(s); err != nil {
		return err
	} else {
		*tv = Time(t)
		return nil
	}
}

func UnsupportedType(rt reflect.Type) error {
	return fmt.Errorf("%w: %v", ErrUnsupportedType, rt)
}

// ByteSlice is a wrapper of []byte to implement Stringable.
// NOTE: we're using base64.StdEncoding here, not base64.URLEncoding.
type ByteSlice []byte

func (bs ByteSlice) ToString() (string, error) {
	return base64.StdEncoding.EncodeToString(bs), nil
}

func (bs *ByteSlice) FromString(s string) error {
	v, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	*bs = ByteSlice(v)
	return nil
}
