package types

import (
	"github.com/shopspring/decimal"
)

type Epoch struct {
	Id                     int64           `json:"id"`
	SeniorLiquidity        decimal.Decimal `json:"seniorLiquidity"`
	JuniorLiquidity        decimal.Decimal `json:"juniorLiquidity"`
	UpsideExposureRate     decimal.Decimal `json:"upsideExposureRate"`
	DownsideProtectionRate decimal.Decimal `json:"downsideProtectionRate"`
	StartDate              *int64          `json:"startDate"`
	EndDate                *int64          `json:"endDate"`
	EntryPrice             decimal.Decimal `json:"entryPrice"`
	JuniorProfits          decimal.Decimal `json:"juniorProfits"`
	SeniorProfits          decimal.Decimal `json:"seniorProfits"`
	JuniorTokenPriceStart  decimal.Decimal `json:"juniorTokenPriceStart"`
	SeniorTokenPriceStart  decimal.Decimal `json:"seniorTokenPriceStart"`
}
