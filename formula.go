package formulae

import "math"

type Vars map[string]float64

func (v Vars) Set(name string, f float64) {
	v[name] = f
}

var DefaultVars = Vars{
	"e":   math.E,
	"E":   math.E,
	"pi":  math.Pi,
	"Pi":  math.Pi,
	"phi": math.Phi,
	"Phi": math.Phi,
}

////////////////

type Formula struct {
	root     Node
	varNames map[string]bool
}

func (f *Formula) Calc(vars Vars) (float64, error) {
	for name, val := range DefaultVars {
		if _, ok := vars[name]; !ok {
			vars[name] = val
		}
	}
	return f.root.Calc(vars)
}

func (f *Formula) Interval(xMin, xStep, xMax float64) ([]float64, []float64, error) {
	n := int((xMax-xMin)/xStep) + 1
	xs := make([]float64, n)
	ys := make([]float64, n)

	vars := Vars{
		"x": 0,
	}

	x := xMin
	for i := 0; i < n; i++ {
		vars["x"] = x
		y, err := f.Calc(vars)
		if err != nil {
			return nil, nil, err
		}

		xs[i] = x
		ys[i] = y

		x += xStep
	}
	return xs, ys, nil
}
