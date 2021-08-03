package types

import (
	"github.com/shopspring/decimal"
)

type Token struct {
	TokenAddress  string      `json:"tokenAddress"`
	TokenSymbol   string      `json:"tokenSymbol"`
	TokenDecimals int64       `json:"tokenDecimals"`
	State         *TokenState `json:"state,omitempty"`
}

type TokenState struct {
	Price          decimal.Decimal `json:"price,omitempty"`
	BlockNumber    int64           `json:"blockNumber,omitempty"`
	BlockTimestamp int64           `json:"blockTimestamp,omitempty"`
}
