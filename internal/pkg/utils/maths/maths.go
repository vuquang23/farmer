package maths

import (
	"errors"
	"math"
	"strconv"
)

func GetPrecision(precisionString string) (int, error) {
	f, err := strconv.ParseFloat(precisionString, 64)
	if err != nil {
		return 0, err
	}
	if f == 0 {
		return 0, errors.New("precision string is zero") // should not happen
	}
	precision := math.Log10(1 / f)
	return int(precision), nil
}

func RoundingUp(value float64, precision int) string {
	p10 := math.Pow10(precision)
	rounded := math.Ceil(value*p10) / p10
	ret := strconv.FormatFloat(rounded, 'f', precision, 64)
	return ret
}
