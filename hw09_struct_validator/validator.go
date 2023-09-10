package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrNotStruct                           = fmt.Errorf(" variable not a struct")
	errNoRule                              = fmt.Errorf(" not find rule for tag validate")
	errLenString                           = fmt.Errorf(" string length > or < validation values")
	errRegexpString                        = fmt.Errorf(" the string does not match the specified Regexp")
	errValueIsLessThanMinValue             = fmt.Errorf(" the value is less than the minimum value")
	errValueIsMoreThanMaxValue             = fmt.Errorf(" the value is more than the maximum value")
	errValueDoesNotMatchSpecifiedValidator = fmt.Errorf(" does not match the specified validator")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationError) Error() string {
	return v.Err.Error()
}

func (v ValidationErrors) Error() string {
	builder := strings.Builder{}
	for i, e := range v {
		builder.WriteString(strconv.Itoa(i+1) + ") " + e.Field + ":" + e.Err.Error() + "\n")
	}
	return builder.String()
}

func Validate(v interface{}) error {
	var errSlice ValidationErrors

	rType := reflect.TypeOf(v)
	rValue := reflect.ValueOf(v)
	if rType.Kind().String() != "struct" {
		return ErrNotStruct
	}

	for i := 0; i < rType.NumField(); i++ {
		xType := rType.Field(i)
		xValue := rValue.Field(i)
		tagValue := xType.Tag.Get("validate")

		if tagValue == "" {
			continue
		}

		switch xType.Type.String() {
		case "string":
			err := tagStringValidate(xValue.String(), tagValue)
			if err != nil {
				errSlice = append(errSlice, ValidationError{
					Field: xType.Name,
					Err:   err,
				})
			}
		case "int":
			err := tagIntValidate(xValue.Interface().(int), tagValue)
			if err != nil {
				errSlice = append(errSlice, ValidationError{
					Field: xType.Name,
					Err:   err,
				})
			}
		case "[]int":
			for _, item := range xValue.Interface().([]int) {
				err := tagIntValidate(item, tagValue)
				if err != nil {
					errSlice = append(errSlice, ValidationError{
						Field: xType.Name,
						Err:   err,
					})
				}
			}
		case "[]string":
			for _, item := range xValue.Interface().([]string) {
				err := tagStringValidate(item, tagValue)
				if err != nil {
					errSlice = append(errSlice, ValidationError{
						Field: xType.Name,
						Err:   err,
					})
				}
			}
		default:
			continue
		}
	}
	return errSlice
}

func tagStringValidate(data string, tag string) error {
	anyRule := strings.Split(tag, "|")

	for _, value := range anyRule {
		rule := strings.Split(value, ":")

		switch rule[0] {
		case "len":
			ato, err := strconv.Atoi(rule[1])
			if err != nil {
				return fmt.Errorf("could not ato : %w", err)
			}
			if len(data) != ato {
				return errLenString
			}
		case "regexp":
			matchString, err := regexp.MatchString(rule[1], data)
			if err != nil {
				return fmt.Errorf("could not rexexp match : %w", err)
			}

			if !matchString {
				return errRegexpString
			}
		case "in":
			for _, item := range strings.Split(rule[1], ",") {
				if item != data {
					return errValueDoesNotMatchSpecifiedValidator
				}
				return nil
			}
		default:
			return errNoRule
		}
	}
	return nil
}

func tagIntValidate(data int, tag string) error {
	anyRule := strings.Split(tag, "|")

	for _, value := range anyRule {
		rule := strings.Split(value, ":")

		switch rule[0] {
		case "min":
			ato, err := strconv.Atoi(rule[1])
			if err != nil {
				return fmt.Errorf("could not ato : %w", err)
			}
			if data < ato {
				return errValueIsLessThanMinValue
			}
		case "max":
			ato, err := strconv.Atoi(rule[1])
			if err != nil {
				return fmt.Errorf("could not ato : %w", err)
			}
			if data > ato {
				return errValueIsMoreThanMaxValue
			}
		case "in":
			for _, item := range strings.Split(rule[1], ",") {
				if item != strconv.Itoa(data) {
					return errValueDoesNotMatchSpecifiedValidator
				}
				return nil
			}
		default:
			return errNoRule
		}
	}
	return nil
}
