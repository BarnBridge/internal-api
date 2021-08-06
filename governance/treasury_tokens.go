package governance

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"

	"github.com/barnbridge/internal-api/response"
	globalTypes "github.com/barnbridge/internal-api/types"
	"github.com/barnbridge/internal-api/utils"
)

func (g *Governance) HandleTreasuryTokens(ctx *gin.Context) {
	treasuryAddress, err := utils.ValidateAccount(ctx.DefaultQuery("address", ""))
	if err != nil {
		response.BadRequest(ctx, err)
		return
	}

	rows, err := g.db.Connection().Query(ctx, `
		select distinct transfers.token_address, tokens.symbol, tokens.decimals
		from account_erc20_transfers transfers
				 inner join tokens on transfers.token_address = tokens.address
		where account = $1;
	`, treasuryAddress)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	var tokens []globalTypes.Token
	for rows.Next() {
		var t globalTypes.Token
		err := rows.Scan(&t.Address, &t.Symbol, &t.Decimals)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		tokens = append(tokens, t)
	}

	response.OKWithBlock(ctx, g.db, tokens)
}
