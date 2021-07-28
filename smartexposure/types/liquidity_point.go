package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type SELiquidityPoint struct {
	Point           time.Time       `json:"point"`
	TokenALiquidity decimal.Decimal `json:"tokenALiquidity"`
	TokenBLiquidity decimal.Decimal `json:"tokenBLiquidity"`
}
