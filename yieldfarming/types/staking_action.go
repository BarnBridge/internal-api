package types

import (
	"github.com/shopspring/decimal"
)

type StakingAction struct {
	UserAddress     string          `json:"userAddress"`
	TokenAddress    string          `json:"tokenAddress"`
	Amount          decimal.Decimal `json:"amount"`
	TransactionHash string          `json:"transactionHash"`
	ActionType      string          `json:"actionType"`
	BlockTimestamp  int64           `json:"blockTimestamp"`
}
