package types

import (
	"github.com/shopspring/decimal"
)

type SeniorBondUser struct {
	SeniorBondId            int64           `json:"seniorBondId"`
	MaturityDate            int64           `json:"maturityDate"`
	Redeemed                bool            `json:"redeemed"`
	AccountAddress          string          `json:"accountAddress"`
	DepositedAmount         decimal.Decimal `json:"depositedAmount"`
	RedeemableAmount        decimal.Decimal `json:"redeemableAmount"`
	UnderlyingTokenAddress  string          `json:"underlyingTokenAddress"`
	UnderlyingTokenSymbol   string          `json:"underlyingTokenSymbol"`
	UnderlyingTokenDecimals int64           `json:"underlyingTokenDecimals"`
	TxHash                  string          `json:"transactionHash"`
	BlockTimestamp          int64           `json:"blockTimestamp"`
}
