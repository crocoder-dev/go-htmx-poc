package main

import (
	"github.com/vicanso/go-charts/v2"
)

func CreateBarsChart(values []float64) string {

	p, err := charts.BarRender(
		[][]float64{values},
		charts.XAxisDataOptionFunc([]string{
			"Jan",
			"Feb",
			"Mar",
			"Apr",
			"May",
			"Jun",
			"Jul",
			"Aug",
			"Sep",
			"Oct",
			"Nov",
			"Dec",
		}),
		func(opt *charts.ChartOption) {
			opt.Type = "svg"
			opt.Height = 400
			opt.Width = 800
			opt.SymbolShow = charts.TrueFlag()
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
