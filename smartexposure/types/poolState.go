package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type PoolState struct {
	PoolLiquidity        decimal.Decimal `json:"poolLiquidity"`
	LastRebalance        int64           `json:"lastRebalance"`
	RebalancingInterval  int64           `json:"rebalancingInterval"`
	RebalancingCondition decimal.Decimal `json:"rebalancingCondition"`
	BlockNumber          int64           `json:"blockNumber"`
	BlockTimestamp       time.Time       `json:"blockTimestamp"`
}
