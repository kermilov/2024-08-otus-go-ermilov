package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrValidationMethodIsNotSupported = errors.New("validation method is not supported")
	ErrInputIsNotAStruct              = errors.New("input is not a struct")
	ErrFieldIsNotSupportedType        = errors.New("field is not supported type")
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
	if v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i).Interface()
			err := validateFieldOrSlice(fieldName+"["+strconv.Itoa(i)+"]", elem, tag)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return validateField(tag, v)
}

func validateField(tag string, kind reflect.Value) error {
	tags := strings.Split(tag, ":")
	validateRegistrar, isOk := validateRegistry[tags[0]]
	if !isOk {
		return ErrValidationMethodIsNotSupported
	}
	if !validateRegistrar.canValidate(kind.Kind()) {
		return ErrFieldIsNotSupportedType
	}
	return validateRegistrar.validate(tags[1], kind)
}
