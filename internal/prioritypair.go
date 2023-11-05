package internal

import (
	"fmt"
	"reflect"
)

// PriorityPair maps a reflect.Type to a pair of (primary, secondary).
type PriorityPair map[reflect.Type][2]any

func NewPriorityPair() PriorityPair {
	return make(PriorityPair)
}

func (m PriorityPair) SetPair(typ reflect.Type, primary, secondary any, ignoreConflict bool) error {
	olds, ok := m[typ]

	if !ok {
		m[typ] = [2]any{primary, secondary}
		return nil
	}

	oldPrimary, _ := olds[0], olds[1]
	if primary != nil { // set primary
		if oldPrimary != nil && !ignoreConflict { // conflict
			return fmt.Errorf("duplicate type: %q", typ)
		}
		olds[0] = primary
	}

	if secondary != nil { // always set secondary
		olds[1] = secondary
	}
	m[typ] = olds // NOTE: this is necessary
	return nil
}

// GetOne returns the primary if it exists, otherwise the secondary.
func (m PriorityPair) GetOne(t reflect.Type) any {
	if pair, ok := m[t]; ok {
		if pair[0] != nil {
			return pair[0]
		}
		return pair[1]
	} else {
		return nil
	}
}

func (m PriorityPair) GetPrimary(t reflect.Type) any {
	if pair, ok := m[t]; ok {
		return pair[0]
	} else {
		return nil
	}
}

func (m PriorityPair) GetSecondary(t reflect.Type) any {
	if pair, ok := m[t]; ok {
		return pair[1]
	} else {
		return nil
	}
}
