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
			function, errs := Parse(test.in)
			if len(errs) > 0 {
				t.Fatal(function, errs)
			}

			y, err := function.Calc(x)
			if err != nil {
				t.Fatal(err)
			}
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
		{"4y", "undefined variable 'y'"},
		{"3/(5-x)", "division by zero"},
	}

	x := 5 + 0i
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			function, errs := Parse(test.in)
			if len(errs) > 0 {
				t.Fatal(function, errs)
			}

			_, err := function.Calc(x)
            if err == nil {
				t.Fatal("nil !=", test.err)
            } else if err.Error() != test.err {
				t.Fatal(err.Error(), "!=", test.err)
			}
		})
	}
}

var function, _ = Parse("sin(x)^2+1/x+0.001x^3")
var N = 100
var xs []complex128
var ys []complex128

func init() {
    xs = make([]complex128, N)
    ys = make([]complex128, N)
    for j := 0; j < N; j++ {
        xs[j] = complex(float64(j)+1.0, 0)
    }
}

func BenchmarkCalcNative(b *testing.B) {
	for i := 0; i < b.N; i++ {
        for j := 0; j < N; j++ {
		    ys[j] = cmplx.Pow(cmplx.Sin(xs[j]), 2) + (1 / xs[j]) + 0.001*cmplx.Pow(xs[j], 3)
        }
	}
}

func BenchmarkCalc(b *testing.B) {
    for i := 0; i < b.N; i++ {
        for i, x := range xs {
            ys[i], _ = function.Calc(x)
        }
    }
}
