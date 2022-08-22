package apicontext

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

func NewValidator() *validator.Validate {
	validate := validator.New()

	// field is required is  (Field1 == A | B) OR (Field2 == C | D)
	// required_if=Field1:A,B
	if err := validate.RegisterValidation(`required_if`, func(fl validator.FieldLevel) bool {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()
		// fields := make(map[string][]string)
		param := strings.Split(fl.Param(), `:`)
		paramField := param[0]
		paramValues := strings.Split(param[1], `'`)
		if paramField == "" {
			return true
		}

		// param field reflect.Value.
		var paramFieldValue reflect.Value

		if fl.Parent().Kind() == reflect.Ptr {
			paramFieldValue = fl.Parent().Elem().FieldByName(paramField)
		} else {
			paramFieldValue = fl.Parent().FieldByName(paramField)
		}

		isHope := false
		for _, value := range paramValues {
			if isEq(paramFieldValue, value) {
				isHope = true
			}
		}
		if !isHope {
			return true
		}

		return hasValue(fl)

	}); err != nil {
		logrus.Errorf("custom validator err: [%s]", err.Error())
	}

	return validate
}

// The following functions are copied from validator.v9 lib.

func hasValue(fl validator.FieldLevel) bool {
	return requireCheckFieldKind(fl, "")
}

func requireCheckFieldKind(fl validator.FieldLevel, param string) bool {
	field := fl.Field()
	if len(param) > 0 {
		if fl.Parent().Kind() == reflect.Ptr {
			field = fl.Parent().Elem().FieldByName(param)
		} else {
			field = fl.Parent().FieldByName(param)
		}
	}
	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return !field.IsNil()
	default:
		_, _, nullable := fl.ExtractType(field)
		if nullable && field.Interface() != nil {
			return true
		}
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}

func isEq(field reflect.Value, value string) bool {
	kind := field.Kind()
	if kind == reflect.Ptr {
		kind = field.Type().Elem().Kind()
		field = field.Elem()
	}
	switch kind {

	case reflect.String:
		return field.String() == value

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(value)

		return int64(field.Len()) == p

	case reflect.Bool:
		return strconv.FormatBool(field.Bool()) == value

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(value)

		return field.Int() == p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(value)

		return field.Uint() == p

	case reflect.Float32, reflect.Float64:
		p := asFloat(value)

		return field.Float() == p
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

func asInt(param string) int64 {

	i, err := strconv.ParseInt(param, 0, 64)
	panicIf(err)

	return i
}

func asUint(param string) uint64 {

	i, err := strconv.ParseUint(param, 0, 64)
	panicIf(err)

	return i
}

func asFloat(param string) float64 {

	i, err := strconv.ParseFloat(param, 64)
	panicIf(err)

	return i
}

func panicIf(err error) {
	if err != nil {
		panic(err.Error())
	}
}
