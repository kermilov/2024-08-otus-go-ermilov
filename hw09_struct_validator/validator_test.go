package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	UserStringLen struct {
		ID     string `json:"id" validate:"len:wrong"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	UserMinForSting struct {
		ID     string          `json:"id" validate:"len:36"`
		Name   string          `validate:"min:18"`
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	UserMaxForSting struct {
		ID     string          `json:"id" validate:"len:36"`
		Name   string          `validate:"max:50"`
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	UserStingMin struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:wrong|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	UserWrongRegExp struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:(?<invalid_group.*"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}
)

func TestSuccessValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			"success User",
			User{
				ID:    "123e4567-e89b-12d3-a456-426655440000",
				Name:  "John Doe",
				Age:   42,
				Email: "H2YtY@example.com",
				Role:  "admin",
				Phones: []string{
					"89998723412",
				},
			},
			nil,
		},
		{
			"success App",
			App{
				Version: "1.0.0",
			},
			nil,
		},
		{
			"success Token",
			Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			nil,
		},
		{
			"success Response",
			Response{
				Code: 200,
				Body: "Body",
			},
			nil,
		},
	}

	runTests(t, tests)
}

func TestErrorsValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			"error User",
			User{
				ID:    "123e4567-e89b-12d3-a456-42665544000",
				Name:  "John Doe",
				Age:   52,
				Email: "H2YtYexample.com",
				Role:  "wrong",
				Phones: []string{
					"89998723412",
					"189998723412",
				},
			},
			ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrFieldHasInvalidLength,
				},
				ValidationError{
					Field: "Age",
					Err:   ErrFieldIsGreaterThanTheMaximumValue,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrFieldDoesNotMatchRegularExpression,
				},
				ValidationError{
					Field: "Role",
					Err:   ErrFieldIsNotInTheAllowedValues,
				},
				ValidationError{
					Field: "Phones",
					Err:   ErrFieldHasInvalidLength,
				},
			},
		},
		{
			"error App",
			App{
				Version: "11.0.0",
			},
			ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   ErrFieldHasInvalidLength,
				},
			},
		},
		{
			"error Response",
			Response{
				Code: 204,
			},
			ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   ErrFieldIsNotInTheAllowedValues,
				},
			},
		},
	}

	runTests(t, tests)
}

func TestErrorsValidateInvalid(t *testing.T) {
	tests := []struct {
		name string
		in   interface{}
	}{
		{
			"UserStringLen",
			UserStringLen{
				ID:    "123e4567-e89b-12d3-a456-42665544000",
				Name:  "John Doe",
				Age:   52,
				Email: "H2YtYexample.com",
				Role:  "wrong",
				Phones: []string{
					"89998723412",
					"189998723412",
				},
			},
		},
		{
			"UserMinForSting",
			UserMinForSting{
				ID:    "123e4567-e89b-12d3-a456-42665544000",
				Name:  "John Doe",
				Age:   52,
				Email: "H2YtYexample.com",
				Role:  "wrong",
				Phones: []string{
					"89998723412",
					"189998723412",
				},
			},
		},
		{
			"UserMaxForSting",
			UserMaxForSting{
				ID:    "123e4567-e89b-12d3-a456-42665544000",
				Name:  "John Doe",
				Age:   52,
				Email: "H2YtYexample.com",
				Role:  "wrong",
				Phones: []string{
					"89998723412",
					"189998723412",
				},
			},
		},
		{
			"UserStingMin",
			UserStingMin{
				ID:    "123e4567-e89b-12d3-a456-42665544000",
				Name:  "John Doe",
				Age:   52,
				Email: "H2YtYexample.com",
				Role:  "wrong",
				Phones: []string{
					"89998723412",
					"189998723412",
				},
			},
		},
		{
			"UserWrongRegExp",
			UserWrongRegExp{
				ID:    "123e4567-e89b-12d3-a456-42665544000",
				Name:  "John Doe",
				Age:   52,
				Email: "H2YtYexample.com",
				Role:  "wrong",
				Phones: []string{
					"89998723412",
					"189998723412",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("case %s", tt.name), func(t *testing.T) {
			t.Parallel()

			actualErr := Validate(tt.in)
			require.NotNil(t, actualErr)
			require.True(t, !errors.Is(actualErr, ErrValidationFailed))
			var validationErrors ValidationErrors
			require.True(t, !errors.As(actualErr, &validationErrors))
		})
	}
}

func runTests(t *testing.T, tests []struct {
	name        string
	in          interface{}
	expectedErr error
},
) {
	t.Helper()
	for _, tt := range tests {
		t.Run(fmt.Sprintf("case %s", tt.name), func(t *testing.T) {
			t.Parallel()

			actualErr := Validate(tt.in)
			require.Equal(t, tt.expectedErr, actualErr)
		})
	}
}
