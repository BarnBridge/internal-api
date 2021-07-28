package types

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/barnbridge/internal-api/types"
)

type Tranche struct {
	PoolAddress string `json:"poolAddress,omitempty"`

	ETokenAddress string `json:"eTokenAddress,omitempty"`
	ETokenSymbol  string `json:"eTokenSymbol,omitempty"`

	SFactorE decimal.Decimal `json:"sFactorE,omitempty"`

	TargetRatio decimal.Decimal `json:"targetRatio,omitempty"`
	TokenARatio decimal.Decimal `json:"tokenARatio,omitempty"`
	TokenBRatio decimal.Decimal `json:"tokenBRatio,omitempty"`

	TokenA types.Token `json:"tokenA,omitempty"`
	TokenB types.Token `json:"tokenB,omitempty"`

	RebalancingInterval  string `json:"rebalancingInterval,omitempty"`
	RebalancingCondition string `json:"rebalancingCondition,omitempty"`

	State TrancheState `json:"state,omitempty"`
}

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
