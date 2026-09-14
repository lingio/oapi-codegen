package codegen

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

const (
	extPropGoType                    = "x-go-type"
	extPropOmitEmpty                 = "x-omitempty"
	extPropGoTypeSkipOptionalPointer = "x-go-type-skip-optional-pointer"
)

// extAs converts the value of an OpenAPI extension into T. kin-openapi hands us
// values which have already been decoded into Go types, but older versions (and
// specs loaded by hand) may still carry raw JSON, so we accept both.
func extAs[T any](extPropValue interface{}) (T, error) {
	var out T

	if extPropValue == nil {
		return out, fmt.Errorf("failed to convert type: %T", extPropValue)
	}

	switch v := extPropValue.(type) {
	case T:
		return v, nil
	case json.RawMessage:
		if err := json.Unmarshal(v, &out); err != nil {
			return out, errors.Wrap(err, "failed to unmarshal json")
		}
		return out, nil
	}

	// Fall back to a JSON round-trip, which copes with values decoded into a
	// different but compatible type, eg. a float64 for an integer.
	raw, err := json.Marshal(extPropValue)
	if err != nil {
		return out, fmt.Errorf("failed to convert type: %T", extPropValue)
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return out, errors.Wrap(err, "failed to unmarshal json")
	}
	return out, nil
}

func extTypeName(extPropValue interface{}) (string, error) {
	return extAs[string](extPropValue)
}

func extParsePropGoTypeSkipOptionalPointer(extPropValue interface{}) (bool, error) {
	return extAs[bool](extPropValue)
}

func extParseOmitEmpty(extPropValue interface{}) (bool, error) {
	return extAs[bool](extPropValue)
}
