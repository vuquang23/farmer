package indicators

import (
	"math"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/cinar/indicator"

	e "farmer/internal/pkg/entities"
)

func CalcPastWavetrendData(
	candles []*e.MinimalKline,
	n1 uint64, n2 uint64, tciLen uint64,
	tciAvgPeriod uint64, difWavetrendLen uint64,
) (*e.PastWavetrend, error) {
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
		// FIXME: d can be zero?
		// if d[i] == 0 {
		// 	for _, c := range candles {
		// 		logger.Logger.Sugar().Infof("%+v", c)
		// 	}
		// 	panic("d is zero")
		// }
		ci = append(ci, (ap[i]-esa[i])/(0.015*d[i]))
	}

	tci := indicator.Ema(int(n2), ci)

	cachedTci := tci[len(tci)-int(tciLen):]
	difWavetrend := calcDifWavetrend(cachedTci, tciAvgPeriod, difWavetrendLen)

	ret := &e.PastWavetrend{
		LastOpenTime: candles[len(candles)-1].OpenTime,
		LastD:        d[len(d)-1],
		LastEsa:      esa[len(esa)-1],
		PastTci:      cachedTci,
		DifWavetrend: difWavetrend,
	}
	return ret, nil
}

// calcDifWavetrend calculate difference between wt1 and wt2.
// constraint: avgPeriod + limit <= len(tci)
// avgPeriod: average along 3 minutes, 3 hours,...
func calcDifWavetrend(tci []float64, avgPeriod uint64, limit uint64) []float64 {
	lenTci := len(tci)
	currentSum := float64(0)
	ret := []float64{}

	for i := lenTci - int(limit) - int(avgPeriod); i < lenTci; i++ {
		currentSum += tci[i]
		if i >= lenTci-int(limit) {
			currentSum -= tci[i-int(avgPeriod)]
			ret = append(ret, tci[i]-currentSum/float64(avgPeriod))
		}
	}

	return ret
}

func UpdatePastWavetrendDataWithNewCandles(
	wt *e.PastWavetrend, candles []*e.MinimalKline,
	n1, n2 uint64, avgPeriod uint64, difWavetrendLen uint64,
) (*e.PastWavetrend, error) {
	for _, c := range candles {
		currentHlc3 := (c.High + c.Low + c.Close) / 3
		lastTci := wt.PastTci[len(wt.PastTci)-1]

		currentEsa := nextEma(wt.LastEsa, currentHlc3, n1)
		currentAbs := math.Abs(currentHlc3 - currentEsa)
		currentD := nextEma(wt.LastD, currentAbs, n1)
		currentCi := (currentHlc3 - currentEsa) / (0.015 * currentD)
		currentTci := nextEma(lastTci, currentCi, n2)

		wt.LastOpenTime = c.OpenTime
		wt.LastD = currentD
		wt.LastEsa = currentEsa
		wt.PastTci = append(wt.PastTci, currentTci)
		wt.PastTci = wt.PastTci[1:]
	}

	difWavetrend := calcDifWavetrend(wt.PastTci, avgPeriod, difWavetrendLen)
	wt.DifWavetrend = difWavetrend

	return wt, nil
}

func nextEma(lastEma, currentVal float64, period uint64) float64 {
	k := float64(2) / float64(1+period)
	return (currentVal * k) + (lastEma * float64(1-k))
}

func CalcCurrentTciAndDifWavetrend(
	wt *e.PastWavetrend, candles []*e.MinimalKline,
	n1, n2 uint64, avgPeriod uint64,
) (float64, float64) {
	wt, _ = UpdatePastWavetrendDataWithNewCandles(wt, candles, n1, n2, avgPeriod, 0)

	sumTci := float64(0)
	for i := len(wt.PastTci) - int(avgPeriod); i < len(wt.PastTci); i++ {
		sumTci += wt.PastTci[i]
	}

	currentTci := wt.PastTci[len(wt.PastTci)-1]
	currentDifWavetrend := currentTci - sumTci/float64(avgPeriod)

	return currentTci, currentDifWavetrend
}

// below is functions used for wtmomentum command

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
