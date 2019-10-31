package validation

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/System-Glitch/goyave/helpers"
)

// Rule function defining a validation rule.
// Passing rules should return true, false otherwise.
//
// Rules can modifiy the validated value if needed.
// For example, the "numeric" rule converts the data to float64 if it's a string.
type Rule func(string, interface{}, []string, map[string]interface{}) bool

func validateRequired(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	val, ok := form[field]
	if ok {
		if val == nil {
			return false
		}
		if str, okStr := val.(string); okStr && str == "" {
			return false
		}
	}
	return ok
}

// -------------------------
// Generic

func validateMin(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	min, err := strconv.ParseFloat(parameters[0], 64)
	if err != nil {
		panic(err)
	}
	switch getFieldType(value) {
	case "numeric":
		floatValue, _ := helpers.ToFloat64(value)
		return floatValue >= min
	case "string":
		return len(value.(string)) >= int(min)
	case "array":
		list := reflect.ValueOf(value)
		return list.Len() >= int(min)
	case "file":
		return false // TODO implement file min size
	}

	return true // Pass if field type cannot be checked (bool, dates, ...)
}

func validateMax(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	max, err := strconv.ParseFloat(parameters[0], 64)
	if err != nil {
		panic(err)
	}
	switch getFieldType(value) {
	case "numeric":
		floatValue, _ := helpers.ToFloat64(value)
		return floatValue <= max
	case "string":
		return len(value.(string)) <= int(max)
	case "array":
		list := reflect.ValueOf(value)
		return list.Len() <= int(max)
	case "file":
		return false // TODO implement file max size
	}

	return true // Pass if field type cannot be checked (bool, dates, ...)
}

func validateBetween(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	min, errMin := strconv.ParseFloat(parameters[0], 64)
	max, errMax := strconv.ParseFloat(parameters[1], 64)
	if errMin != nil {
		panic(errMin)
	}
	if errMax != nil {
		panic(errMax)
	}

	switch getFieldType(value) {
	case "numeric":
		floatValue, _ := helpers.ToFloat64(value)
		return floatValue >= min && floatValue <= max
	case "string":
		length := len(value.(string))
		return length >= int(min) && length <= int(max)
	case "array":
		list := reflect.ValueOf(value)
		length := list.Len()
		return length >= int(min) && length <= int(max)
	case "file":
		return false // TODO implement file between size
	}

	return true // Pass if field type cannot be checked (bool, dates, ...)
}

func validateGreaterThan(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	valueType := getFieldType(value)

	compared, exists := form[parameters[0]]
	if !exists || valueType != getFieldType(compared) {
		return false // Can't compare two different types or missing field
	}

	switch valueType {
	case "numeric":
		floatValue, _ := helpers.ToFloat64(value)
		comparedFloatValue, _ := helpers.ToFloat64(compared)
		return floatValue > comparedFloatValue
	case "string":
		return len(value.(string)) > len(compared.(string))
	case "array":
		return reflect.ValueOf(value).Len() > reflect.ValueOf(compared).Len()
	case "file":
		return false // TODO implement file greater than size
	}

	return false
}

func validateGreaterThanEqual(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	valueType := getFieldType(value)

	compared, exists := form[parameters[0]]
	if !exists || valueType != getFieldType(compared) {
		return false // Can't compare two different types or missing field
	}

	switch valueType {
	case "numeric":
		floatValue, _ := helpers.ToFloat64(value)
		comparedFloatValue, _ := helpers.ToFloat64(compared)
		return floatValue >= comparedFloatValue
	case "string":
		return len(value.(string)) >= len(compared.(string))
	case "array":
		return reflect.ValueOf(value).Len() >= reflect.ValueOf(compared).Len()
	case "file":
		return false // TODO implement file greater than size
	}

	return false
}

func validateLowerThan(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	valueType := getFieldType(value)

	compared, exists := form[parameters[0]]
	if !exists || valueType != getFieldType(compared) {
		return false // Can't compare two different types or missing field
	}

	switch valueType {
	case "numeric":
		floatValue, _ := helpers.ToFloat64(value)
		comparedFloatValue, _ := helpers.ToFloat64(compared)
		return floatValue < comparedFloatValue
	case "string":
		return len(value.(string)) < len(compared.(string))
	case "array":
		return reflect.ValueOf(value).Len() < reflect.ValueOf(compared).Len()
	case "file":
		return false // TODO implement file greater than size
	}

	return false
}

func validateLowerThanEqual(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	valueType := getFieldType(value)

	compared, exists := form[parameters[0]]
	if !exists || valueType != getFieldType(compared) {
		return false // Can't compare two different types or missing field
	}

	switch valueType {
	case "numeric":
		floatValue, _ := helpers.ToFloat64(value)
		comparedFloatValue, _ := helpers.ToFloat64(compared)
		return floatValue <= comparedFloatValue
	case "string":
		return len(value.(string)) <= len(compared.(string))
	case "array":
		return reflect.ValueOf(value).Len() <= reflect.ValueOf(compared).Len()
	case "file":
		return false // TODO implement file greater than size
	}

	return false
}

// -------------------------
// Numeric

func validateNumeric(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	rv := reflect.ValueOf(value)
	kind := rv.Kind().String()
	switch {
	case strings.HasPrefix(kind, "int"), strings.HasPrefix(kind, "uint") && kind != "uintptr", strings.HasPrefix(kind, "float"):
		return true
	case kind == "string":
		floatVal, err := strconv.ParseFloat(value.(string), 64)
		if err == nil {
			form[field] = floatVal
		}
		return err == nil
	default:
		return false
	}
}

func validateInteger(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	rv := reflect.ValueOf(value)
	kind := rv.Kind().String()
	switch {
	case strings.HasPrefix(kind, "int"), strings.HasPrefix(kind, "uint") && kind != "uintptr":
		return true
	case strings.HasPrefix(kind, "float"):
		if kind == "float64" {
			val, _ := value.(float64)
			if val-float64(int(val)) > 0 {
				return false
			}
			form[field] = int(val)
			return true
		}

		val, _ := value.(float32)
		if val-float32(int(val)) > 0 {
			return false
		}
		form[field] = int(val)
		return true
	case kind == "string":
		intVal, err := strconv.Atoi(value.(string))
		if err == nil {
			form[field] = intVal
		}
		return err == nil
	default:
		return false
	}
}

// -------------------------
// Strings

func validateString(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	_, ok := value.(string)
	return ok
}

// digits
// alpha
// alpha_dash
// alphanumeric
// email
// starts_with
// ends_with

// length
// min
// max

// ip address
// json
// regex
// timezone
// url
// uuid

// -------------------------
// Arrays

func validateArray(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	return getFieldType(value) == "array"
}

func validateDistinct(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	if getFieldType(value) != "array" {
		return false // Can't validate if not an array
	}

	found := []interface{}{}
	list := reflect.ValueOf(value)
	for i := 0; i < list.Len(); i++ {
		v := list.Index(i).Interface()
		if helpers.Contains(found, v) {
			return false
		}
		found = append(found, v)
	}

	return true
}

// -------------------------
// Dates

// date:format

// date before + before or equal
// date after + after or equal
// between + between or equal
// date equals
//

// -------------------------
// Files

// file
// mime
// image
// size
// extension

// -------------------------
// Misc

// boolean (accept 1, "on", "true", "yes")
// different
// in + not_in
// in_array + not_in_array
// nullable
// confirmed
