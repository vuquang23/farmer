package spotworker

import (
	"fmt"
	"os"
	"time"

	c "farmer/internal/pkg/constants"
	"farmer/internal/pkg/utils/logger"
)

func (w *spotWorker) runMainProcessor() {
	ID := w.setting.loadSymbol()
	log := logger.WithDescription(fmt.Sprintf("%s - Main proccessor", ID))
	log.Sugar().Infof("Worker started")

	ticker := time.NewTicker(c.ProcessingFrequencyTime)
	for ; !w.getStopSignal(); <-ticker.C {
		log.Sugar().Infof("Value TCI: %f", w.waveTrendDat.loadCurrentTci())
		log.Sugar().Infof("Value Dif wavetrend: %f", w.waveTrendDat.loadCurrentDifWavetrend())

		if w.shouldBuy() && !w.isDowntrendOnSecondaryWavetrend() {
			log.Sugar().Infof("Should buy now %s", time.Now().String())
			os.Exit(0)
			continue
		}

	}
}

func (w *spotWorker) shouldBuy() bool {
	currentTci := w.waveTrendDat.loadCurrentTci()
	if currentTci > c.WavetrendOversold {
		return false
	}

	currentDifWt := w.waveTrendDat.loadCurrentDifWavetrend()
	if currentDifWt <= 0 {
		return false
	}

	pastWtDat := w.waveTrendDat.loadPastWaveTrendData()
	for i := len(pastWtDat.PastTci) - c.OversoldRequiredTime; i < len(pastWtDat.PastTci); i++ {
		if pastWtDat.PastTci[i] > c.WavetrendOversold {
			return false
		}
	}

	for i := len(pastWtDat.DifWavetrend) - c.OversoldNegativeDifWtRequiredTime - c.OversoldPositiveDifWtRequiredTime; i < len(pastWtDat.DifWavetrend); i++ {
		if i < len(pastWtDat.DifWavetrend)-c.OversoldPositiveDifWtRequiredTime {
			if pastWtDat.DifWavetrend[i] > 0 {
				return false
			}
		} else {
			if pastWtDat.DifWavetrend[i] <= 0 {
				return false
			}
		}
	}

	return true
}
