package main

import (
	"fmt"

	"github.com/vicanso/go-charts/v2"
)

func CreateLineChart(values []float64) string {

	p, err := charts.LineRender(
		[][]float64{values},
		charts.XAxisDataOptionFunc([]string{
			"Jul 24",
			"Jul 31",
			"Aug 7",
			"Aug 14",
			"Aug 21",
			"Aug 28",
			"Sep 4",
			"Sep 11",
			"Sep 18",
		}),
		func(opt *charts.ChartOption) {
			opt.Type = "svg"
			opt.Height = 400
			opt.Width = 800
			opt.SymbolShow = charts.TrueFlag()
			opt.LineStrokeWidth = 2
			opt.ValueFormatter = func(f float64) string {
				return fmt.Sprintf("%.0f", f)
			}
		},
	)

	if err != nil {
		panic(err)
	}

	buf, err := p.Bytes()
	if err != nil {
		panic(err)
	}
	return string(buf[:])
}
