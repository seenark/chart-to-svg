package mychart

import (
	"bytes"
	"fmt"

	"github.com/cnkei/gospline"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

var (
	greenColor = drawing.Color{
		R: 50, G: 161, B: 128, A: 255,
	}
	redColor = drawing.Color{
		R: 230,
		G: 66,
		B: 66,
		A: 255,
	}
)

func GenerateSVG(series []float64) *bytes.Buffer {
	selectedColor := redColor
	if series[0] <= series[len(series)-1] {
		selectedColor = greenColor
	}

	xAxis, yAxis := smoothSeries(series)
	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: xAxis,
				YValues: yAxis,
				Style: chart.Style{
					Show:        true,
					StrokeWidth: 5,
					StrokeColor: selectedColor,
				},
			},
		},
		Width:  720,
		Height: 250,
		Background: chart.Style{
			StrokeColor: drawing.ColorBlack.WithAlpha(1),
			FillColor:   drawing.ColorBlack.WithAlpha(1),
		},
		Canvas: chart.Style{
			StrokeColor: drawing.ColorBlack.WithAlpha(1),
			FillColor:   drawing.ColorBlack.WithAlpha(1),
		},
	}
	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.SVG, buffer)
	if err != nil {
		fmt.Println(err)
	}
	return buffer
}

func smoothSeries(series []float64) (xAxis, yAxis []float64) {
	xAxis = []float64{}
	for index := range series {
		xAxis = append(xAxis, float64(index))
	}
	s := gospline.NewCubicSpline(xAxis, series)
	yAxis = s.Range(0, float64(len(series)-1), 0.1)

	xAxis = []float64{}
	for index := range yAxis {
		xAxis = append(xAxis, float64(index))
	}
	return
}
