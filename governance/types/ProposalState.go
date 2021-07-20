package types

type ProposalState string

const (
	CREATED   ProposalState = "CREATED"
	WARMUP    ProposalState = "WARMUP"
	ACTIVE    ProposalState = "ACTIVE"
	CANCELED  ProposalState = "CANCELED"
	FAILED    ProposalState = "FAILED"
	ACCEPTED  ProposalState = "ACCEPTED"
	QUEUED    ProposalState = "QUEUED"
	GRACE     ProposalState = "GRACE"
	EXPIRED   ProposalState = "EXPIRED"
	EXECUTED  ProposalState = "EXECUTED"
	ABROGATED ProposalState = "ABROGATED"
)
