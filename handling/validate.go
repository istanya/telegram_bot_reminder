package handling

import (
	"strconv"
)

func validDigitDiapason(digitStr string, minDigit int, maxDigit int) bool {
	digitInt, err := strconv.Atoi(digitStr)
	switch {
	case err != nil:
		return false
	case digitInt >= minDigit || digitInt <= maxDigit:
		return true
	}
	return false
}
