package types

import (
	"github.com/shopspring/decimal"
)

type Token struct {
	Address  string      `json:"address"`
	Symbol   string      `json:"symbol"`
	Decimals int64       `json:"decimals"`
	State    *TokenState `json:"state,omitempty"`
}

type TokenState struct {
	Price          decimal.Decimal `json:"price,omitempty"`
	BlockNumber    int64           `json:"blockNumber,omitempty"`
	BlockTimestamp int64           `json:"blockTimestamp,omitempty"`
}
