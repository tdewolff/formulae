package formulae

import (
	"fmt"
	"math"
)

type Vars map[string]complex128

func (v Vars) Set(name string, f complex128) {
	v[name] = f
}

var DefaultVars = Vars{
	"e":   complex(math.E, 0),
	"pi":  complex(math.Pi, 0),
	"phi": complex(math.Phi, 0),
}

////////////////

type Formula struct {
	root     Node
	varNames map[string]bool
}

func (f *Formula) Calc(vars Vars) (complex128, error) {
	for name, val := range DefaultVars {
		if _, ok := vars[name]; !ok {
			vars[name] = val
		}
	}
	return f.root.Calc(vars)
}

func (f *Formula) Interval(xMin, xStep, xMax float64) ([]float64, []complex128, []error) {
	n := int((xMax-xMin)/xStep) + 1
	xs := make([]float64, n)
	ys := make([]complex128, n)

	vars := Vars{
		"x": 0,
	}

	x := xMin
	var errs []error
	for i := 0; i < n; i++ {
		vars["x"] = complex(x, 0)
		y, err := f.Calc(vars)
		if err != nil {
			errs = append(errs, fmt.Errorf("%v (x = %v)", err, x))
			continue
		}

		xs[i] = x
		ys[i] = y

		x += xStep
	}
	return xs, ys, errs
}
