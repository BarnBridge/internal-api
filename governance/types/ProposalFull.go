package types

import (
	"github.com/shopspring/decimal"
)

type ProposalFull struct {
	ProposalBase

	Targets    []string `json:"targets"`
	Values     []string `json:"values"`
	Signatures []string `json:"signatures"`
	Calldatas  []string `json:"calldatas"`

	BlockTimestamp      int64 `json:"blockTimestamp"`
	WarmUpDuration      int64 `json:"warmUpDuration"`
	ActiveDuration      int64 `json:"activeDuration"`
	QueueDuration       int64 `json:"queueDuration"`
	GracePeriodDuration int64 `json:"gracePeriodDuration"`
	AcceptanceThreshold int64 `json:"acceptanceThreshold"`
	MinQuorum           int64 `json:"minQuorum"`

	BondStaked decimal.Decimal `json:"-"`

	History []HistoryEvent `json:"history"`
}
