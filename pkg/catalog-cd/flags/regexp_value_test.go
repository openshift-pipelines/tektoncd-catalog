package flags

import (
	"regexp"
	"testing"

	"github.com/onsi/gomega"
)

func TestNewRegexpValue(t *testing.T) {
	g := gomega.NewWithT(t)

	// slice where the regular expressions are stored, informed as a pointer to the RegexpValue
	// instance, by calling Set() successfully the regexp is appended
	values := []regexp.Regexp{}
	r := NewRegexpValue(&values)

	tests := []struct {
		name    string
		raw     string
		wantErr bool
	}{{
		name:    "invalid expression",
		raw:     ")))(((",
		wantErr: true,
	}, {
		name:    "valid expression",
		raw:     "^$",
		wantErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.Set(tt.raw); (err != nil) != tt.wantErr {
				t.Errorf("RegexpValue.Set() error = %v, wantErr %v", err, tt.wantErr)
			}

			// when Set is successful asserting the corresponding expresion is appended on the
			// original slice instance
			if !tt.wantErr {
				g.Expect(len(values)).NotTo(gomega.BeZero())
				g.Expect(values[len(values)-1].String()).To(gomega.Equal(tt.raw))
			}
		})
	}
}
