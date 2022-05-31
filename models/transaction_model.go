package models

type Transaction struct {
	From   string  `json:"from,omitempty" validate:"required"`
	To     string  `json:"to,omitempty" validate:"required"`
	Amount float64 `json:"amount,omitempty" validate:"required"`
	Fee    float64 `json:"fee,omitempty" validate:"required"`
}
