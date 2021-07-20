package types

import (
	"github.com/shopspring/decimal"
)

type ProposalBase struct {
	Id          uint64 `json:"proposalId"`
	Proposer    string `json:"proposer"`
	Description string `json:"description"`
	Title       string `json:"title"`
	CreateTime  int64  `json:"createTime"`

	State         ProposalState   `json:"state"`
	StateTimeLeft *int64          `json:"stateTimeLeft"`
	ForVotes      decimal.Decimal `json:"forVotes"`
	AgainstVotes  decimal.Decimal `json:"againstVotes"`
}
