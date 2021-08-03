package smartyield

import (
	"strings"

	"github.com/barnbridge/internal-api/query"
	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/smartyield/types"
	"github.com/barnbridge/internal-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func (h *SmartYield) UserTransactionHistory(ctx *gin.Context) {
	builder := query.New()

	address := ctx.Param("address")

	userAddress, err := utils.ValidateAccount(address)
	if err != nil {
		response.BadRequest(ctx, errors.New("invalid user address"))
		return
	}

	builder.Filters.Add("user_address", userAddress)

	originator := strings.ToLower(ctx.DefaultQuery("originator", "all"))
	if originator != "all" {
		if !isSupportedOriginator(originator) {
			response.BadRequest(ctx, errors.New("invalid originator parameter"))
			return
		}

		builder.Filters.Add("protocol_id", originator)
	}

	token := strings.ToLower(ctx.DefaultQuery("token", "all"))
	if token != "all" {
		tokenAddress, err := utils.ValidateAccount(token)
		if err != nil {
			response.BadRequest(ctx, errors.New("invalid token address"))
			return
		}

		builder.Filters.Add("underlying_token_address", tokenAddress)
	}

	txType := strings.ToUpper(ctx.DefaultQuery("transactionType", "all"))
	if txType != "ALL" {
		if !isSupportedTxType(txType) {
			response.BadRequest(ctx, errors.New("invalid transactionType parameter"))
			return
		}

		builder.Filters.Add("transaction_type", txType)
	}

	query, params := builder.WithPaginationFromCtx(ctx).Run(`
		select 
			h.protocol_id,
			h.pool_address,
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
		var hist types.PoolHistory
		var underlyingDecimals int64

		err := rows.Scan(&hist.ProtocolId, &hist.Pool, &hist.UnderlyingTokenAddress, &underlyingDecimals, &hist.UnderlyingTokenSymbol, &hist.Amount,
			&hist.Tranche, &hist.TransactionType, &hist.TransactionHash, &hist.BlockTimestamp, &hist.BlockNumber,
		)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		tenPowDec := decimal.NewFromInt(10).Pow(decimal.NewFromInt(underlyingDecimals))
		hist.Amount = hist.Amount.DivRound(tenPowDec, int32(underlyingDecimals))

		history = append(history, hist)
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
