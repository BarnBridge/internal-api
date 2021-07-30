package smartyield

import (
	"strings"

	"github.com/barnbridge/internal-api/query"
	"github.com/barnbridge/internal-api/smartyield/types"
	"github.com/barnbridge/internal-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/barnbridge/internal-api/response"
)

func (h *SmartYield) JuniorPastPositions(ctx *gin.Context) {
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

	query, params := builder.WithPaginationFromCtx(ctx).Run(`
		select h.protocol_id as protocol_id,
			   	h.pool_address,
			   	underlying_token_address,
				( select underlying_decimals from smart_yield.pools p where h.pool_address = p.pool_address) as underlying_token_decimals,
			   	( select underlying_symbol from smart_yield.pools p where h.pool_address = p.pool_address ) as underlying_token_symbol,
			    coalesce(
				   (select tokens_in from smart_yield.junior_instant_withdraw_events s where s.tx_hash = h.tx_hash and s.log_index = h.log_index ),
				   (select tokens_in from smart_yield.junior_2step_redeem_events r inner join smart_yield.junior_2step_withdraw_events jwe on r.junior_bond_address = jwe.junior_bond_address and r.junior_bond_id = jwe.junior_bond_id where r.tx_hash = h.tx_hash and r.log_index = h.log_index )
			    ) as tokens_in,
			    coalesce(
				   ( select underlying_out from smart_yield.junior_instant_withdraw_events s where s.tx_hash = h.tx_hash and s.log_index = h.log_index ),
				   ( select underlying_out from smart_yield.junior_2step_redeem_events r where r.tx_hash = h.tx_hash and r.log_index = h.log_index )
			    ) as underlying_out,
			    coalesce(
				   ( select forfeits from smart_yield.junior_instant_withdraw_events s where s.tx_hash = h.tx_hash and s.log_index = h.log_index),
				   0
			    ) as forfeits,
			    transaction_type,
			    tx_hash,
			    block_timestamp
		from smart_yield.transaction_history h
		$filters$
			and transaction_type = ANY(ARRAY['JUNIOR_INSTANT_WITHDRAW'::sy_tx_history_tx_type, 'JUNIOR_REDEEM'::sy_tx_history_tx_type])
		order by included_in_block desc, tx_index desc, log_index desc
		$offset$ $limit$
	`)

	rows, err := h.db.Connection().Query(ctx, query, params...)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	var positions []types.JuniorPastPosition
	for rows.Next() {
		var p types.JuniorPastPosition
		var underlyingDecimals int64

		err := rows.Scan(&p.ProtocolId, &p.SmartYieldAddress, &p.UnderlyingTokenAddress, &underlyingDecimals, &p.UnderlyingTokenSymbol,
			&p.TokensIn, &p.UnderlyingOut, &p.Forfeits, &p.TransactionType, &p.TxHash, &p.BlockTimestamp,
		)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		tenPowDec := decimal.NewFromInt(10).Pow(decimal.NewFromInt(underlyingDecimals))

		p.TokensIn = p.TokensIn.DivRound(tenPowDec, int32(underlyingDecimals))
		p.UnderlyingOut = p.UnderlyingOut.DivRound(tenPowDec, int32(underlyingDecimals))
		p.Forfeits = p.Forfeits.DivRound(tenPowDec, int32(underlyingDecimals))

		positions = append(positions, p)
	}

	query, params = builder.Run(`
		select count(*) 
			from smart_yield.transaction_history as h
		$filters$
			and transaction_type = ANY(ARRAY['JUNIOR_INSTANT_WITHDRAW'::sy_tx_history_tx_type, 'JUNIOR_REDEEM'::sy_tx_history_tx_type])
	`)
	var count int
	err = h.db.Connection().QueryRow(ctx, query, params...).Scan(&count)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.OKWithBlock(ctx, h.db, positions, response.Meta().Set("count", count))
}
