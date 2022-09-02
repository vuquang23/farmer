package maths

import "strconv"

func StrToFloat(value string) float64 {
	ret, _ := strconv.ParseFloat(value, 64)
	return ret
}
