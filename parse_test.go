package formulae

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"1+2*3", "('1' + ('2' * '3'))"},
		{"sin 5", "(sin '5')"},
		{"-5", "(- '5')"},
		{"-5--5", "((- '5') - (- '5'))"},
		{"5sqrt 2^2*3+1", "(('5' * (sqrt (('2' ^ '2') * '3'))) + '1')"},
		{"sin sqrt 2", "(sin (sqrt '2'))"},
		{"Sin 4", "(sin '4')"},
		{"1^2^3", "('1' ^ ('2' ^ '3'))"},
		{"4x", "('4' * 'x')"},
		{"5.5.5", "('5.5' * '0.5')"},
		{"5.5e-6.4e+4", "('5.5e-06' * '4000')"},
		{"(2+4)*3", "(('2' + '4') * '3')"},
		{"sin(x)^3", "((sin x)^3)"},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			formula, errs := Parse(test.in)
			if len(errs) > 0 {
				t.Fatal(formula, errs)
			}
			if formula.root.String() != test.out {
				t.Fatal(formula.root.String(), "!=", test.out)
			}
		})
	}
}

func TestParseErr(t *testing.T) {
	tests := []struct {
		in  string
		err string
	}{
		{"", "empty formula"},
		{"1++2", "2: operator has no operands"},
		{"4&4", "2: bad input"},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			_, errs := Parse(test.in)
			if len(errs) == 0 {
				t.Fatal("nil !=", test.err)
			}
			if errs[0].Error() != test.err {
				t.Fatal(errs[0].Error(), "!=", test.err)
			}
		})
	}
}
