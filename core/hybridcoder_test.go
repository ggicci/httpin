package core

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHybridCoder_MarshalText_only(t *testing.T) {
	apple := &textMarshalerApple{}
	rv := reflect.ValueOf(apple)
	stringable := hybridizeCoder(rv)
	assert.NotNil(t, stringable)

	text, err := stringable.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "apple", text)

	assert.ErrorContains(t, stringable.FromString("red apple"), "StringUnmarshaler not implemented")
}

func TestHybridCoder_UnmarshalText_only(t *testing.T) {
	banana := &textUnmarshalerBanana{}
	rv := reflect.ValueOf(banana)
	stringable := hybridizeCoder(rv)
	assert.NotNil(t, stringable)

	text, err := stringable.ToString()
	assert.ErrorContains(t, err, "StringMarshaler not implemented")
	assert.Empty(t, text)

	err = stringable.FromString("yellow banana")
	assert.NoError(t, err)
	assert.Equal(t, "yellow banana", banana.Content)
}

func TestHybridCoder_MarshalText_and_UnmarshalText(t *testing.T) {
	orange := &textMarshalerAndUnmarshalerOrange{Content: "orange"}
	rv := reflect.ValueOf(orange)
	stringable := hybridizeCoder(rv)
	assert.NotNil(t, stringable)

	text, err := stringable.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "orange", text)

	err = stringable.FromString("red orange")
	assert.NoError(t, err)
	assert.Equal(t, "red orange", orange.Content)
}

func TestHybridCoder_StringMarshaler_takePrecedence(t *testing.T) {
	peach := &stringMarshalerAndTextMarshalerPeach{Content: "peach"}
	rv := reflect.ValueOf(peach)
	stringable := hybridizeCoder(rv)
	assert.NotNil(t, stringable)

	text, err := stringable.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "ToString:peach", text)

	err = stringable.FromString("red peach")
	assert.ErrorContains(t, err, "StringUnmarshaler not implemented")
}

func TestHybridCoder_StringUnmarshaler_takePrecedence(t *testing.T) {
	peach := &stringUnmarshalerAndTextUnmarshalerPeach{Content: "peach"}
	rv := reflect.ValueOf(peach)
	stringable := hybridizeCoder(rv)
	assert.NotNil(t, stringable)

	text, err := stringable.ToString()
	assert.ErrorContains(t, err, "StringMarshaler not implemented")
	assert.Empty(t, text)

	err = stringable.FromString("red peach")
	assert.NoError(t, err)
	assert.Equal(t, "FromString:red peach", peach.Content)
}

func TestHybridCoder_StringMarshaler_and_TextUnmarshaler(t *testing.T) {
	pineapple := &stringMarshalerAndTextUnmarshalerPineapple{Content: "pineapple"}
	rv := reflect.ValueOf(pineapple)
	stringable := hybridizeCoder(rv)
	assert.NotNil(t, stringable)

	text, err := stringable.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "ToString:pineapple", text)

	err = stringable.FromString("red pineapple")
	assert.NoError(t, err)
	assert.Equal(t, "UnmarshalText:red pineapple", pineapple.Content)
}

func TestHybridCoder_MarshalText_Error(t *testing.T) {
	watermelon := &textMarshalerSpoiledWatermelon{}
	rv := reflect.ValueOf(watermelon)
	stringable := hybridizeCoder(rv)
	assert.NotNil(t, stringable)

	text, err := stringable.ToString()
	assert.ErrorContains(t, err, "spoiled")
	assert.Empty(t, text)
}

func TestHybridCoder_cannotInterfaceError(t *testing.T) {
	type mystruct struct {
		unexportedName string
	}
	v := mystruct{unexportedName: "mystruct"}
	rv := reflect.ValueOf(v)

	stringable := hybridizeCoder(rv.Field(0))
	assert.Nil(t, stringable)
}

func TestHybridCoder_nilOnNoInterfacesDetected(t *testing.T) {
	var zero zeroInterface
	rv := reflect.ValueOf(zero)

	stringable := hybridizeCoder(rv)
	assert.Nil(t, stringable)
}

type textMarshalerApple struct{} // only implements encoding.TextMarshaler

func (t *textMarshalerApple) MarshalText() ([]byte, error) {
	return []byte("apple"), nil
}

type textUnmarshalerBanana struct{ Content string } // only implements encoding.TextUnmarshaler

func (t *textUnmarshalerBanana) UnmarshalText(text []byte) error {
	t.Content = string(text)
	return nil
}

type textMarshalerAndUnmarshalerOrange struct{ Content string } // implements both encoding.TextMarshaler and encoding.TextUnmarshaler

func (t *textMarshalerAndUnmarshalerOrange) MarshalText() ([]byte, error) {
	return []byte(t.Content), nil
}

func (t *textMarshalerAndUnmarshalerOrange) UnmarshalText(text []byte) error {
	t.Content = string(text)
	return nil
}

// implements internal.StringMarshaler and encoding.TextMarshaler
// will use internal.StringMarshaler
type stringMarshalerAndTextMarshalerPeach struct{ Content string }

func (s *stringMarshalerAndTextMarshalerPeach) ToString() (string, error) {
	return "ToString:" + s.Content, nil
}

func (s *stringMarshalerAndTextMarshalerPeach) MarshalText() ([]byte, error) {
	return []byte("MarshalText:" + s.Content), nil
}

// implements internal.StringUnmarshaler and encoding.TextUnmarshaler
// will use internal.StringUnmarshaler
type stringUnmarshalerAndTextUnmarshalerPeach struct{ Content string }

func (s *stringUnmarshalerAndTextUnmarshalerPeach) FromString(text string) error {
	s.Content = "FromString:" + text
	return nil
}

func (s *stringUnmarshalerAndTextUnmarshalerPeach) UnmarshalText(text []byte) error {
	s.Content = "UnmarshalText:" + string(text)
	return nil
}

type stringMarshalerAndTextUnmarshalerPineapple struct{ Content string }

func (s *stringMarshalerAndTextUnmarshalerPineapple) ToString() (string, error) {
	return "ToString:" + s.Content, nil
}

func (s *stringMarshalerAndTextUnmarshalerPineapple) UnmarshalText(text []byte) error {
	s.Content = "UnmarshalText:" + string(text)
	return nil
}

type textMarshalerSpoiledWatermelon struct{}

func (t *textMarshalerSpoiledWatermelon) MarshalText() ([]byte, error) {
	return nil, errors.New("spoiled")
}

type zeroInterface struct{}
