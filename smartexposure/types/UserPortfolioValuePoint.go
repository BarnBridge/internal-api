package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type UserPortfolioValuePoint struct {
	Point          time.Time       `json:"point"`
	PortfolioValue decimal.Decimal `json:"portfolioValueSE,omitempty"`
}
