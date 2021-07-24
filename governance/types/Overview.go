package types

import (
	"github.com/shopspring/decimal"
)

type Overview struct {
	AvgLockTimeSeconds     int64           `json:"avgLockTimeSeconds"`
	Holders                int64           `json:"holders"`
	TotalDelegatedPower    decimal.Decimal `json:"totalDelegatedPower"`
	TotalVBond             decimal.Decimal `json:"totalVbond"`
	Voters                 int64           `json:"voters"`
	BarnUsers              int64           `json:"barnUsers"`
	HoldersStakingExcluded int64           `json:"holdersStakingExcluded"`
}
