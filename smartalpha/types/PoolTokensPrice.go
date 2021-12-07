package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type PoolTokensPrice struct {
	JuniorTokenPrice decimal.Decimal `json:"juniorTokenPrice"`
	SeniorTokenPrice decimal.Decimal `json:"seniorTokenPrice"`
	Timestamp        time.Time       `json:"timestamp"`
}
