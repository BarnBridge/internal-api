package types

import (
	globalTypes "github.com/barnbridge/internal-api/types"
)

type Pool struct {
	PoolName    string            `json:"poolName"`
	PoolAddress string            `json:"poolAddress"`
	PoolToken   globalTypes.Token `json:"poolToken"`

	JuniorTokenAddress string `json:"juniorTokenAddress"`
	SeniorTokenAddress string `json:"seniorTokenAddress"`
	OracleAddress      string `json:"oracleAddress"`
	OracleAssetSymbol  string `json:"oracleAssetSymbol"`
}
