package formulae

import (
	"fmt"
	"math"

	"github.com/aarzilli/golua/lua"
)

type Vars map[string]complex128

func (v Vars) Set(name string, value complex128) {
	v[name] = value
}

////////////////

type Formula struct {
	root Node
	Vars
}

func NewFormula(root Node) *Formula {
	vars := Vars{}
	vars.Set("e", complex(math.E, 0))
	vars.Set("pi", complex(math.Pi, 0))
	vars.Set("phi", complex(math.Phi, 0))
	return &Formula{root: root, Vars: vars}
}

func (f *Formula) Calc(x complex128) (complex128, error) {
	return f.root.Calc(x, f.Vars)
}

func (f *Formula) Interval(xMin, xStep, xMax float64) ([]float64, []complex128, []error) {
	n := int((xMax-xMin)/xStep) + 1
	xs := make([]float64, n)
	ys := make([]complex128, n)

	x := xMin
	var errs []error
	for i := 0; i < n; i++ {
		y, err := f.Calc(complex(x, 0))
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

func (f *Formula) Compile() (LuaFormula, error) {
	L := lua.NewState()
	L.OpenLibs()
	L.OpenMath()

	luaFunc := "function formula(x) return " + f.root.Lua() + " end"
	err := L.DoString(luaFunc)
	if err != nil {
		return LuaFormula{nil}, err
	}
	return LuaFormula{L}, nil
}

type LuaFormula struct {
	*lua.State
}

func (L LuaFormula) Calc(x complex128) (complex128, error) {
	L.GetGlobal("formula")
	L.PushNumber(real(x))
	if err := L.Call(1, 1); err != nil {
		return 0, err
	}
	y := L.ToNumber(1)
	L.Pop(1)
	return complex(y, 0), nil
}
