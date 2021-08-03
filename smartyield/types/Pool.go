package types

type Pool struct {
	ProtocolId         string `json:"protocolId"`
	ControllerAddress  string `json:"controllerAddress"`
	ModelAddress       string `json:"modelAddress"`
	ProviderAddress    string `json:"providerAddress"`
	PoolAddress        string `json:"smartYieldAddress"`
	OracleAddress      string `json:"oracleAddress"`
	JuniorBondAddress  string `json:"juniorBondAddress"`
	SeniorBondAddress  string `json:"seniorBondAddress"`
	CTokenAddress      string `json:"cTokenAddress"`
	UnderlyingAddress  string `json:"underlyingAddress"`
	UnderlyingSymbol   string `json:"underlyingSymbol"`
	UnderlyingDecimals int64  `json:"underlyingDecimals"`
	RewardPoolAddress  string `json:"rewardPoolAddress"`

	State PoolState `json:"state"`
}
