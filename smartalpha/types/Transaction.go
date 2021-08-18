package types

import (
	"github.com/shopspring/decimal"
)

type Transaction struct {
	PoolAddress        string          `json:"poolAddress"`
	UserAddress        string          `json:"userAddress"`
	Tranche            string          `json:"tranche"`
	TransactionType    string          `json:"transactionType"`
	TokenSymbol        string          `json:"tokenSymbol"`
	Amount             decimal.Decimal `json:"amount"`
	AmountInQuoteAsset decimal.Decimal `json:"amountInQuoteAsset"`
	AmountInUSD        decimal.Decimal `json:"amountInUSD"`
	TransactionHash    string          `json:"transactionHash"`
	BlockTimestamp     int64           `json:"blockTimestamp"`
}
