package types

import (
	"github.com/shopspring/decimal"
)

type JuniorBondUser struct {
	JuniorBondId            int64           `json:"juniorBondId"`
	AccountAddress          string          `json:"accountAddress"`
	DepositedAmount         decimal.Decimal `json:"depositedAmount"`
	MaturityDate            int64           `json:"maturityDate"`
	Redeemed                bool            `json:"redeemed"`
	UnderlyingTokenAddress  string          `json:"underlyingTokenAddress"`
	UnderlyingTokenSymbol   string          `json:"underlyingTokenSymbol"`
	UnderlyingTokenDecimals int64           `json:"underlyingTokenDecimals"`
	TxHash                  string          `json:"transactionHash"`
	BlockTimestamp          int64           `json:"blockTimestamp"`
}
