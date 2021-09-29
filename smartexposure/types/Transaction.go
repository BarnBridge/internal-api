package types

import (
	"github.com/shopspring/decimal"

	globalTypes "github.com/barnbridge/internal-api/types"
)

type Transaction struct {
	ETokenAddress   string            `json:"eTokenAddress"`
	AccountAddress  string            `json:"accountAddress"`
	TokenA          globalTypes.Token `json:"tokenA"`
	TokenB          globalTypes.Token `json:"tokenB"`
	AmountA         decimal.Decimal   `json:"amountA"`
	AmountB         decimal.Decimal   `json:"amountB"`
	AmountEToken    decimal.Decimal   `json:"amountEToken"`
	TransactionType string            `json:"transactionType"`
	TransactionHash string            `json:"transactionHash"`
	BlockTimestamp  int64             `json:"blockTimestamp"`
	BlockNumber     int64             `json:"blockNumber"`
	SFactorE        decimal.Decimal   `json:"sFactorE"`
	TokenAPrice     decimal.Decimal   `json:"tokenAPrice"`
	TokenBPrice     decimal.Decimal   `json:"tokenBPrice"`
	ETokenPrice     decimal.Decimal   `json:"eTokenPrice"`
	ETokenSymbol    string            `json:"eTokenSymbol"`
}
