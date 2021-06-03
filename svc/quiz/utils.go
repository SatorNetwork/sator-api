package quiz

import "math"

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func calcRate(questionTimeSec, answerTimeSec float64) int {
	return int(4 - math.Ceil(answerTimeSec/(questionTimeSec/4)))
}
