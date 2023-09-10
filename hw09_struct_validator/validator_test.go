package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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
			User{
				ID:     "1",
				Name:   "test",
				Age:    1,
				Email:  "testtest.ru",
				Role:   "test",
				Phones: []string{"1111111111"},
				meta:   json.RawMessage(""),
			},
			ValidationErrors{
				ValidationError{"ID", errLenString},
				ValidationError{"Age", errValueIsLessThanMinValue},
				ValidationError{"Email", errRegexpString},
				ValidationError{"Role", errValueDoesNotMatchSpecifiedValidator},
				ValidationError{"Phones", errLenString},
			},
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
				Code: 400,
				Body: "",
			},
			ValidationErrors{
				ValidationError{"Code", errValueDoesNotMatchSpecifiedValidator},
			},
		},
		{
			Response{
				Code: 200,
				Body: "",
			},
			ValidationErrors{},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			require.EqualError(t, Validate(tt.in), tt.expectedErr.Error())
		})
	}
}
