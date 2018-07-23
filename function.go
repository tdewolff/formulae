package formulae

import (
	"math"
	"strings"
)

type Vars map[string]complex128

func (v Vars) Set(name string, value complex128) {
	v[strings.ToLower(name)] = value
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
    vars Vars
    buf []complex128
}

func (f *Function) Calc(xs []complex128) ([]complex128, error) {
    ys := make([]complex128, len(xs))
    copy(ys, xs)

    if cap(f.buf) < len(xs) {
        f.buf = make([]complex128, len(xs))
    }
    f.buf = f.buf[:len(xs)]

    return f.root.Calc(xs, ys)
}

func (f *Function) Interval(xMin, xStep, xMax float64) ([]float64, []complex128, error) {
	n := int((xMax-xMin)/xStep) + 1
	xsReal := make([]float64, n)
	xs := make([]complex128, n)

	x := xMin
    for i := 0; i < n; i++ {
        xsReal[i] = x
        xs[i] = complex(x, 0)
        x += xStep
    }

    ys, err := f.Calc(xs)
    if err != nil {
        return nil, nil, err
    }
	return xsReal, ys, nil
}
