package types

type PriceDTO struct {
	Inclusive Float64 `json:"inclusive"`
	Exclusive Float64 `json:"exclusive"`
	Tax       Float64 `json:"tax"`
	Currency  string  `json:"currency"`
}

type Price struct {
	Inclusive Float64
	Exclusive Float64
	Tax       Float64
	Currency  string
}

func (price *Price) DTO() PriceDTO {
	return PriceDTO{
		Inclusive: price.Inclusive,
		Exclusive: price.Exclusive,
		Tax:       price.Tax,
		Currency:  price.Currency,
	}
}

func FromPrice(price Price) PriceDTO {
	return PriceDTO{
		Inclusive: price.Inclusive,
		Exclusive: price.Exclusive,
		Tax:       price.Tax,
		Currency:  price.Currency,
	}
}

func NewPrice(amount float64, currency string) Price {
	return Price{
		Currency:  currency,
		Inclusive: Float64(amount),
		Tax:       Float64(amount - (amount / 1.21)),
		Exclusive: Float64(amount / 1.21),
	}
}
