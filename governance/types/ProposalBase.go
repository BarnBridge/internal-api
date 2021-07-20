package types

type ProposalBase struct {
	Id          uint64 `json:"proposalId"`
	Proposer    string `json:"proposer"`
	Description string `json:"description"`
	Title       string `json:"title"`
	CreateTime  int64  `json:"createTime"`

	State         ProposalState `json:"state"`
	StateTimeLeft *int64        `json:"stateTimeLeft"`
	ForVotes      string        `json:"forVotes"`
	AgainstVotes  string        `json:"againstVotes"`
}
