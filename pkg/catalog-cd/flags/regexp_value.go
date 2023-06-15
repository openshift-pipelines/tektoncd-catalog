package flags

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/spf13/pflag"
)

// RegexpValue implements pflag.Value interface to manage an array of regular expressions, thus each
// command-line flag represents one entry on the array.
type RegexpValue struct {
	values *[]regexp.Regexp // array of regexp
}

var _ pflag.Value = &RegexpValue{}

// ErrInvalidRegexp the regular expression doesn't compile.
var ErrInvalidRegexp = errors.New("invalid regular expression")

// Type exposes the "type" to cobra's pflag.Value, a "stringArray" means each command-line flag
// becomes a single slice entry.
func (*RegexpValue) Type() string {
	return "stringArray"
}

// Set adds a regexp array entry, the raw expression must be compiled successfully.
func (r *RegexpValue) Set(raw string) error {
	value, err := regexp.Compile(raw)
	if err != nil {
		return fmt.Errorf("%w: '%v' %w", ErrInvalidRegexp, raw, err)
	}
	*r.values = append(*r.values, *value)
	return nil
}

// String shows the current array entries as string.
func (r *RegexpValue) String() string {
	slice := []string{}
	for _, value := range *r.values {
		slice = append(slice, value.String())
	}
	return fmt.Sprintf("%#v", slice)
}

// NewRegexpValue instantiates the RegexpValue with a pointer to the slice, the pointer receives the
// slice (array) entries.
func NewRegexpValue(values *[]regexp.Regexp) *RegexpValue {
	return &RegexpValue{values: values}
}
