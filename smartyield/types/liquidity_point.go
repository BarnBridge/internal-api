package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type LiquidityPoint struct {
	Point           time.Time       `json:"point"`
	SeniorLiquidity decimal.Decimal `json:"seniorLiquidity"`
	JuniorLiquidity decimal.Decimal `json:"juniorLiquidity"`
}
