package indicators

import (
	"math"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/cinar/indicator"

	e "farmer/internal/pkg/entities"
)

func WaveTrendMomentumValue(candles []*e.MinimalKline, n1, n2 uint64) (float64, error) {
	ap := Hlc3(candles)
	esa := indicator.Ema(int(n1), ap)

	// ap[0] = esa[0]
	ap = ap[1:]
	esa = esa[1:]

	abs := []float64{}
	for i := 0; i < len(esa); i++ {
		abs = append(abs, math.Abs(ap[i]-esa[i]))
	}

	d := indicator.Ema(int(n1), abs)

	ci := []float64{}
	for i := 0; i < len(d); i++ {
		ci = append(ci, (ap[i]-esa[i])/(0.015*d[i]))
	}

	tci := indicator.Ema(int(n2), ci)
	return tci[len(tci)-1], nil

}

func Hlc3(candles []*e.MinimalKline) []float64 {
	ret := []float64{}
	for _, k := range candles {
		ret = append(ret, (k.High+k.Low+k.Close)/3)
	}
	return ret
}

func SpotKlineToMinimalKline(candles []*binance.Kline) []*e.MinimalKline {
	ret := []*e.MinimalKline{}
	for _, c := range candles {
		high, _ := strconv.ParseFloat(c.High, 64)
		low, _ := strconv.ParseFloat(c.Low, 64)
		close, _ := strconv.ParseFloat(c.Close, 64)
		openTime := uint64(c.OpenTime)
		ret = append(ret, &e.MinimalKline{
			OpenTime: openTime,
			High:     high,
			Low:      low,
			Close:    close,
		})
	}
	return ret
}

func FutureKlineToMinimalKline(candles []*futures.Kline) []*e.MinimalKline {
	ret := []*e.MinimalKline{}
	for _, c := range candles {
		high, _ := strconv.ParseFloat(c.High, 64)
		low, _ := strconv.ParseFloat(c.Low, 64)
		close, _ := strconv.ParseFloat(c.Close, 64)
		openTime := uint64(c.OpenTime)
		ret = append(ret, &e.MinimalKline{
			OpenTime: openTime,
			High:     high,
			Low:      low,
			Close:    close,
		})
	}
	return ret
}
