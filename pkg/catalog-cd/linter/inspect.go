package linter

import (
	"errors"
	"fmt"
	"unicode"
)

var (
	ErrInvalidUppercase  = errors.New("uppercase name")
	ErrInvalidLowercase  = errors.New("lowercase name")
	ErrRequiredAttribute = errors.New("required attribute is not set")
)

// getAttribute gets a attribute value from the set, if the entry does not exist or is empty it
// returns error.
func getAttribute(set map[string]interface{}, entry string) (string, error) {
	value, ok := set[entry]
	if !ok {
		return "", fmt.Errorf("%w: %q", ErrRequiredAttribute, entry)
	}
	if value == "" {
		return "", fmt.Errorf("%w: %q is empty", ErrRequiredAttribute, entry)
	}
	return value.(string), nil
}

// isNameValidFn checks if the name attribute is valid.
type isNameValidFn func(string) error

// isSetValid inspect the set (map) informed making sure it contains a description, and the name
// matches the informed validation function.
func isSetValid(set map[string]interface{}, fn isNameValidFn) error {
	name, err := getAttribute(set, "name")
	if err != nil {
		return err
	}
	if err = fn(name); err != nil {
		return err
	}
	_, err = getAttribute(set, "description")
	return err
}

// isUppercaseNameValid asserts the informed name is all uppercase with underscore ("_") separator.
func isUppercaseNameValid(name string) error {
	for _, c := range name {
		if c != '_' && !unicode.IsUpper(c) {
			return fmt.Errorf("%w: %q contains lowercase %q", ErrInvalidUppercase, name, c)
		}
	}
	return nil
}

// isLowercaseNameValid asserts the informed name is all lowercase with dash ("-") separator.
func isLowercaseNameValid(name string) error {
	for _, c := range name {
		if c != '-' && !unicode.IsLower(c) {
			return fmt.Errorf("%w: %q contains uppercase %q", ErrInvalidLowercase, name, c)
		}
	}
	return nil
}

// isLowercaseSetValid  a set (map) of attributes to assert if contains a description and lowercase
// name.
func isLowercaseSetValid(set map[string]interface{}) error {
	return isSetValid(set, isLowercaseNameValid)
}

// isUppercaseSetValid lints a set (map) of attributes to assert if contains a description and
// uppercase name.
func isUppercaseSetValid(set map[string]interface{}) error {
	return isSetValid(set, isUppercaseNameValid)
}

type isSetValidFn func(map[string]interface{}) error

// loopSliceMap applies the informed linter function for each slice entry.
func loopSliceMap(slice []interface{}, fn isSetValidFn) error {
	for i, entry := range slice {
		var set map[string]interface{}
		set, ok := entry.(map[string]interface{})
		if !ok {
			return fmt.Errorf("entry %d is invalid '%#v'", i, set)
		}
		if err := fn(set); err != nil {
			return fmt.Errorf("entry %d %w", i, err)
		}
	}
	return nil
}

// lowercaseSliceMapLinter lints each slice entry as a map of attributes, each map entry must contain
// a description and lowercase name.
func lowercaseSliceMapLinter(slice []interface{}) error {
	return loopSliceMap(slice, isLowercaseSetValid)
}

// uppercaseSliceMapLinter lints each slice entry as a map of attributes, each map entry must contain
// a description and uppercase name.
func uppercaseSliceMapLinter(slice []interface{}) error {
	return loopSliceMap(slice, isUppercaseSetValid)
}
