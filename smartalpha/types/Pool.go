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

	UserHasActivePosition *bool `json:"userHasActivePosition,omitempty"`
}

type PoolTVL struct {
	EpochJuniorTVL      float64 `json:"epochJuniorTVL"`
	EpochSeniorTVL      float64 `json:"epochSeniorTVL"`
	JuniorEntryQueueTVL float64 `json:"juniorEntryQueueTVL"`
	SeniorEntryQueueTVL float64 `json:"seniorEntryQueueTVL"`
	JuniorExitQueueTVL  float64 `json:"juniorExitQueueTVL"`
	SeniorExitQueueTVL  float64 `json:"seniorExitQueueTVL"`
	JuniorExitedTVL     float64 `json:"juniorExitedTVL"`
	SeniorExitedTVL     float64 `json:"seniorExitedTVL"`
}
