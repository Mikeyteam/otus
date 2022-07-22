package hw09structvalidator

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	stringType      = "string"
	stringSliceType = "[]string"
	intType         = "int"
	intSliceType    = "[]int"
)

var (
	validatorKey = "validate"
	ErrorRegexp  = errors.New("error regular expression")
)

var validationFuncMap = map[string]validationFunc{
	"len":    validateLen,
	"regexp": validateRegexp,
	"min":    validateMaxOrMin,
	"max":    validateMaxOrMin,
	"in":     validateIn,
}

type ErrorValidation struct {
	Field string
	Err   error
}

type (
	ErrorsValidation []ErrorValidation
	ErrorValidator   struct {
		message string
	}
)

type validationFunc func(tagKeyMethod string,
	structFieldName string,
	validationLimit string,
	validateValue reflect.Value,
	typeValue reflect.Type) error

func (e LenError) Error() string {
	builder := strings.Builder{}
	builder.WriteString("invalid actual field len is ")
	builder.WriteString(strconv.Itoa(e.CurrentValue))
	builder.WriteString(" should be equal ")
	builder.WriteString(strconv.Itoa(e.Limit))
	return builder.String()
}

func (e MaxError) Error() string {
	builder := strings.Builder{}
	builder.WriteString("value is ")
	builder.WriteString(strconv.Itoa(e.CurrentValue))
	builder.WriteString(" should be less than or equal to ")
	builder.WriteString(strconv.Itoa(e.Limit))
	return builder.String()
}

func (e MinError) Error() string {
	builder := strings.Builder{}
	builder.WriteString("value is ")
	builder.WriteString(strconv.Itoa(e.CurrentValue))
	builder.WriteString(" should be greater than or equal to ")
	builder.WriteString(strconv.Itoa(e.Limit))
	return builder.String()
}

func (e InError) Error() string {
	builder := strings.Builder{}
	builder.WriteString("value is ")
	builder.WriteString(e.CurrentValue)
	builder.WriteString(" should be one of ")
	builder.WriteString(e.Limit)
	return builder.String()
}

func (e ErrorValidation) Error() string {
	builder := strings.Builder{}
	builder.WriteString(e.Field)
	builder.WriteString(": ")
	builder.WriteString(e.Err.Error())
	return builder.String()
}

func (e ErrorValidator) Error() string {
	return e.message
}

type LenError struct {
	Limit        int
	CurrentValue int
}

type MaxError struct {
	Limit        int
	CurrentValue int
}

type MinError struct {
	Limit        int
	CurrentValue int
}

type InError struct {
	Limit        string
	CurrentValue string
}

func validateLen(method string, validationLimit string, structFieldName string,
	valueToCheck reflect.Value, valueType reflect.Type,
) error {
	limitInt, err := strconv.Atoi(validationLimit)
	if err != nil {
		return ErrorValidator{fmt.Sprintf("%s must be integer", validationLimit)}
	}
	switch {
	case valueType.String() == stringSliceType:
		elementInterface := valueToCheck.Interface()
		elementSlice, ok := elementInterface.([]string)
		if !ok {
			return ErrorValidator{fmt.Sprintf("error %#v string slice", elementInterface)}
		}
		for _, element := range elementSlice {
			if len(element) != limitInt {
				return ErrorValidation{
					Field: structFieldName,
					Err:   LenError{Limit: limitInt, CurrentValue: len(element)},
				}
			}
		}
	case valueType.Kind().String() == stringType:
		if len(valueToCheck.String()) != limitInt {
			return ErrorValidation{
				Field: structFieldName,
				Err:   LenError{Limit: limitInt, CurrentValue: len(valueToCheck.String())},
			}
		}
	default:
		return ErrorValidator{fmt.Sprintf(
			"field %s wrong type %s for method \"%s\"", structFieldName, valueType, method)}
	}

	return nil
}

