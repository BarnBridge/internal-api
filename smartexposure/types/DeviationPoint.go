package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type DeviationPoint struct {
	Point     time.Time       `json:"point"`
	Deviation decimal.Decimal `json:"deviation"`
}
