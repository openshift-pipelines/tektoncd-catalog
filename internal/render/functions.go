package render

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"
)

// templateFuncMap map with the functions available on the template.
var templateFuncMap = template.FuncMap{
	"chomp":          chomp,
	"formatType":     formatType,
	"formatValue":    formatValue,
	"formatOptional": formatOptional,
}

// chomp removes new lines.
func chomp(s string) string {
	return strings.TrimSuffix(strings.ReplaceAll(s, "\n", " "), " ")
}

// formatType when type is not informed the type "string" returned.
func formatType(s interface{}) string {
	if s == nil || s.(string) == "" {
		return "string"
	}
	return s.(string)
}

func anySliceJoin(slice []interface{}, separator string) string {
	stringSlice := []string{}
	for _, j := range slice {
		stringSlice = append(stringSlice, j.(string))
	}
	return strings.Join(stringSlice, separator)
}

// formatValues highlights the informed value is required or empty.
func formatValue(value interface{}) string {
	if value == nil {
		return "(required)"
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		if v.String() == "" {
			return "\"\" (empty)"
		}
		return fmt.Sprintf("`%s`", v.String())
	case reflect.Slice:
		if v.Len() == 0 {
			return "`[]` (empty)"
		}
		return fmt.Sprintf("`[ %s ]`", anySliceJoin(v.Interface().([]interface{}), ", "))
	case reflect.Map:
		iter := v.MapRange()
		slice := []string{}
		for iter.Next() {
			slice = append(slice, fmt.Sprintf("%s=\"%s\"", iter.Key(), iter.Value()))
		}
		return fmt.Sprintf("`{ %s }`", strings.Join(slice, ", "))
	default:
		panic(fmt.Sprintf("unsupported param type %q", v.Kind()))
	}
}

// formatOptional makes sure "false" is printed when the informed variable is nil.
func formatOptional(s interface{}) string {
	if s == nil || !s.(bool) {
		return "false"
	}
	return "true"
}
