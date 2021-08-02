package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type ETokenPricePoint struct {
	Point       time.Time       `json:"point"`
	ETokenPrice decimal.Decimal `json:"eTokenPrice"`
}
