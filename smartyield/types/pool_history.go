package types

import (
	"github.com/shopspring/decimal"
)

type PoolHistory struct {
	AccountAddress         *string         `json:"accountAddress,omitempty"`
	ProtocolId             string          `json:"protocolId"`
	Pool                   string          `json:"pool"`
	UnderlyingTokenAddress string          `json:"underlyingTokenAddress"`
	UnderlyingTokenSymbol  string          `json:"underlyingTokenSymbol"`
	Amount                 decimal.Decimal `json:"amount"`
	Tranche                string          `json:"tranche"`
	TransactionType        string          `json:"transactionType"`
	TransactionHash        string          `json:"transactionHash"`
	BlockTimestamp         int64           `json:"blockTimestamp"`
	BlockNumber            int64           `json:"blockNumber"`
}
