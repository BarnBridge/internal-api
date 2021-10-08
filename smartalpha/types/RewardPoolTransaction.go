package types

import (
	"github.com/shopspring/decimal"
)

type RewardPoolTransaction struct {
	UserAddress     string          `json:"userAddress"`
	TransactionType string          `json:"transactionType"`
	Amount          decimal.Decimal `json:"amount"`
	BlockTimestamp  int64           `json:"blockTimestamp"`
	BlockNumber     int64           `json:"blockNumber"`
	TxHash          string          `json:"transactionHash"`
}
