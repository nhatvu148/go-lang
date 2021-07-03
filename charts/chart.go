package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	input := os.Args[1]
	end, err := strconv.Atoi(input)
	if err != nil {
		log.Fatal("Input is not a number!")
		return
	} else if end > 50 {
		log.Fatal("Input is too large!")
		return
	}

	// Get some random points
	rand.Seed(int64(0))
	n := 50
	// scatterData := randomPoints(n)
	lineData := randomPoints(n)
	fmt.Printf("%T, %v\n", lineData, lineData)
	// linePointsData := randomPoints(n)

	// Create a new plot, set its title and
	// axis labels.
	p := plot.New()

	p.Title.Text = "DFP"
	// p.X.Label.Text = "X"
	p.Y.Label.Text = "Stress [MPa]"
	// Draw a grid behind the data
	p.Add(plotter.NewGrid())

	// Make a scatter plotter and set its style.
	// s, err := plotter.NewScatter(scatterData)
	// if err != nil {
	// 	panic(err)
	// }
	// s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}

	// Make a line plotter and set its style.
	l, err := plotter.NewLine(lineData)
	if err != nil {
		panic(err)
	}
	l.LineStyle.Width = vg.Points(1)
	// l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}

	// Make a line plotter with points and set its style.
	// lpLine, lpPoints, err := plotter.NewLinePoints(linePointsData)
	// if err != nil {
	// 	panic(err)
	// }
	// lpLine.Color = color.RGBA{G: 255, A: 255}
	// lpPoints.Shape = draw.PyramidGlyph{}
	// lpPoints.Color = color.RGBA{R: 255, A: 255}

	// Add the plotters to the plot, with a legend
	// entry for each
	p.Add(l)
	// p.Add(s, l, lpLine, lpPoints)
	// p.Legend.Add("scatter", s)
	p.Legend.Add("line", l)
	// p.Legend.Add("line points", lpLine, lpPoints)

	start := time.Now()

	numRange := makeRange(1, end)

	for _, num := range numRange {
		outName := fmt.Sprintf("result/output_%d.png", num)

		// Save the plot to a PNG file.
		if err := p.Save(4*vg.Inch, 4*vg.Inch, outName); err != nil {
			panic(err)
		}
	}

	elapsed := time.Since(start)
	log.Printf("Exporting images took %s", elapsed)
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

// randomPoints returns some random x, y points.
func randomPoints(n int) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := range pts {
		if i == 0 {
			pts[i].X = rand.Float64()
		} else {
			pts[i].X = pts[i-1].X + rand.Float64()
		}
		pts[i].Y = pts[i].X + 10*rand.Float64()
	}
	return pts
}
