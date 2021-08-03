package types

import (
	"github.com/shopspring/decimal"
)

type SeniorRedeem struct {
	SeniorBondAddress string          `json:"seniorBondAddress"`
	UserAddress       string          `json:"userAddress"`
	SeniorBondID      int64           `json:"seniorBondId"`
	PoolAddress       string          `json:"smartYieldAddress"`
	Fee               decimal.Decimal `json:"fee"`
	UnderlyingIn      decimal.Decimal `json:"underlyingIn"`
	Gain              decimal.Decimal `json:"gain"`
	ForDays           int64           `json:"forDays"`
	BlockTimestamp    int64           `json:"blockTimestamp"`
	TxHash            string          `json:"transactionHash"`
}
