package smartyield

import (
	"strings"

	"github.com/barnbridge/internal-api/query"
	"github.com/barnbridge/internal-api/smartyield/types"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/utils"
)

func (h *SmartYield) PoolTransactions(ctx *gin.Context) {
	builder := query.New()

	pool := ctx.Param("address")

	poolAddress, err := utils.ValidateAccount(pool)
	if err != nil {
		response.BadRequest(ctx, errors.New("invalid pool address"))
		return
	}

	builder.Filters.Add("pool_address", poolAddress)

	txType := strings.ToUpper(ctx.DefaultQuery("transactionType", "all"))
	if txType != "ALL" {
		if !isSupportedTxType(txType) {
			response.BadRequest(ctx, errors.New("invalid transactionType parameter"))
			return
		}

		builder.Filters.Add("transaction_type", txType)
	}

	query, params := builder.WithPaginationFromCtx(ctx).Run(`
		select h.protocol_id,
			   h.pool_address,
               user_address, 
			   underlying_token_address,
               (select underlying_decimals from smart_yield.pools p where h.pool_address = p.pool_address) as underlying_token_decimals, 
               (select underlying_symbol from smart_yield.pools p where h.pool_address = p.pool_address) as underlying_token_symbol, 
			   amount,
			   tranche,
			   transaction_type,
			   tx_hash,
			   block_timestamp,
			   included_in_block
		from smart_yield.transaction_history h
		$filters$
		order by included_in_block desc, tx_index desc, log_index desc
		$offset$ $limit$
	`)

	rows, err := h.db.Connection().Query(ctx, query, params...)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	var history []types.PoolHistory
	for rows.Next() {
		var ph types.PoolHistory
		var underlyingDecimals int64
		var userAddr string
		err := rows.Scan(&ph.ProtocolId, &ph.Pool, &userAddr,
			&ph.UnderlyingTokenAddress, &underlyingDecimals, &ph.UnderlyingTokenSymbol, &ph.Amount, &ph.Tranche,
			&ph.TransactionType, &ph.TransactionHash, &ph.BlockTimestamp, &ph.BlockNumber,
		)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		ph.AccountAddress = &userAddr
		tenPowDec := decimal.NewFromInt(10).Pow(decimal.NewFromInt(underlyingDecimals))
		ph.Amount = ph.Amount.DivRound(tenPowDec, int32(underlyingDecimals))

		history = append(history, ph)
	}

	query, params = builder.Run(`select count(*) from smart_yield.transaction_history $filters$`)
	var count int
	err = h.db.Connection().QueryRow(ctx, query, params...).Scan(&count)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.OKWithBlock(ctx, h.db, history, response.Meta().Set("count", count))
}
