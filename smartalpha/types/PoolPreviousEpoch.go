package types

import (
	"github.com/barnbridge/internal-api/types"
)

type PoolDetails struct {
	PoolAddress       string      `json:"poolAddress"`
	PoolName          string      `json:"poolName"`
	PoolToken         types.Token `json:"poolToken"`
	OracleAssetSymbol string      `json:"oracleAssetSymbol"`
}
