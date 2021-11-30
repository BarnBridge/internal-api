package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type PoolState struct {
	BlockNumber           int64           `json:"blockNumber"`
	BlockTimestamp        time.Time       `json:"blockTimestamp"`
	SeniorLiquidity       decimal.Decimal `json:"seniorLiquidity"`
	JuniorLiquidity       decimal.Decimal `json:"juniorLiquidity"`
	JTokenPrice           decimal.Decimal `json:"jTokenPrice"`
	SeniorAPY             float64         `json:"seniorApy"`
	JuniorAPY             float64         `json:"juniorApy"`
	JuniorAPYPast30dAvg   float64         `json:"juniorAPYPast30DAvg"`
	OriginatorApy         float64         `json:"originatorApy"`
	OriginatorNetApy      float64         `json:"originatorNetApy"`
	AvgSeniorMaturityDays float64         `json:"avgSeniorMaturityDays"`
	NumberOfSeniors       int64           `json:"numberOfSeniors"`
	NumberOfJuniors       int64           `json:"numberOfJuniors"`
	NumberOfJuniorsLocked int64           `json:"numberOfJuniorsLocked"`
	JuniorLiquidityLocked decimal.Decimal `json:"juniorLiquidityLocked"`
}
