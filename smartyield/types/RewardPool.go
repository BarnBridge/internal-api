package types

type RewardPool struct {
	PoolAddress        string `json:"poolAddress"`
	PoolTokenAddress   string `json:"poolTokenAddress"`
	RewardTokenAddress string `json:"rewardTokenAddress"`
	PoolTokenDecimals  int64  `json:"poolTokenDecimals"`
	ProtocolID         string `json:"protocolId"`
	UnderlyingSymbol   string `json:"underlyingSymbol"`
	UnderlyingAddress  string `json:"underlyingAddress"`
}

type RewardPoolV2RewardToken struct {
	Address  string `json:"address"`
	Symbol   string `json:"symbol"`
	Decimals int64  `json:"decimals"`
}

type RewardPoolV2 struct {
	PoolType              string                    `json:"poolType"`
	PoolAddress           string                    `json:"poolAddress"`
	PoolControllerAddress string                    `json:"poolControllerAddress"`
	PoolTokenAddress      string                    `json:"poolTokenAddress"`
	RewardTokens          []RewardPoolV2RewardToken `json:"rewardTokens"`
	PoolTokenDecimals     int64                     `json:"poolTokenDecimals"`
	ProtocolID            string                    `json:"protocolId"`
	UnderlyingSymbol      string                    `json:"underlyingSymbol"`
	UnderlyingAddress     string                    `json:"underlyingAddress"`
}
