package types

import (
	"github.com/barnbridge/internal-api/types"
)

type UserQueuePosition struct {
	PoolAddress       string      `json:"poolAddress"`
	PoolName          string      `json:"poolName"`
	PoolToken         types.Token `json:"poolToken"`
	OracleAssetSymbol string      `json:"oracleAssetSymbol"`
	Tranche           string      `json:"tranche"`
	QueueType         string      `json:"queueType"`
	BlockTimestamp    int64       `json:"blockTimestamp"`
}
