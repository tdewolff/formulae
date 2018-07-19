package main

import (
	"fmt"
	"log"

	"github.com/tdewolff/formulae"
)

func main() {
	in := "2+x"
	formula, errs := formulae.Parse(in)
	if len(errs) > 0 {
		log.Fatal(errs)
	}

	compiled, err := formula.Compile()
	if err != nil {
		log.Fatal(err)
	}
	defer compiled.Close()

	y, err := compiled.Calc(5 + 0i)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(y)

	y, err = compiled.Calc(8 + 0i)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(y)

	// xs, ys, err := formula.Interval(0.0, 0.1, 5.0)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// xys := make(plotter.XYs, len(xs))
	// for i := range xs {
	// 	xys[i].X = xs[i]
	// 	xys[i].Y = ys[i]
	// }
	// _, _, ymin, _ := plotter.XYRange(xys)
	// if ymin > 0 {
	// 	ymin = 0
	// }

	// p, err := plot.New()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// p.Title.Text = "Formulae"
	// p.X.Label.Text = "X"
	// p.Y.Label.Text = "Y"
	// p.Y.Min = ymin

	// line, err := plotter.NewLine(xys)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// p.Add(line)
	// p.Legend.Add(in, line)

	// if err := p.Save(4*vg.Inch, 4*vg.Inch, "formula.png"); err != nil {
	// 	log.Fatal(err)
	// }
}
