package main

import (
	"log"

	"github.com/tdewolff/formulae"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	in := "2+x"
	formula, errs := formulae.Parse(in)
	if len(errs) > 0 {
		log.Fatal(errs)
	}

	err := formula.Compile(formulae.DefaultVars)
	if err != nil {
		log.Fatal(err)
	}

	xs, ys := formula.Interval(0.0, 0.1, 5.0)

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
