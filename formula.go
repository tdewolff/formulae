package formulae

import (
	"math"
	"strings"

	"github.com/aarzilli/golua/lua"
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

type Formula struct {
	root Node
	Calc
}

func (f *Formula) Compile(vars Vars) error {
	var err error
	f.Calc, err = f.root.Compile(vars)
	return err
}

func (f *Formula) Interval(xMin, xStep, xMax float64) ([]float64, []complex128) {
	n := int((xMax-xMin)/xStep) + 1
	xs := make([]float64, n)
	ys := make([]complex128, n)

	x := xMin
	for i := 0; i < n; i++ {
		y := f.Calc(complex(x, 0))

		xs[i] = x
		ys[i] = y

		x += xStep
	}
	return xs, ys
}

type LuaFormula struct {
	*lua.State
}

func (f *Formula) CompileLua() (LuaFormula, error) {
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
