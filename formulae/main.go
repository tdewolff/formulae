package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/tdewolff/formulae"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	// Parse formula
	in := "sin(cos(x))^2+1/x-1"
	f, errs := formulae.Parse(in)
	if len(errs) > 0 {
		log.Fatal(errs)
	}
	df := f.Derivative()

	err := writeHTML("math.html", f.LaTeX(), df.LaTeX())
	if err != nil {
		log.Fatal(err)
	}

	// Calculate function
	xs, ys, errs := f.Interval(0.5, 0.01, 5.0)
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

	// Calculate function derivative
	xs2, ys2, errs := df.Interval(0.5, 0.01, 5.0)
	if len(errs) > 0 {
		log.Fatal(errs)
	}
	xys2 := make(plotter.XYs, len(xs2))
	for i := range xs2 {
		xys2[i].X = xs2[i]
		xys2[i].Y = real(ys2[i])
	}

	// Plot functions
	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}

	p.Title.Text = "Formula"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"
	p.Y.Min = ymin

	line, err := plotter.NewLine(xys)
	if err != nil {
		log.Fatal(err)
	}
	line2, err := plotter.NewLine(xys2)
	if err != nil {
		log.Fatal(err)
	}
	line2.LineStyle.Color = color.Gray{192}

	p.Add(plotter.NewGrid())
	p.Add(line)
	p.Add(line2)
	p.Legend.Add("f", line)
	p.Legend.Add("df/dx", line2)

	if err := p.Save(8*vg.Inch, 4*vg.Inch, "formula.png"); err != nil {
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
	fmt.Fprintf(f, `<p style="text-align:center;"><img src="formula.png"></p>
</body>
</html>
`)
	return f.Close()
}
