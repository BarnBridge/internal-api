package types

import "github.com/shopspring/decimal"

type Vote struct {
	User           string          `json:"address"`
	Support        bool            `json:"support"`
	BlockTimestamp int64           `json:"blockTimestamp"`
	Power          decimal.Decimal `json:"power"`
}
