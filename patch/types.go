package patch

import (
	"encoding/json"
	"time"
)

type Bool struct {
	Value bool
	Valid bool // Valid is true if the corresponding key were found in the source
}

func (t Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Bool) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Int struct {
	Value int
	Valid bool
}

func (t Int) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Int) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Int8 struct {
	Value int8
	Valid bool
}

func (t Int8) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Int8) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Int16 struct {
	Value int16
	Valid bool
}

func (t Int16) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Int16) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Int32 struct {
	Value int32
	Valid bool
}

func (t Int32) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Int32) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Int64 struct {
	Value int64
	Valid bool
}

func (t Int64) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Int64) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Uint struct {
	Value uint
	Valid bool
}

func (t Uint) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Uint) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Uint8 struct {
	Value uint8
	Valid bool
}

func (t Uint8) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Uint8) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Uint16 struct {
	Value uint16
	Valid bool
}

func (t Uint16) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Uint16) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Uint32 struct {
	Value uint32
	Valid bool
}

func (t Uint32) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Uint32) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Uint64 struct {
	Value uint64
	Valid bool
}

func (t Uint64) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Uint64) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Float32 struct {
	Value float32
	Valid bool
}

func (t Float32) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Float32) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Float64 struct {
	Value float64
	Valid bool
}

func (t Float64) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Float64) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Complex64 struct {
	Value complex64
	Valid bool
}

func (t Complex64) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Complex64) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Complex128 struct {
	Value complex128
	Valid bool
}

func (t Complex128) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Complex128) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type String struct {
	Value string
	Valid bool
}

func (t String) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *String) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Time struct {
	Value time.Time
	Valid bool
}

func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Time) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type BoolArray struct {
	Value []bool
	Valid bool
}

func (t BoolArray) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *BoolArray) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type IntArray struct {
	Value []int
	Valid bool
}

func (t IntArray) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *IntArray) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Int8Array struct {
	Value []int8
	Valid bool
}

func (t Int8Array) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Int8Array) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Int16Array struct {
	Value []int16
	Valid bool
}

func (t Int16Array) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Int16Array) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Int32Array struct {
	Value []int32
	Valid bool
}

func (t Int32Array) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Int32Array) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Int64Array struct {
	Value []int64
	Valid bool
}

func (t Int64Array) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Int64Array) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type UintArray struct {
	Value []uint
	Valid bool
}

func (t UintArray) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *UintArray) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Uint8Array struct {
	Value []uint8
	Valid bool
}

func (t Uint8Array) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Uint8Array) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Uint16Array struct {
	Value []uint16
	Valid bool
}

func (t Uint16Array) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Uint16Array) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Uint32Array struct {
	Value []uint32
	Valid bool
}

func (t Uint32Array) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Uint32Array) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Uint64Array struct {
	Value []uint64
	Valid bool
}

func (t Uint64Array) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Uint64Array) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Float32Array struct {
	Value []float32
	Valid bool
}

func (t Float32Array) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Float32Array) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Float64Array struct {
	Value []float64
	Valid bool
}

func (t Float64Array) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Float64Array) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Complex64Array struct {
	Value []complex64
	Valid bool
}

func (t Complex64Array) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Complex64Array) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type Complex128Array struct {
	Value []complex128
	Valid bool
}

func (t Complex128Array) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *Complex128Array) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type StringArray struct {
	Value []string
	Valid bool
}

func (t StringArray) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *StringArray) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}

type TimeArray struct {
	Value []time.Time
	Valid bool
}

func (t TimeArray) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *TimeArray) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}
