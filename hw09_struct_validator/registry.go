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
	ErrFieldHasInvalidLength              = errors.New("field has invalid length")
	ErrFieldDoesNotMatchRegularExpression = errors.New("field does not match regular expression")
	ErrFieldIsNotInTheAllowedValues       = errors.New("field is not in the allowed values")
	ErrFieldIsLessThanTheMinimumValue     = errors.New("field is less than the minimum value")
	ErrFieldIsGreaterThanTheMaximumValue  = errors.New("field is greater than the maximum value")
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
	length, _ := strconv.Atoi(s)
	if kind.Len() != length {
		return ErrFieldHasInvalidLength
	}
	return nil
}

func regexpValidate(s string, kind reflect.Value) error {
	match, _ := regexp.MatchString(s, kind.String())
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
	minValue, _ := strconv.Atoi(s)
	if int(kind.Int()) < minValue {
		return ErrFieldIsLessThanTheMinimumValue
	}
	return nil
}

func maxValidate(s string, kind reflect.Value) error {
	maxValue, _ := strconv.Atoi(s)
	if int(kind.Int()) > maxValue {
		return ErrFieldIsGreaterThanTheMaximumValue
	}
	return nil
}
