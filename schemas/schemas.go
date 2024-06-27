package schemas

import (
	"fmt"
	"net/mail"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type FieldError interface {
	Error() string
	Name() string
	Value() string
	Kind() reflect.Kind
	Type() reflect.Type
}

// Compile time errors
var (
	_ FieldError = new(fieldError)
	_ error      = new(fieldError)
)

type ConstRuleValues struct {
	IDMinValue              int16
	UsernameMinLength       uint16
	UsernameMaxLength       uint16
	PasswordMinLength       uint16
	PasswordMaxLength       uint16
	CodeLength              uint16
	CodeRegex               string
	ProductTitleMinLength   uint16
	ProductTitleMaxLength   uint16
	ProductCaloriesMinValue int16
	ProductCaloriesMaxValue int16
	ProductNutrientMinValue int16
	ProductNutrientMaxValue int16
}

var DefRV ConstRuleValues = ConstRuleValues{
	IDMinValue:              1,
	UsernameMinLength:       2,
	UsernameMaxLength:       12,
	PasswordMinLength:       6,
	PasswordMaxLength:       30,
	CodeLength:              6,
	CodeRegex:               "^[A-Z0-9]+$",
	ProductTitleMinLength:   4,
	ProductTitleMaxLength:   128,
	ProductCaloriesMinValue: 0,
	ProductCaloriesMaxValue: 1000,
	ProductNutrientMinValue: 0,
	ProductNutrientMaxValue: 100,
}

type RulesMap map[string]func(field reflect.Value, structField reflect.StructField, v string) error

var Rules RulesMap = RulesMap{
	"ge":         geF,
	"le":         leF,
	"max_length": maxLengthF,
	"min_length": minLengthF,
	"regex":      regexF,
	"email":      emailF,
}

var Formats map[string]string = map[string]string{
	"id":    fmt.Sprintf("ge=%d", DefRV.IDMinValue),
	"email": "email",
	"username": fmt.Sprintf("min_length=%d,max_length=%d",
		DefRV.UsernameMinLength, DefRV.UsernameMaxLength),
	"password": fmt.Sprintf("min_length=%d,max_length=%d",
		DefRV.PasswordMinLength, DefRV.PasswordMaxLength),
	"code": fmt.Sprintf("min_length=%d,max_length=%d,regex=%s",
		DefRV.CodeLength, DefRV.CodeLength, DefRV.CodeRegex),
	"product_title": fmt.Sprintf("min_length=%d,max_length=%d",
		DefRV.ProductTitleMinLength, DefRV.ProductTitleMaxLength),
	"product_calories": fmt.Sprintf("ge=%d,le=%d",
		DefRV.ProductCaloriesMinValue, DefRV.ProductCaloriesMaxValue),
	"product_nutrient": fmt.Sprintf("ge=%d,le=%d",
		DefRV.ProductNutrientMinValue, DefRV.ProductNutrientMaxValue),
}

func emailF(field reflect.Value, structField reflect.StructField, v string) error {
	fieldKind := field.Kind()

	switch fieldKind {
	case reflect.String:
		_, err := mail.ParseAddress(field.String())
		if err != nil {
			return err
		}
	default:
		err := fmt.Errorf("Mismatched type and value at validation tag. Field name - %s. Expected type - %s",
			structField.Name, "String")
		panic(err.Error())
	}

	return nil
}

func regexF(field reflect.Value, structField reflect.StructField, v string) error {
	fieldKind := field.Kind()

	switch fieldKind {
	case reflect.String:
		validRegex := regexp.MustCompile(v)
		if validRegex == nil {
			err := fmt.Errorf("Invalid regex prompt at validation tag. Field name - %s. Regex - %s",
				structField.Name, v)
			panic(err.Error())
		}
		if !validRegex.MatchString(field.String()) {
			return fmt.Errorf("%s does not match regex %s", structField.Name, v)
		}

	default:
		err := fmt.Errorf("Mismatched type and value at validation tag. Field name - %s. Real - %s, expected - %s",
			structField.Name, fieldKind.String(), "String")
		panic(err.Error())
	}

	return nil
}

func minLengthF(field reflect.Value, structField reflect.StructField, v string) error {
	fieldKind := field.Kind()

	switch fieldKind {
	case reflect.String:
		cmp, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			panic(err.Error())
		}
		if uint64(len(field.String())) < cmp {
			return fmt.Errorf("%s length is less than %d", structField.Name, cmp)
		}
	default:
		err := fmt.Errorf("Mismatched type and value at validation tag. Field name - %s. Real - %s, expected - %s",
			structField.Name, fieldKind.String(), "String")
		panic(err.Error())
	}

	return nil
}

