package types

type HistoryEvent struct {
	Name    string `json:"name"`
	StartTs int64  `json:"startTimestamp"`
	EndTs   int64  `json:"endTimestamp"`
	TxHash  string `json:"txHash"`
}
