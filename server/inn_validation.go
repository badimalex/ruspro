package server

import (
	"strconv"
)

func isValidINN(inn string) bool {
	length := len(inn)

	switch length {
	case 10:
		return checkINN10(inn)
	case 12:
		return checkINN12(inn)
	default:
		return false
	}
}

func checkINN10(inn string) bool {
	coefficients := []int{2, 4, 10, 3, 5, 9, 4, 6, 8, 0}
	return checkINN(inn, coefficients)
}

func checkINN12(inn string) bool {
	coefficients1 := []int{7, 2, 4, 10, 3, 5, 9, 4, 6, 8, 0, 0}
	coefficients2 := []int{3, 7, 2, 4, 10, 3, 5, 9, 4, 6, 8, 0}
	return checkINN(inn, coefficients1) && checkINN(inn, coefficients2)
}

func checkINN(inn string, coefficients []int) bool {
	var sum int
	for i, c := range coefficients {
		digit, err := strconv.Atoi(string(inn[i]))
		if err != nil {
			return false
		}
		sum += digit * c
	}
	return sum%11%10 == int(inn[len(inn)-1]-'0')
}
