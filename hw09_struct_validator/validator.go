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
	ErrMethodTypeIsNotSupported           = errors.New("method type is not supported")
	ErrInputIsNotAStruct                  = errors.New("input is not a struct")
	ErrFieldIsNotSupportedType            = errors.New("field is not supported type")
	ErrFieldIsNotAString                  = errors.New("field is not a string")
	ErrFieldHasInvalidLength              = errors.New("field has invalid length")
	ErrFieldDoesNotMatchRegularExpression = errors.New("field does not match regular expression")
	ErrFieldIsNotInTheAllowedValues       = errors.New("field is not in the allowed values")
	ErrFieldIsLessThanTheMinimumValue     = errors.New("field is less than the minimum value")
	ErrFieldIsNotAInt                     = errors.New("field is not a int")
	ErrFieldIsGreaterThanTheMaximumValue  = errors.New("field is greater than the maximum value")
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("invalid %s: %s", v.Field, v.Err)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errorsCount := len(v)
	if errorsCount == 0 {
		return "no validation errors"
	}
	result := make([]string, errorsCount)
	for i := 0; i < errorsCount; i++ {
		result[i] = v[i].Error()
	}
	return strings.Join(result, "; ")
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return errors.New("input is not a struct")
	}

	var validationErrors ValidationErrors

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i)
		tag := field.Tag.Get("validate")

		if tag == "" {
			continue
		}

		tags := strings.Split(tag, "|")
		for _, t := range tags {
			err := validateFieldOrSlice(field.Name, value.Interface(), t)
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateFieldOrSlice(fieldName string, value interface{}, tag string) error {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int:
		return validateField(tag, v)
	case reflect.String:
		return validateField(tag, v)
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i).Interface()
			err := validateFieldOrSlice(fieldName+"["+strconv.Itoa(i)+"]", elem, tag)
			if err != nil {
				return err
			}
		}
	default:
		return ErrFieldIsNotSupportedType
	}

	return nil
}

func validateField(tag string, kind reflect.Value) error {
	switch {
	case strings.HasPrefix(tag, "len:"):
		length, _ := strconv.Atoi(tag[4:])
		if kind.Kind() != reflect.String {
			return ErrFieldIsNotAString
		}
		if kind.Len() != length {
			return ErrFieldHasInvalidLength
		}
	case strings.HasPrefix(tag, "regexp:"):
		pattern := tag[7:]
		if kind.Kind() != reflect.String {
			return ErrFieldIsNotAString
		}
		match, _ := regexp.MatchString(pattern, kind.String())
		if !match {
			return ErrFieldDoesNotMatchRegularExpression
		}
	case strings.HasPrefix(tag, "in:"):
		allowed := strings.Split(tag[3:], ",")
		found := false
		if kind.Kind() != reflect.String && kind.Kind() != reflect.Int {
			return ErrFieldIsNotSupportedType
		}
		for _, allowVal := range allowed {
			allowValInt, errAllowValInt := strconv.Atoi(allowVal)
			if kind.Kind() == reflect.String && kind.String() == allowVal {
				found = true
				break
			} else if kind.Kind() == reflect.Int {
				if errAllowValInt != nil {
					return fmt.Errorf("invalid allowed value: %s", allowVal)
				}
				if int(kind.Int()) == allowValInt {
					found = true
					break
				}
			}
		}
		if !found {
			return ErrFieldIsNotInTheAllowedValues
		}
	case strings.HasPrefix(tag, "min:"):
		min, _ := strconv.Atoi(tag[4:])
		if kind.Kind() != reflect.Int {
			return ErrFieldIsNotAInt
		}
		if int(kind.Int()) < min {
			return ErrFieldIsLessThanTheMinimumValue
		}
	case strings.HasPrefix(tag, "max:"):
		max, _ := strconv.Atoi(tag[4:])
		if kind.Kind() != reflect.Int {
			return ErrFieldIsNotAInt
		}
		if int(kind.Int()) > max {
			return ErrFieldIsGreaterThanTheMaximumValue
		}
	}
	return nil
}
