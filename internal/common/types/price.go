package types

type Price struct {
	Inclusive Float64 `json:"inclusive"`
	Exclusive Float64 `json:"exclusive"`
	Tax       Float64 `json:"tax"`
	Currency  string  `json:"currency"`
}

func NewPrice(amount float64, currency string) Price {
	return Price{
		Currency:  currency,
		Inclusive: Float64(amount),
		Tax:       Float64(amount - (amount / 1.21)),
		Exclusive: Float64(amount / 1.21),
	}
}
