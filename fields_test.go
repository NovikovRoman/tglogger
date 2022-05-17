package tglogger

import (
	"regexp"
	"testing"
)

func TestFields_String(t *testing.T) {
	pointerInt := 123

	tests := []struct {
		name  string
		f     Fields
		wantS string
		re    *regexp.Regexp
	}{
		{
			name: "strings",
			f: Fields{
				"s1": "string1",
			},
			wantS: "s1: string1\n",
		},
		{
			name: "int",
			f: Fields{
				"int": 12,
			},
			wantS: "int: 12\n",
		},
		{
			name: "int8",
			f: Fields{
				"int8": int8(-12),
			},
			wantS: "int8: -12\n",
		},
		{
			name: "uint64",
			f: Fields{
				"uint64": uint64(18446744073709551615),
			},
			wantS: "uint64: 18446744073709551615\n",
		},
		{
			name: "pointer",
			f: Fields{
				"pointer": &pointerInt,
			},
			wantS: "",
			re:    regexp.MustCompile(`(!si)^pointer: 0x[0-9a-h]+\n$`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotS := tt.f.String()

			if tt.re == nil {
				if gotS != tt.wantS {
					t.Errorf("String() = %v, want %v", gotS, tt.wantS)
				}
				return
			}

			if tt.re.MatchString(gotS) {
				t.Errorf("String() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}
