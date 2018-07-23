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

func TestCalcN(t *testing.T) {
    function, errs := Parse("x+1")
    if len(errs) > 0 {
        t.Fatal(function, errs)
    }

    N := 100
    xs := make([]complex128, N)
    for i := range xs {
        xs[i] = complex(float64(i), 0)
    }

    ys, err := function.CalcN(xs)
    if err != nil {
        t.Fatal(err)
    }

    y := ys[len(ys)-1]
    if real(y) != 100.0 {
	    t.Fatal(y, "!= 100.0")
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

    x := 5+0i
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			function, errs := Parse(test.in)
			if len(errs) > 0 {
				t.Fatal(function, errs)
			}

			_, err := function.Calc(x)
			if err.Error() != test.err {
				t.Fatal(err.Error(), "!=", test.err)
			}
		})
	}
}

var function, _ = Parse("sin(x)^2+1/x+0.001x^3")
var N = 100
var x []complex128
var y []complex128

func init() {
    x = make([]complex128, N)
    y = make([]complex128, N)
    for j := 0; j < N; j++ {
        x[j] = complex(float64(j)+1.0, 0)
    }
}

func BenchmarkCalcNative(b *testing.B) {
	for i := 0; i < b.N; i++ {
        for j := 0; j < N; j++ {
		    y[j] = cmplx.Pow(cmplx.Sin(x[j]), 2) + (1 / x[j]) + 0.001*cmplx.Pow(x[j], 3)
        }
	}
}

func BenchmarkCalc(b *testing.B) {
	for i := 0; i < b.N; i++ {
        for j := 0; j < N; j++ {
		    y[j], _ = function.Calc(x[j])
        }
	}
}

func BenchmarkCalcN(b *testing.B) {
    for i := 0; i < b.N; i++ {
        y, _ = function.CalcN(x)
    }
}
