package builder

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"farmer/internal/pkg/enum"
	"farmer/internal/pkg/services"
	"farmer/internal/pkg/utils/logger"
)

type IWaveTrendCalculator interface {
	Run(ctx context.Context, market enum.Market, interval string, symlistFilePath string, resultFilePath string) error
}

type waveTrendCalculator struct {
	wtMomentumSvc services.IWavetrendMomentumService
}

func NewWaveTrendCalculator(wtMomentumSvc services.IWavetrendMomentumService) IWaveTrendCalculator {
	return &waveTrendCalculator{
		wtMomentumSvc: wtMomentumSvc,
	}
}

func (w *waveTrendCalculator) Run(
	ctx context.Context, market enum.Market, interval string, symlistFilePath string, resultFilePath string,
) error {
	fi, err := os.Open(symlistFilePath)
	if err != nil {
		return err
	}
	defer fi.Close()

	symbolList := []string{}
	scanner := bufio.NewScanner(fi)
	for scanner.Scan() {
		symbol := scanner.Text()
		symbolList = append(symbolList, symbol)
	}

	ret, svcErr := w.wtMomentumSvc.Calculate(ctx, market, symbolList, interval)
	if svcErr != nil {
		return svcErr
	}

	fo, err := os.Create(resultFilePath)
	if err != nil {
		return err
	}
	defer fo.Close()

	for _, r := range ret {
		fo.Write([]byte(
			fmt.Sprintf("%s%s%f\n", r.Symbol, strings.Repeat(" ", 25-len(r.Symbol)), r.Value),
		))
	}

	logger.Infof(ctx, "[Run] calculated for %d symbols", len(ret))
	return nil
}
