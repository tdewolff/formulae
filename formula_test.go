package formulae

import (
	"math/cmplx"
	"testing"
)

var Epsilon = 1e-6

func TestCalc(t *testing.T) {
	tests := []struct {
		in  string
		out complex128
	}{
		{"1+2*3", 7},
		{"4x", 20},
		{"sin(pi)", 0},
		{"sin(pi/2)", 1},
		{"ln(e)", 1},
		{"LN(E)", 1},
		{"5+1i+6+2i", 11 + 3i},
		{"i", 1i},
	}

	x := 5 + 0i
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			formula, errs := Parse(test.in)
			if len(errs) > 0 {
				t.Fatal(formula, errs)
			}

			err := formula.Compile(DefaultVars)
			if err != nil {
				t.Fatal(err)
			}

			y := formula.Calc(x)
			if cmplx.Abs(y-test.out) > Epsilon {
				t.Fatal(y, "!=", test.out)
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
		// {"3/(5-x)", "2: division by zero"},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			formula, errs := Parse(test.in)
			if len(errs) > 0 {
				t.Fatal(formula, errs)
			}

			err := formula.Compile(DefaultVars)
			if err.Error() != test.err {
				t.Fatal(err.Error(), "!=", test.err)
			}
		})
	}
}

var formula, _ = Parse("sin(x)^2+1/x+0.001x^3")
var x = 5 + 1i
var y = 0 + 0i

func BenchmarkCalcNativeGo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		y = cmplx.Pow(cmplx.Sin(x), 2) + (1 / x) + 0.001*cmplx.Pow(x, 3)
	}
}

func BenchmarkCalcGo(b *testing.B) {
	_ = formula.Compile(DefaultVars)
	for i := 0; i < b.N; i++ {
		y = formula.Calc(x)
	}
}

func BenchmarkCalcLua(b *testing.B) {
	compiled, _ := formula.CompileLua()
	defer compiled.Close()

	for i := 0; i < b.N; i++ {
		y, _ = compiled.Calc(x)
	}
}
