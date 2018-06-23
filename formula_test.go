package formulae

import (
	"math"
	"testing"
)

var Epsilon = 1e-6

func TestCalc(t *testing.T) {
	tests := []struct {
		in  string
		out float64
	}{
		{"1+2*3", 7},
		{"4x", 20},
		{"sin(pi)", 0},
		{"sin(pi/2)", 1},
		{"ln(e)", 1},
	}

	vars := Vars{
		"x": 5,
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			formula, errs := Parse(test.in)
			if len(errs) > 0 {
				t.Fatal(formula, errs)
			}

			f, err := formula.Calc(vars)
			if err != nil {
				t.Fatal(err)
			}
			if math.Abs(f-test.out) > Epsilon {
				t.Fatal(f, "!=", test.out)
			}
		})
	}
}

func TestCalcErr(t *testing.T) {
	tests := []struct {
		in  string
		err string
	}{
		{"4y", "2: undeclared variable 'y'"},
		{"3/(5-x)", "2: division by zero"},
	}

	vars := Vars{
		"x": 5,
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			formula, errs := Parse(test.in)
			if len(errs) > 0 {
				t.Fatal(formula, errs)
			}

			_, err := formula.Calc(vars)
			if err.Error() != test.err {
				t.Fatal(err.Error(), "!=", test.err)
			}
		})
	}
}
