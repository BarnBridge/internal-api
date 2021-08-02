package types

import (
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

	RebalancingInterval  int64  `json:"rebalancingInterval,omitempty"`
	RebalancingCondition string `json:"rebalancingCondition,omitempty"`

	State TrancheState `json:"state,omitempty"`
}
