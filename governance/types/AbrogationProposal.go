package types

import (
	"github.com/shopspring/decimal"
)

type AbrogationProposal struct {
	ProposalID   uint64          `json:"proposalId"`
	Creator      string          `json:"caller"`
	CreateTime   uint64          `json:"createTime"`
	Description  string          `json:"description"`
	ForVotes     decimal.Decimal `json:"forVotes"`
	AgainstVotes decimal.Decimal `json:"againstVotes"`
}
