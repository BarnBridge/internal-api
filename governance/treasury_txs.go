package governance

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/shopspring/decimal"

	"github.com/barnbridge/internal-api/governance/types"
	"github.com/barnbridge/internal-api/query"
	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/utils"
)

func (g *Governance) HandleTreasuryTxs(ctx *gin.Context) {
	builder := query.New()

	treasuryAddress, err := utils.ValidateAccount(ctx.DefaultQuery("address", ""))
	if err != nil {
		response.BadRequest(ctx, err)
		return
	}
	builder.Filters.Add("t.account", treasuryAddress)

	tokenAddress := strings.ToLower(ctx.DefaultQuery("tokenAddress", "all"))
	if tokenAddress != "all" {
		builder.Filters.Add("t.token_address", tokenAddress)
	}

	txDirection := strings.ToUpper(ctx.DefaultQuery("transactionDirection", "all"))
	if txDirection != "ALL" {
		builder.Filters.Add("tx_direction", txDirection)
	}

	q, params := builder.WithPaginationFromCtx(ctx).Run(`
		select t.token_address,
			   t.account,
			   t.counterparty,
			   t.amount,
			   t.tx_direction,
			   t.block_timestamp,
			   t.included_in_block,
			   t.tx_hash,
			   e20t.symbol,
			   e20t.decimals,
			   coalesce(( select label from labels where address = t.account ), '')      as accountLabel,
			   coalesce(( select label from labels where address = t.counterparty ), '') as counterpartyLabel
		from account_erc20_transfers as t
				 inner join tokens e20t on t.token_address = e20t.address
		$filters$
		order by included_in_block desc, t.tx_index desc, t.log_index desc
		$offset$ $limit$
	`)

	rows, err := g.db.Connection().Query(ctx, q, params...)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	var list []types.TreasuryTx
	for rows.Next() {
		var t types.TreasuryTx

		err := rows.Scan(&t.TokenAddress, &t.AccountAddress, &t.CounterpartyAddress, &t.Amount, &t.TransactionDirection, &t.BlockTimestamp, &t.BlockNumber, &t.TransactionHash, &t.TokenSymbol, &t.TokenDecimals, &t.AccountLabel, &t.CounterpartyLabel)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		t.Amount = t.Amount.DivRound(decimal.NewFromInt(10).Pow(decimal.NewFromInt(t.TokenDecimals)), int32(t.TokenDecimals))
		list = append(list, t)
	}

	q, params = builder.Run(`select count(*) from account_erc20_transfers t $filters$`)

	var count int
	err = g.db.Connection().QueryRow(ctx, q, params...).Scan(&count)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.OKWithBlock(ctx, g.db, list, response.Meta().Set("count", count))
}
