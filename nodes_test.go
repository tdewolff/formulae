package formulae

import "testing"

func TestOptimize(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"0+2", "2"},
		{"2+0", "2"},
		{"2-0", "2"},
		{"0-2", "-2"},
		{"0*2", "0"},
		{"2*0", "0"},
		{"1*2", "2"},
		{"2*1", "2"},
		{"-1*2", "-2"},
		{"2*-1", "-2"},
		{"2/1", "2"},
		{"2/-1", "-2"},
		{"2^0", "1"},
		{"2^1", "2"},
		{"a^-1", "1/a"},
		{"1^2", "1"},
		{"--2", "2"},
		{"a+-1", "a-1"},
		{"a--2", "a+2"},
		{"ln(e)", "1"},
		{"log(e)", "1"},
		{"log10(10)", "1"},
		{"e^ln(x)", "x"},
		{"10^log10(x)", "x"},
		{"sin(0)", "0"},
		{"cos(0)", "1"},
		{"tan(0)", "0"},
		{"sin(-1)", "-sin(1)"},
		{"cos(-1)", "cos(1)"},
		{"tan(-1)", "-tan(1)"},
		{"5+(-1/x)", "5-1/x"},
		{"x+x", "2*x"},
		{"-a*-b", "a*b"},
		{"a*-b", "-(a*b)"},
		{"(-a)^2", "a^2"},
		{"(-a)^3", "-(a^3)"},
		{"x*2", "2*x"},
		{"2*(x*3)", "6*x"},
		{"(2*x)*3", "6*x"},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			f, errs := Parse(test.in)
			if len(errs) > 0 {
				t.Fatal(f, errs)
			}
			f.Optimize()
			if f.root.String() != test.out {
				t.Fatal(f.root.String(), "!=", test.out)
			}
		})
	}
}

func TestDerivative(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"x", "1"},
		{"sin x", "cos(x)"},
		{"x^2", "2*x"},
		{"2^x", "2^x*log(2)"},
		{"e^x", "e^x"},
		{"ln x", "1/x"},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			f, errs := Parse(test.in)
			if len(errs) > 0 {
				t.Fatal(f, errs)
			}
			df := f.Derivative()
			if df.String() != test.out {
				t.Fatal(df.String(), "!=", test.out)
			}
		})
	}
}
