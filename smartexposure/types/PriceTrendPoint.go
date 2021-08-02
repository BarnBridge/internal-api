package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type PriceTrendPoint struct {
	Point       time.Time       `json:"point"`
	TokenAPrice decimal.Decimal `json:"tokenAPrice"`
	TokenBPrice decimal.Decimal `json:"tokenBPrice"`
}
