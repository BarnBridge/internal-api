package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type UserPortfolioPoint struct {
	Point           time.Time       `json:"point"`
	JuniorValue     decimal.Decimal `json:"juniorValue"`
	SeniorValue     decimal.Decimal `json:"seniorValue"`
	EntryQueueValue decimal.Decimal `json:"entryQueueValue"`
	ExitQueueValue  decimal.Decimal `json:"exitQueueValue"`
}
