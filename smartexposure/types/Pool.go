package types

import (
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