func validateRegexp(method string, validationLimit string, structFieldName string,
	valueToCheck reflect.Value, valueType reflect.Type,
) error {
	regex, err := regexp.Compile(validationLimit)
	if err != nil {
		return ErrorValidator{fmt.Sprintf("error compile regex %#v", validationLimit)}
	}
	switch {
	case valueType.String() == stringSliceType:
		elementInterface := valueToCheck.Interface()
		elementSlice, ok := elementInterface.([]string)
		if !ok {
			return ErrorValidator{fmt.Sprintf("error %#v string slice", elementInterface)}
		}
		for _, elem := range elementSlice {
			if !regex.MatchString(elem) {
				return ErrorValidation{
					Field: structFieldName,
					Err:   ErrorRegexp,
				}
			}
		}
	case valueType.Kind().String() == stringType:
		if !regex.MatchString(valueToCheck.String()) {
			return ErrorValidation{
				Field: structFieldName,
				Err:   ErrorRegexp,
			}
		}

	default:
		return ErrorValidator{fmt.Sprintf(
			"field %s wrong type %s for method \"%s\"", structFieldName, valueType, method)}
	}

	return nil
}

func validateMaxOrMin(method string, validationLimit string, structFieldName string,
	valueToCheck reflect.Value, valueType reflect.Type,
) error {
	limitInt, err := strconv.Atoi(validationLimit)
	if err != nil {
		return ErrorValidator{fmt.Sprintf("%s must be integer", validationLimit)}
	}
	switch {
	case valueType.Kind().String() == intType:
		switch method {
		case "min":
			if int(valueToCheck.Int()) < limitInt {
				return ErrorValidation{
					Field: structFieldName,
					Err:   MinError{Limit: limitInt, CurrentValue: int(valueToCheck.Int())},
				}
			}
		case "max":
			if int(valueToCheck.Int()) > limitInt {
				return ErrorValidation{
					Field: structFieldName,
					Err:   MaxError{Limit: limitInt, CurrentValue: int(valueToCheck.Int())},
				}
			}
		}

	case valueType.String() == intSliceType:
		elementInterface := valueToCheck.Interface()
		elementSlice, ok := elementInterface.([]int)
		if !ok {
			return ErrorValidator{fmt.Sprintf("error %#v string slice", elementInterface)}
		}
		for _, element := range elementSlice {
			switch method {
			case "min":
				if element < limitInt {
					return ErrorValidation{
						Field: structFieldName,
						Err:   MinError{Limit: limitInt, CurrentValue: element},
					}
				}
			case "max":
				if element > limitInt {
					return ErrorValidation{
						Field: structFieldName,
						Err:   MaxError{Limit: limitInt, CurrentValue: element},
					}
				}
			}
		}
	default:
		return ErrorValidator{fmt.Sprintf(
			"field %s wrong type %s for method \"%s\"", structFieldName, valueType, method)}
	}

	return nil
}

func validateIn(method string, validationLimit string, structFieldName string,
	valueToCheck reflect.Value, valueType reflect.Type,
) error {
	limitSlice := strings.Split(validationLimit, ",")
	switch {
	case valueType.Kind().String() == intType:
		limitIntMap, err := sliceToAtoiMap(limitSlice)
		if err != nil {
			return ErrorValidator{fmt.Sprintf("for %s limit %s must be integer, %s",
				valueToCheck, validationLimit, err.Error())}
		}
		if _, ok := limitIntMap[int(valueToCheck.Int())]; !ok {
			return ErrorValidation{
				Field: structFieldName,
				Err:   InError{Limit: validationLimit, CurrentValue: fmt.Sprintf("%d", valueToCheck.Int())},
			}
		}
	case valueType.String() == intSliceType:
		limitIntMap, err := sliceToAtoiMap(limitSlice)
		if err != nil {
			return ErrorValidator{fmt.Sprintf("for %s limit %s must be integer, %s",
				valueToCheck, validationLimit, err.Error())}
		}
		elementInterface := valueToCheck.Interface()
		elementSlice, ok := elementInterface.([]int)
		if !ok {
			return ErrorValidator{fmt.Sprintf("error %#v string slice", elementInterface)}
		}
		for _, elem := range elementSlice {
			if _, ok := limitIntMap[elem]; !ok {
				return ErrorValidation{
					Field: structFieldName,
					Err:   InError{Limit: validationLimit, CurrentValue: fmt.Sprintf("%d", elem)},
				}
			}
		}
	case valueType.Kind().String() == stringType:
		limitStringMap := sliceToStringMap(limitSlice)
		if _, ok := limitStringMap[valueToCheck.String()]; !ok {
			return ErrorValidation{
				Field: structFieldName,
				Err:   InError{Limit: validationLimit, CurrentValue: valueToCheck.String()},
			}
		}
	case valueType.String() == stringSliceType:
		elementInterface := valueToCheck.Interface()
		elementSlice, ok := elementInterface.([]string)
		if !ok {
			return ErrorValidator{fmt.Sprintf("error %#v string slice", elementInterface)}
		}
		limitStringMap := sliceToStringMap(limitSlice)
		for _, element := range elementSlice {
			if _, ok = limitStringMap[element]; !ok {
				return ErrorValidation{
					Field: structFieldName,
					Err:   InError{Limit: validationLimit, CurrentValue: element},
				}
			}
		}
	default:
		return ErrorValidator{fmt.Sprintf(
			"field %s wrong type %s for method \"%s\"", structFieldName, valueType, method)}
	}

	return nil
}

