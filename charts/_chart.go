package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/wcharczuk/go-chart/v2"
)

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

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

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:         "The XAxis",
			TickPosition: chart.TickPositionBetweenTicks,
			ValueFormatter: func(v interface{}) string {
				typed := v.(float64)
				typedDate := chart.TimeFromFloat64(typed)
				return fmt.Sprintf("%d-%d-%d", typedDate.Month(), typedDate.Day(), typedDate.Year())
			},
		},
		YAxis: chart.YAxis{
			AxisType: chart.YAxisSecondary,
			Name:     "Stress (MPa)",
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
				},
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},
		},
	}

	start := time.Now()

	numRange := makeRange(1, end)

	for _, num := range numRange {
		outName := fmt.Sprintf("result/output_%d.png", num)
		f, _ := os.Create(outName)
		defer f.Close()
		graph.Render(chart.PNG, f)
	}

	elapsed := time.Since(start)
	log.Printf("Exporting images took %s", elapsed)
	// io.WriteString(os.Stdout, "Exporting images took %s")
}
