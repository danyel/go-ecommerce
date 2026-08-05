package types

type Price struct {
	Inclusive Float64 `json:"inclusive"`
	Exclusive Float64 `json:"exclusive"`
	Tax       Float64 `json:"tax"`
	Currency  string  `json:"currency"`
}