func sliceToAtoiMap(sliceWithStrings []string) (map[int]struct{}, error) {
	resultIntMap := make(map[int]struct{}, len(sliceWithStrings))
	for _, strValue := range sliceWithStrings {
		intValue, err := strconv.Atoi(strValue)
		if err != nil {
			return nil, err
		}
		resultIntMap[intValue] = struct{}{}
	}
	return resultIntMap, nil
}

func sliceToStringMap(sliceWithStrings []string) map[string]struct{} {
	resultIntMap := make(map[string]struct{}, len(sliceWithStrings))
	for _, strValue := range sliceWithStrings {
		resultIntMap[strValue] = struct{}{}
	}
	return resultIntMap
}

// ErrorsValidation.Error converts slice of errors into single string.
func (v ErrorsValidation) Error() string {
	builder := strings.Builder{}
	for i, err := range v {
		builder.WriteString(err.Field)
		builder.WriteString(": ")
		builder.WriteString(err.Err.Error())
		if len(v) != 0 && i < len(v)-1 {
			builder.WriteString(" \\\n")
		}
	}
	return builder.String()
}

// Validate This function validate fields of struct, if they have "validatorKey" flag.
func Validate(v interface{}) error {
	var errs ErrorsValidation
	var valerr ErrorValidation
	rVal := reflect.ValueOf(v)
	if rVal.Kind() != reflect.Struct {
		return ErrorValidator{"Validate support for structs only"}
	}
	structRval := rVal.Type()
	errs = make(ErrorsValidation, 0, structRval.NumField())
	for i := 0; i < structRval.NumField(); i++ {
		fld := structRval.Field(i)
		var (
			fieldName  = fld.Name
			fieldType  = fld.Type
			fieldTag   = fld.Tag
			fieldValue = rVal.Field(i)
		)
		tagValue, ok := fieldTag.Lookup(validatorKey)
		if !ok {
			continue
		}
		if tagValue == "" {
			continue
		}
		extrValMap, err := extractValidators(tagValue)
		if err != nil {
			return err
		}
		for key, limit := range extrValMap {
			validationFn, ok := validationFuncMap[key]
			if !ok {
				log.Fatalf("validator %s not implemented", key)
			}
			err = validationFn(key, limit, fieldName, fieldValue, fieldType)
			if err != nil {
				if ok = errors.As(err, &valerr); !ok {
					return err
				}
				errs = append(errs, valerr)
			}
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func extractValidators(val string) (map[string]string, error) {
	if !strings.Contains(val, ":") {
		return nil, ErrorValidator{""}
	}
	validateCandidates := strings.Split(val, "|")
	extractedMap := make(map[string]string, len(validateCandidates))
	for _, candidate := range validateCandidates {
		keyWithVal := strings.SplitN(candidate, ":", 2)
		if len(keyWithVal) != 2 {
			return nil, ErrorValidator{message: fmt.Sprintf("validator should be func:limit but got %s", val)}
		}
		if _, ok := extractedMap[strings.Trim(keyWithVal[0], " ")]; ok {
			return nil, ErrorValidator{message: fmt.Sprintf(
				"duplicate key \"%s\" in the same validator: %s", strings.Trim(keyWithVal[0], " "), val)}
		}
		extractedMap[strings.Trim(keyWithVal[0], " ")] = keyWithVal[1]
	}
	return extractedMap, nil
}
