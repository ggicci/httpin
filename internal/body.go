package internal

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
)

// BodyEncodeDecoder is the interface for encoding and decoding the request body.
// Common body formats are: json, xml, yaml, etc.
type BodyEncodeDecoder interface {
	// Decode decodes the request body into the specified object.
	Decode(src io.Reader, dst any) error
	// Encode encodes the specified object into a reader for the request body.
	Encode(src any) (io.Reader, error)
}

type JSONBody struct{}

func (de *JSONBody) Decode(src io.Reader, dst any) error {
	return json.NewDecoder(src).Decode(dst)
}

func (en *JSONBody) Encode(src any) (io.Reader, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(src); err != nil {
		return nil, err
	}
	return &buf, nil
}

type XMLBody struct{}

func (de *XMLBody) Decode(src io.Reader, dst any) error {
	return xml.NewDecoder(src).Decode(dst)
}

func (en *XMLBody) Encode(src any) (io.Reader, error) {
	var buf bytes.Buffer
	if err := xml.NewEncoder(&buf).Encode(src); err != nil {
		return nil, err
	}
	return &buf, nil
}
