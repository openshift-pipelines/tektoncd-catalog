package linter

import (
	"testing"
)

func Test_lowercaseSliceMapLinter(t *testing.T) {
	tests := []struct {
		name    string
		slice   []interface{}
		wantErr bool
	}{{
		name:    "empty slice",
		slice:   []interface{}{},
		wantErr: false,
	}, {
		name: "single valid entry",
		slice: []interface{}{map[string]interface{}{
			"name":        "lowercase",
			"description": "description",
		}},
		wantErr: false,
	}, {
		name: "single invalid entry",
		slice: []interface{}{map[string]interface{}{
			"name":        "UPPERCASE",
			"description": "description",
		}},
		wantErr: true,
	}, {
		name: "single invalid entry (no description)",
		slice: []interface{}{map[string]interface{}{
			"name":        "lowercase",
			"description": "",
		}},
		wantErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := lowercaseSliceMapLinter(tt.slice); (err != nil) != tt.wantErr {
				t.Errorf("lowercaseSliceMapLinter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_uppercaseSliceMapLinter(t *testing.T) {
	tests := []struct {
		name    string
		slice   []interface{}
		wantErr bool
	}{{
		name:    "empty slice",
		slice:   []interface{}{},
		wantErr: false,
	}, {
		name: "single valid entry",
		slice: []interface{}{map[string]interface{}{
			"name":        "UPPERCASE",
			"description": "description",
		}},
		wantErr: false,
	}, {
		name: "single invalid entry",
		slice: []interface{}{map[string]interface{}{
			"name":        "lowercase",
			"description": "description",
		}},
		wantErr: true,
	}, {
		name: "single invalid entry (no description)",
		slice: []interface{}{map[string]interface{}{
			"name":        "UPPERCASE",
			"description": "",
		}},
		wantErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := uppercaseSliceMapLinter(tt.slice); (err != nil) != tt.wantErr {
				t.Errorf("uppercaseSliceMapLinter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
