package types

import (
	"github.com/barnbridge/internal-api/types"
)

type RewardPool struct {
	PoolType     string        `json:"poolType"`
	PoolAddress  string        `json:"poolAddress"`
	PoolToken    types.Token   `json:"poolToken"`
	RewardTokens []types.Token `json:"rewardTokens"`
}
