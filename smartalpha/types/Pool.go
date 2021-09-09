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

	OracleAddress     string `json:"oracleAddress"`
	OracleAssetSymbol string `json:"oracleAssetSymbol"`

	SeniorRateModelAddress string `json:"seniorRateModelAddress"`
	AccountingModelAddress string `json:"accountingModelAddress"`

	Epoch1Start   int64 `json:"epoch1Start"`
	EpochDuration int64 `json:"epochDuration"`

	State PoolState `json:"state"`
	TVL   PoolTVL   `json:"tvl"`
}

type PoolTVL struct {
	EpochJuniorTVL      float64
	EpochSeniorTVL      float64
	JuniorEntryQueueTVL float64
	SeniorEntryQueueTVL float64
	JuniorExitedTVL     float64
	SeniorExitedTVL     float64
}
