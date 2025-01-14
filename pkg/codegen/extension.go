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

func extTypeName(extPropValue interface{}) (string, error) {
	raw, ok := extPropValue.(json.RawMessage)
	if !ok {
		return "", fmt.Errorf("failed to convert type: %T", extPropValue)
	}
	var name string
	if err := json.Unmarshal(raw, &name); err != nil {
		return "", errors.Wrap(err, "failed to unmarshal json")
	}

	return name, nil
}

func extParsePropGoTypeSkipOptionalPointer(extPropValue interface{}) (bool, error) {
	raw, ok := extPropValue.(json.RawMessage)
	if !ok {
		return false, fmt.Errorf("failed to convert type: %T", extPropValue)
	}

	var goTypeSkipOptionalPointer bool
	if err := json.Unmarshal(raw, &goTypeSkipOptionalPointer); err != nil {
		return false, errors.Wrap(err, "failed to unmarshal json")
	}
	return goTypeSkipOptionalPointer, nil
}

func extParseOmitEmpty(extPropValue interface{}) (bool, error) {

	raw, ok := extPropValue.(json.RawMessage)
	if !ok {
		return false, fmt.Errorf("failed to convert type: %T", extPropValue)
	}

	var omitEmpty bool
	if err := json.Unmarshal(raw, &omitEmpty); err != nil {
		return false, errors.Wrap(err, "failed to unmarshal json")
	}

	return omitEmpty, nil
}
