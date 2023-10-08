package main

import (
	"bytes"
	"github.com/wcharczuk/go-chart/v2"
)

func CreateTestChart() string {
	mainSeries := chart.ContinuousSeries{
		Name:    "A test series",
		XValues: chart.Seq{Sequence: chart.NewLinearSequence().WithStart(1.0).WithEnd(100.0)}.Values(),
		YValues: chart.Seq{Sequence: chart.NewRandomSequence().WithLen(100).WithMin(0).WithMax(100)}.Values(),
	}

	smaSeries := &chart.SMASeries{
		InnerSeries: mainSeries,
	}

	graph := chart.Chart{
		Series: []chart.Series{
			mainSeries,
			smaSeries,
		},
	}

	var buffer = bytes.NewBuffer([]byte{})
	var err = graph.Render(chart.SVG, buffer)

	if err != nil {
		panic(err)
	}

	return buffer.String()
}
