package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
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

		conditions := strings.Split(tag, "|")
		for _, condition := range conditions {
			split := strings.Split(condition, ":")
			validator, isOk := validateRegistry[split[0]]
			if !isOk {
				return ErrValidationMethodIsNotSupported
			}
			parse, errParse := validator.parse(split[1])
			if errParse != nil {
				return errParse
			}
			err := validateField(validator, value, parse)
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

func validateField(validator validateRegistrar, value reflect.Value, parse interface{}) error {
	if value.Kind() == reflect.Slice {
		sliceKind := value.Type().Elem().Kind()
		if !validator.canValidate(sliceKind) {
			return ErrFieldIsNotSupportedType
		}
		for i := 0; i < value.Len(); i++ {
			elem := value.Index(i).Interface()
			err := validateField(validator, reflect.ValueOf(elem), parse)
			if err != nil {
				return err
			}
		}
		return nil
	}
	if !validator.canValidate(value.Kind()) {
		return ErrFieldIsNotSupportedType
	}
	return validator.validate(parse, value)
}
