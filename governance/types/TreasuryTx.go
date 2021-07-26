package types

import "github.com/shopspring/decimal"

type TreasuryTx struct {
	AccountAddress string `json:"accountAddress"`
	AccountLabel   string `json:"accountLabel"`

	CounterpartyAddress string `json:"counterpartyAddress"`
	CounterpartyLabel   string `json:"counterpartyLabel"`

	Amount               decimal.Decimal `json:"amount"`
	TransactionDirection string          `json:"transactionDirection"`
	TokenAddress         string          `json:"tokenAddress"`
	TokenSymbol          string          `json:"tokenSymbol"`
	TokenDecimals        int64           `json:"tokenDecimals"`

	TransactionHash string `json:"transactionHash"`
	BlockTimestamp  int64  `json:"blockTimestamp"`
	BlockNumber     int64  `json:"blockNumber"`
}
