package types

import (
	"github.com/barnbridge/internal-api/types"
)

type PoolPreviousEpoch struct {
	PoolAddress       string `json:"poolAddress"`
	PoolName          string `json:"poolName"`
	PoolToken         types.Token
	OracleAssetSymbol string  `json:"oracleAssetSymbol"`
	Epochs            []Epoch `json:"epochs"`
}
