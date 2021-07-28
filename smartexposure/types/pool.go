package types

import (
	globalTypes "github.com/barnbridge/internal-api/types"
)

type Pool struct {
	EPoolAddress string
	ProtocolId   string

	TokenA globalTypes.Token
	TokenB globalTypes.Token

	StartAtBlock int64
}
