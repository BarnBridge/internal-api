package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type Chart map[time.Time]Aggregate

type Aggregate struct {
	SumDeposits    decimal.Decimal `json:"sumDeposits"`
	SumWithdrawals decimal.Decimal `json:"sumWithdrawals"`
}
