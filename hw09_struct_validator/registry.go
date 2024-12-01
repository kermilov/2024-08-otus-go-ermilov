package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrValidationFailed                   = errors.New("validation failed")
	ErrFieldHasInvalidLength              = fmt.Errorf("%w : field has invalid length", ErrValidationFailed)
	ErrFieldDoesNotMatchRegularExpression = fmt.Errorf("%w : field does not match regular expression", ErrValidationFailed)
	ErrFieldIsNotInTheAllowedValues       = fmt.Errorf("%w : field is not in the allowed values", ErrValidationFailed)
	ErrFieldIsLessThanTheMinimumValue     = fmt.Errorf("%w : field is less than the minimum value", ErrValidationFailed)
	ErrFieldIsGreaterThanTheMaximumValue  = fmt.Errorf("%w : field is greater than the maximum value", ErrValidationFailed)
)

type validateRegistrar struct {
	kinds    []reflect.Kind
	validate func(string, reflect.Value) error
}

func (v validateRegistrar) canValidate(kind reflect.Kind) bool {
	for _, k := range v.kinds {
		if k == kind {
			return true
		}
	}
	return false
}

var validateRegistry = map[string]validateRegistrar{
	"len": {
		kinds:    []reflect.Kind{reflect.String},
		validate: lenValidate,
	},
	"regexp": {
		kinds:    []reflect.Kind{reflect.String},
		validate: regexpValidate,
	},
	"in": {
		kinds:    []reflect.Kind{reflect.String, reflect.Int},
		validate: inValidate,
	},
	"min": {
		kinds:    []reflect.Kind{reflect.Int},
		validate: minValidate,
	},
	"max": {
		kinds:    []reflect.Kind{reflect.Int},
		validate: maxValidate,
	},
}

func lenValidate(s string, kind reflect.Value) error {
	length, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("invalid length: %s", s)
	}
	if kind.Len() != length {
		return ErrFieldHasInvalidLength
	}
	return nil
}

func regexpValidate(s string, kind reflect.Value) error {
	match, err := regexp.MatchString(s, kind.String())
	if err != nil {
		return fmt.Errorf("invalid regular expression: %s", s)
	}
	if !match {
		return ErrFieldDoesNotMatchRegularExpression
	}
	return nil
}

func inValidate(s string, kind reflect.Value) error {
	allowed := strings.Split(s, ",")
	for _, allowVal := range allowed {
		if kind.Kind() == reflect.String && kind.String() == allowVal {
			return nil
		} else if kind.Kind() == reflect.Int {
			allowValInt, errAllowValInt := strconv.Atoi(allowVal)
			if errAllowValInt != nil {
				return fmt.Errorf("invalid allowed value: %s", allowVal)
			}
			if int(kind.Int()) == allowValInt {
				return nil
			}
		}
	}
	return ErrFieldIsNotInTheAllowedValues
}

func minValidate(s string, kind reflect.Value) error {
	minValue, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("invalid minimum value: %s", s)
	}
	if int(kind.Int()) < minValue {
		return ErrFieldIsLessThanTheMinimumValue
	}
	return nil
}

func maxValidate(s string, kind reflect.Value) error {
	maxValue, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("invalid maximum value: %s", s)
	}
	if int(kind.Int()) > maxValue {
		return ErrFieldIsGreaterThanTheMaximumValue
	}
	return nil
}
