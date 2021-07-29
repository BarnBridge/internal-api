package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type PortfolioValuePoint struct {
	Point          time.Time       `json:"point"`
	PortfolioValue decimal.Decimal `json:"portfolioValueSE,omitempty"`
}
