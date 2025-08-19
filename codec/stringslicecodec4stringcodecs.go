package codec

import "fmt"

// StringCodecs is a slice of StringCodec, which implements StringSliceCodec.
type StringCodecs []StringCodec

func (sc StringCodecs) ToStringSlice() ([]string, error) {
	values := make([]string, len(sc))
	for i, s := range sc {
		if value, err := s.ToString(); err != nil {
			return nil, fmt.Errorf("cannot stringify %q at index %d: %w", s, i, err)
		} else {
			values[i] = value
		}
	}
	return values, nil
}

func (sc StringCodecs) FromStringSlice(values []string) error {
	for i, s := range values {
		if err := sc[i].FromString(s); err != nil {
			return fmt.Errorf("cannot convert from string %q at index %d: %w", s, i, err)
		}
	}
	return nil
}
