package util

const (
	USD = "USD"
	TK  = "Tk"
	EUR = "EUR"
	CAD = "CAD"
)

func IsSupporpedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, TK, CAD:
		return true
	}
	return false
}
