package main

import (
	"math/rand"

	"github.com/vicanso/go-charts/v2"
)

func CreateBarsChart() string {

	// Generate random data for the bars
	var values []float64
	for i := 0; i < 12; i++ {
		// Generate random float64 values between 0 and 200
		value := rand.Float64() * 200
		values = append(values, value)
	}

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
