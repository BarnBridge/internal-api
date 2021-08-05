package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type TokensPricePoint struct {
	Point            time.Time       `json:"point"`
	SeniorTokenPrice decimal.Decimal `json:"seniorTokenPrice"`
	JuniorTokenPrice decimal.Decimal `json:"juniorTokenPrice"`
}
