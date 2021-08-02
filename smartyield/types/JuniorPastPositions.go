package types

import (
	"github.com/shopspring/decimal"
)

type JuniorPastPosition struct {
	ProtocolId             string `json:"protocolId"`
	SmartYieldAddress      string `json:"smartYieldAddress"`
	UnderlyingTokenAddress string `json:"underlyingTokenAddress"`
	UnderlyingTokenSymbol  string `json:"underlyingTokenSymbol"`

	TokensIn        decimal.Decimal `json:"tokensIn"`
	UnderlyingOut   decimal.Decimal `json:"underlyingOut"`
	Forfeits        decimal.Decimal `json:"forfeits"`
	TransactionType string          `json:"transactionType"`

	BlockTimestamp int64  `json:"blockTimestamp"`
	TxHash         string `json:"transactionHash"`
}
