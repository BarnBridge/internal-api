package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type TrancheState struct {
	TokenALiquidity decimal.Decimal `json:"tokenALiquidity,omitempty"`
	TokenBLiquidity decimal.Decimal `json:"tokenBLiquidity,omitempty"`

	ETokenPrice decimal.Decimal `json:"eTokenPrice,omitempty"`

	CurrentRatio       decimal.Decimal `json:"currentRatio,omitempty"`
	TokenACurrentRatio decimal.Decimal `json:"tokenACurrentRatio,omitempty"`
	TokenBCurrentRatio decimal.Decimal `json:"tokenBCurrentRatio,omitempty"`

	LastRebalance int64 `json:"lastRebalance,omitempty"`

	BlockNumber    int64     `json:"blockNumber,omitempty"`
	BlockTimestamp time.Time `json:"blockTimestamp,omitempty"`
}
