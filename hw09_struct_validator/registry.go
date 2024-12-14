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
	parse    func(string) (interface{}, error)
	validate func(interface{}, reflect.Value) error
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
		parse:    lenParse,
		validate: lenValidate,
	},
	"regexp": {
		kinds:    []reflect.Kind{reflect.String},
		parse:    regexpParse,
		validate: regexpValidate,
	},
	"in": {
		kinds:    []reflect.Kind{reflect.String, reflect.Int},
		parse:    inParse,
		validate: inValidate,
	},
	"min": {
		kinds:    []reflect.Kind{reflect.Int},
		parse:    minParse,
		validate: minValidate,
	},
	"max": {
		kinds:    []reflect.Kind{reflect.Int},
		parse:    maxParse,
		validate: maxValidate,
	},
}

func lenParse(s string) (interface{}, error) {
	return strconv.Atoi(s)
}

func lenValidate(length interface{}, kind reflect.Value) error {
	if kind.Len() != length.(int) {
		return ErrFieldHasInvalidLength
	}
	return nil
}

func regexpParse(s string) (interface{}, error) {
	return regexp.Compile(s)
}

func regexpValidate(regExp interface{}, kind reflect.Value) error {
	if !regExp.(*regexp.Regexp).MatchString(kind.String()) {
		return ErrFieldDoesNotMatchRegularExpression
	}
	return nil
}

func inParse(s string) (interface{}, error) {
	return strings.Split(s, ","), nil
}

func inValidate(allowed interface{}, kind reflect.Value) error {
	for _, allowVal := range allowed.([]string) {
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

func minParse(s string) (interface{}, error) {
	return strconv.Atoi(s)
}

func minValidate(minValue interface{}, kind reflect.Value) error {
	if int(kind.Int()) < minValue.(int) {
		return ErrFieldIsLessThanTheMinimumValue
	}
	return nil
}

func maxParse(s string) (interface{}, error) {
	return strconv.Atoi(s)
}

func maxValidate(maxValue interface{}, kind reflect.Value) error {
	if int(kind.Int()) > maxValue.(int) {
		return ErrFieldIsGreaterThanTheMaximumValue
	}
	return nil
}
