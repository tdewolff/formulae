package formulae

import (
	"strings"
	"testing"
)

func TestNumeric(t *testing.T) {
	tests := []struct {
		in    string
		valid bool
	}{
		{"1", true},
		{"1.0", true},
		{".0", true},
		{"1.", false},

		{"1e+1", true},
		{"1e-1", true},
		{"1E5", true},
		{"1.5e5", true},
		{"1.5e5.5", false},

		{"i", true},
		{"1i", true},
		{"1.0i", true},
		{"1.0e5i", true},

		{"1+5i", true},
		{"1 + 5i", true},
		{"1.3e-3 + 4.5e-7i", true},
		{"1 + i", false},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			l := NewLexer(strings.NewReader(test.in))
			tt, data := l.Next()

			if tt != NumericToken {
				t.Fatal("not numeric")
			}
			valid := len(data) == len(test.in)
			if valid != test.valid {
				if valid {
					t.Fatal(string(data), "==", test.in)
				} else {
					t.Fatal(string(data), "!=", test.in)
				}
			}
		})
	}
}
