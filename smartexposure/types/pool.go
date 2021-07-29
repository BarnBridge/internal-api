package types

import (
	"time"

	"github.com/shopspring/decimal"

	globalTypes "github.com/barnbridge/internal-api/types"
)

type Pool struct {
	PoolAddress string `json:"poolAddress"`
	PoolName    string `json:"poolName"`

	TokenA globalTypes.Token `json:"tokenA"`
	TokenB globalTypes.Token `json:"tokenB"`

	StartAtBlock int64
	State        PoolState `json:"state,omitempty"`
}

type PoolState struct {
	PoolLiquidity        decimal.Decimal `json:"poolLiquidity"`
	LastRebalance        int64           `json:"lastRebalance"`
	RebalancingInterval  int64           `json:"rebalancingInterval"`
	RebalancingCondition decimal.Decimal `json:"rebalancingCondition"`
	BlockNumber          int64           `json:"blockNumber"`
	BlockTimestamp       time.Time       `json:"blockTimestamp"`
}
