package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrValidationInvalid              = errors.New("validation failed")
	ErrValidationMethodIsNotSupported = fmt.Errorf("%w : validation method is not supported", ErrValidationInvalid)
	ErrInputIsNotAStruct              = fmt.Errorf("%w : input is not a struct", ErrValidationInvalid)
	ErrFieldIsNotSupportedType        = fmt.Errorf("%w : field is not supported type", ErrValidationInvalid)
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
		tag := field.Tag.Get("validate")

		if tag == "" {
			continue
		}
		value := val.Field(i)
		v := reflect.ValueOf(value.Interface())

		tags := strings.Split(tag, "|")
		for _, tag := range tags {
			split := strings.Split(tag, ":")
			validateRegistrar, isOk := validateRegistry[split[0]]
			if !isOk {
				return ErrValidationMethodIsNotSupported
			}
			err := validateField(validateRegistrar, field.Name, v, split)
			if err != nil {
				if errors.Is(err, ErrValidationFailed) {
					validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
				} else {
					return err
				}
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateField(validateRegistrar validateRegistrar, fieldName string, v reflect.Value, tag []string) error {
	if v.Kind() == reflect.Slice {
		x := v.Type().Elem().Kind()
		if !validateRegistrar.canValidate(x) {
			return ErrFieldIsNotSupportedType
		}
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i).Interface()
			err := validateField(validateRegistrar, fieldName+"["+strconv.Itoa(i)+"]", reflect.ValueOf(elem), tag)
			if err != nil {
				return err
			}
		}
		return nil
	}
	if !validateRegistrar.canValidate(v.Kind()) {
		return ErrFieldIsNotSupportedType
	}
	return validateRegistrar.validate(tag[1], v)
}
