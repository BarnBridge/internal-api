package types

import (
	"github.com/shopspring/decimal"
)

type PoolState struct {
	Epoch                  int64           `json:"epoch"`
	SeniorLiquidity        decimal.Decimal `json:"seniorLiquidity"`
	JuniorLiquidity        decimal.Decimal `json:"juniorLiquidity"`
	UpsideExposureRate     decimal.Decimal `json:"upsideExposureRate"`
	DownsideProtectionRate decimal.Decimal `json:"downsideProtectionRate"`
}
