package indicators

import (
	"math"

	"github.com/cinar/indicator"

	e "farmer/internal/pkg/entities"
)

func CalculatePastWavetrendData(candles []*e.MinimalKline, n1, n2 uint64) (*e.PastWavetrend, error) {
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

	ret := &e.PastWavetrend{
		LastOpenTime: candles[len(candles)-1].OpenTime,
		LastD:        d[len(d)-1],
		LastEsa:      esa[len(esa)-1],
		PastTci:      tci[len(tci)-30:],
	}
	return ret, nil
}

func CalculatePastWavetrendDataWithNewCandles(wt *e.PastWavetrend, candles []*e.MinimalKline, n1, n2 uint64) (*e.PastWavetrend, error) {
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

	return wt, nil
}

func nextEma(lastEma, currentVal float64, period uint64) float64 {
	k := float64(2) / float64(1+period)
	return (currentVal * k) + (lastEma * float64(1-k))
}

func CalculateCurrentTciFromPastWavetrendDatAndCurrentCandle(wt *e.PastWavetrend, candles []*e.MinimalKline, n1, n2 uint64) float64 {
	wt, _ = CalculatePastWavetrendDataWithNewCandles(wt, candles, n1, n2)
	return wt.PastTci[len(wt.PastTci)-1]
}
