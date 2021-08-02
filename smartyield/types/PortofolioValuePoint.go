package types

import (
	"time"
)

type PortfolioValuePoint struct {
	Timestamp   time.Time `json:"timestamp"`
	SeniorValue *float64  `json:"seniorValue,omitempty"`
	JuniorValue *float64  `json:"juniorValue,omitempty"`
}
