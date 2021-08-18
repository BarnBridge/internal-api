package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type PerformancePoint struct {
	Point           time.Time       `json:"point"`
	SeniorWithSA    decimal.Decimal `json:"seniorWithSA"`
	SeniorWithoutSA decimal.Decimal `json:"seniorWithoutSA"`
	JuniorWithSA    decimal.Decimal `json:"juniorWithSA"`
	JuniorWithoutSA decimal.Decimal `json:"juniorWithoutSA"`
	UnderlyingPrice decimal.Decimal `json:"underlyingPrice"`
}
