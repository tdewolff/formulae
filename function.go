package formulae

import (
	"math"
    "fmt"
)

type Vars map[string]complex128

func (v Vars) Set(name string, value complex128) {
	v[name] = value
}

func (v Vars) Duplicate() Vars {
	v2 := make(map[string]complex128, len(v))
	for key, val := range v {
		v2[key] = val
	}
	return v2
}

var DefaultVars = Vars{
	"e":   complex(math.E, 0),
	"pi":  complex(math.Pi, 0),
	"phi": complex(math.Phi, 0),
}

////////////////

type Function struct {
	root Node
    Vars
    nthDerivative int
}

func (f *Function) String() string {
    return f.root.String()
}

func (f *Function) LaTeX() string {
    fs := "f"
    if f.nthDerivative == 1 {
        fs = fmt.Sprintf("\\frac{\\partial}{\\partial x}")
    } else if f.nthDerivative > 1 {
        fs = fmt.Sprintf("\\frac{\\partial^{%v}}{\\partial x^{%v}}", f.nthDerivative, f.nthDerivative)
    }
    return fmt.Sprintf("%s(x) = %s", fs, f.root.LaTeX())
}

func (f *Function) Optimize() {
    f.root = Optimize(f.root)
}

func (f *Function) Derivative() *Function {
    root := f.root.Derivative()
    root = Optimize(root)
    return &Function{
        root: root,
        Vars: f.Vars,
        nthDerivative: f.nthDerivative + 1,
    }
}

func (f *Function) Calc(x complex128) (complex128, error) {
    return f.root.Calc(x, f.Vars)
}

func (f *Function) Interval(xMin, xStep, xMax float64) ([]float64, []complex128, []error) {
	n := int((xMax-xMin)/xStep) + 1
	xs := make([]float64, n)
	ys := make([]complex128, n)

	x := xMin
    var err error
    var errs []error
    for i := 0; i < n; i++ {
        xs[i] = x
        ys[i], err = f.Calc(complex(x, 0))
        if err != nil {
            errs = append(errs, err)
        }
        x += xStep
    }
	return xs, ys, errs
}
