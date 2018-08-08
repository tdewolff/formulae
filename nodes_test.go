package formulae

import "testing"

func TestOptimize(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"0+2", "'2'"},
		{"2+0", "'2'"},
		{"2-0", "'2'"},
		{"0-2", "'-2'"},
		{"0*2", "'0'"},
		{"2*0", "'0'"},
		{"1*2", "'2'"},
		{"2*1", "'2'"},
		{"-1*2", "'-2'"},
		{"2*-1", "'-2'"},
		{"2/1", "'2'"},
		{"2/-1", "'-2'"},
		{"2^0", "'1'"},
		{"2^1", "'2'"},
		{"a^-1", "('1' / 'a')"},
		{"1^2", "'1'"},
		{"--2", "'2'"},
		{"a+-1", "('a' - '1')"},
		{"a--2", "('a' + '2')"},
		{"ln(e)", "'1'"},
		{"log(e)", "'1'"},
		{"sin(0)", "'0'"},
		{"cos(0)", "'1'"},
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
		{"x", "'1'"},
		{"sin x", "(cos 'x')"},
		{"x^2", "('2' * 'x')"},
		{"2^x", "(('2' ^ 'x') * (ln '2'))"},
		{"e^x", "('e' ^ 'x')"},
		{"ln x", "('1' / 'x')"},
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
