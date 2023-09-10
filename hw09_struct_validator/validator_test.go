package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type UserRole string

type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
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
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			1,
			ErrNotStruct,
		},
		{
			App{
				Version: "12345",
			},
			ValidationErrors{},
		},
		{
			Token{
				Header:    []byte{1, 2},
				Payload:   []byte{3, 4},
				Signature: []byte{5, 6},
			},
			ValidationErrors{},
		},
		{
			Response{
				Code: 200,
				Body: "",
			},
			ValidationErrors{},
		},
		{
			User{
				ID:     "1",
				Name:   "test",
				Age:    1,
				Email:  "testtest.ru",
				Role:   "admin",
				Phones: []string{"1111111111"},
				meta:   json.RawMessage(""),
			},
			ValidationErrors{
				ValidationError{"ID", errLenString},
				ValidationError{"Age", errValueIsLessThanMinValue},
				ValidationError{"Email", errRegexpString},
				ValidationError{"Phones", errLenString},
			},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			if validationErr, ok := err.(ValidationErrors); ok {
				for _, e := range validationErr {
					if e.Field == "Code" && !errors.Is(e.Err, tt.expectedErr) {
						t.Errorf("Error: Expected: %v, but received: %v", tt.expectedErr, e.Err)
					}
				}
			} else {
				// Обработать случай, когда err - это одиночная ошибка, а не ValidationErrors.
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("Error: Expected: %v, but received: %v", tt.expectedErr, err)
				}
			}

		})
	}
}
