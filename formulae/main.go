package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tdewolff/formulae"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	in := "sin(cos(x))^2+1/(x*2)+0.001x^(3+x)"
	formula, errs := formulae.Parse(in)
	if len(errs) > 0 {
		log.Fatal(errs)
	}

	dformula := formula.Derivative()

	err := writeHTML("math.html", formula.LaTeX(), dformula.LaTeX())
	if err != nil {
		log.Fatal(err)
	}

	xs, ys, errs := formula.Interval(0.1, 0.1, 5.0)
	if len(errs) > 0 {
		log.Fatal(errs)
	}

	xys := make(plotter.XYs, len(xs))
	for i := range xs {
		xys[i].X = xs[i]
		xys[i].Y = real(ys[i])
	}
	_, _, ymin, _ := plotter.XYRange(xys)
	if ymin > 0 {
		ymin = 0
	}

	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}

	p.Title.Text = "Formulae"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Y.Min = ymin

	line, err := plotter.NewLine(xys)
	if err != nil {
		log.Fatal(err)
	}

	p.Add(line)
	p.Legend.Add(in, line)

	if err := p.Save(4*vg.Inch, 4*vg.Inch, "formula.png"); err != nil {
		log.Fatal(err)
	}
}

func writeHTML(filename string, latex ...string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	fmt.Fprintf(f, `<!DOCTYPE html>
<html>
<head>
    <script type="text/javascript" async src="https://cdnjs.cloudflare.com/ajax/libs/mathjax/2.7.5/MathJax.js?config=TeX-MML-AM_CHTML"></script>
</head>
<body>
`)
	for _, x := range latex {
		fmt.Fprintf(f, "    <p>$$%s$$</p>\n", x)
	}
	fmt.Fprintf(f, `</body>
</html>
`)
	return f.Close()
}