func maxLengthF(field reflect.Value, structField reflect.StructField, v string) error {
	fieldKind := field.Kind()

	switch fieldKind {
	case reflect.String:
		cmp, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			panic(err.Error())
		}
		if uint64(len(field.String())) > cmp {
			return fmt.Errorf("%s exceedes max length of %d", structField.Name, cmp)
		}
	default:
		err := fmt.Errorf("Mismatched type and value at validation tag. Field name - %s. Real - %s, expected - %s",
			structField.Name, fieldKind.String(), "String")
		panic(err.Error())
	}

	return nil
}

func leF(field reflect.Value, structField reflect.StructField, v string) error {
	fieldKind := field.Kind()

	switch fieldKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		cmp, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic(err.Error())
		}
		if field.Int() > cmp {
			return fmt.Errorf("%s is greater than %d", structField.Name, cmp)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		cmp, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			panic(err.Error())
		}
		if field.Uint() > cmp {
			return fmt.Errorf("%s is greater than %d", structField.Name, cmp)
		}
	default:
		err := fmt.Errorf("Mismatched type and value at validation tag. Field name - %s. Real - %s, expected - %s",
			structField.Name, fieldKind.String(), "Integer")
		panic(err.Error())
	}

	return nil
}

func geF(field reflect.Value, structField reflect.StructField, v string) error {
	fieldKind := field.Kind()

	switch fieldKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		cmp, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic(err.Error())
		}
		if field.Int() < cmp {
			return fmt.Errorf("%s is less than %d", structField.Name, cmp)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		cmp, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			panic(err.Error())
		}
		if field.Uint() < cmp {
			return fmt.Errorf("%s is less than %d", structField.Name, cmp)
		}
	default:
		err := fmt.Errorf("Mismatched type and value at validation tag. Field name - %s. Real - %s, expected - %s",
			structField.Name, fieldKind.String(), "Integer")
		panic(err.Error())
	}

	return nil
}

type fieldError struct {
	err   error
	name  string
	value string
	kind  reflect.Kind
	typ   reflect.Type
}

func (fe *fieldError) Error() string {
	return fe.err.Error()
}

func (fe *fieldError) Name() string {
	return fe.name
}

func (fe *fieldError) Value() string {
	return fe.value
}

func (fe *fieldError) Kind() reflect.Kind {
	return fe.kind
}

func (fe *fieldError) Type() reflect.Type {
	return fe.typ
}

type ValidationErrors []FieldError

func (ve ValidationErrors) Error() string {
	out := ""
	for _, err := range ve {
		out += err.Error() + "; "
	}
	return out
}

func validate(field reflect.Value, structField reflect.StructField, validateStr string) FieldError {
	var fe FieldError = nil
	var omitzero bool = false
	rules := strings.Split(validateStr, ",")

	for _, rule := range rules {
		var key, value string = "", ""
		parts := strings.Split(rule, "=")
		partsLen := len(parts)
		switch partsLen {
		case 1:
			key = parts[0]
		case 2:
			key = parts[0]
			value = parts[1]
		default:
			err := fmt.Errorf("Error at tag syntax. Field name: \"%s\". Validate string: \"%s\".",
				structField.Name, validateStr)
			panic(err.Error())
		}

		if key == "omitzero" {
			omitzero = true
			if fe != nil {
				break
			}
			continue
		}
		if fe != nil {
			continue
		}

		ruleFunction, exists := Rules[key]
		if !exists {
			err := fmt.Errorf("Error at tag value. Field name: \"%s\". Validate string: \"%s\".",
				structField.Name, validateStr)
			panic(err.Error())
		}

		err := ruleFunction(field, structField, value)
		if err != nil {
			fe = &fieldError{
				err:   err,
				name:  structField.Name,
				value: field.String(),
				kind:  field.Kind(),
				typ:   field.Type(),
			}
		}

	}

	if omitzero && field.IsZero() {
		return nil
	}
	return fe
}

func ValidateStruct(s interface{}) ValidationErrors {
	val := reflect.Indirect(reflect.ValueOf(s))

	ve := ValidationErrors{}
	for i := 0; i < val.NumField(); i++ {
		var tag string = ""
		var validateStr string = ""
		field := val.Field(i)
		structField := val.Type().Field(i)

		tag = structField.Tag.Get("format")
		if tag != "" {
			format, exists := Formats[tag]
			if !exists {
				err := fmt.Errorf("Unknown format type. Field name: \"%s\". Given type: \"%s\".",
					structField.Name, tag)
				panic(err.Error())
			}
			validateStr = format
		}

		tag = structField.Tag.Get("validate")
		if tag != "" {
			if validateStr == "" {
				validateStr = tag
			} else {
				validateStr += "," + tag
			}
		}

		if validateStr == "" {
			continue
		}
		fe := validate(field, structField, validateStr)
		if fe != nil {
			// ve.Add(fe)
			ve = append(ve, fe)
		}
	}
	if len(ve) != 0 {
		return ve
	}
	return nil
}

func IsZero(v any) bool {
	vr := reflect.ValueOf(v)
	return vr.IsZero()
}
