package types

import (
	"github.com/shopspring/decimal"
)

type Voter struct {
	Address             string          `json:"address"`
	BondStaked          decimal.Decimal `json:"bondStaked"`
	LockedUntil         int64           `json:"lockedUntil"`
	DelegatedPower      decimal.Decimal `json:"delegatedPower"`
	Votes               int64           `json:"votes"`
	Proposals           int64           `json:"proposals"`
	VotingPower         decimal.Decimal `json:"votingPower"`
	HasActiveDelegation bool            `json:"hasActiveDelegation"`
}
