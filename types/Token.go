package types

type Token struct {
	TokenAddress  string `json:"tokenAddress"`
	TokenSymbol   string `json:"tokenSymbol"`
	TokenDecimals int64  `json:"tokenDecimals"`
}
