package smartyield

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/barnbridge/internal-api/query"
	"github.com/barnbridge/internal-api/smartyield/types"
	"github.com/barnbridge/internal-api/utils"

	"github.com/barnbridge/internal-api/response"
)

func (h *SmartYield) UserSeniorRedeems(ctx *gin.Context) {
	builder := query.New()

	address := ctx.Param("address")

	userAddress, err := utils.ValidateAccount(address)
	if err != nil {
		response.BadRequest(ctx, errors.New("invalid user address"))
		return
	}

	builder.Filters.Add("owner_address", userAddress)

	originator := strings.ToLower(ctx.DefaultQuery("originator", "all"))
	if originator != "all" {
		if !isSupportedOriginator(originator) {
			response.BadRequest(ctx, errors.New("invalid originator parameter"))
			return
		}

		builder.Filters.Add("(select protocol_id from smart_yield.pools p where p.pool_address = r.pool_address)", originator)
	}

	token := strings.ToLower(ctx.DefaultQuery("token", "all"))
	if token != "all" {
		tokenAddress, err := utils.ValidateAccount(token)
		if err != nil {
			response.BadRequest(ctx, errors.New("invalid token address"))
			return
		}

		builder.Filters.Add("(select underlying_address from smart_yield.pools p where p.pool_address = r.pool_address)", tokenAddress)
	}

	query, params := builder.WithPaginationFromCtx(ctx).Run(`
		select 
			r.pool_address,
			(select underlying_decimals from smart_yield.pools p where r.pool_address = p.pool_address) as underlying_token_decimals,
			r.owner_address,
			r.senior_bond_address,
			r.fee,
			r.block_timestamp,
			r.senior_bond_id,
			r.tx_hash,
			b.underlying_in,
			b.gain,
			b.for_days
		from smart_yield.senior_redeem_events as r
			inner join smart_yield.senior_entry_events as b
					on r.senior_bond_address = b.senior_bond_address and r.senior_bond_id = b.senior_bond_id
		$filters$
		order by r.included_in_block desc, r.tx_index desc, r.log_index desc
		$offset$ $limit$
	`)

	rows, err := h.db.Connection().Query(ctx, query, params...)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	var seniorBondRedeems []types.SeniorRedeem
	for rows.Next() {
		var redeem types.SeniorRedeem
		var underlyingDecimals int64

		err := rows.Scan(&redeem.PoolAddress, &underlyingDecimals, &redeem.UserAddress, &redeem.SeniorBondAddress, &redeem.Fee, &redeem.BlockTimestamp, &redeem.SeniorBondID, &redeem.TxHash, &redeem.UnderlyingIn, &redeem.Gain, &redeem.ForDays)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		tenPowDec := decimal.NewFromInt(10).Pow(decimal.NewFromInt(underlyingDecimals))

		redeem.UnderlyingIn = redeem.UnderlyingIn.DivRound(tenPowDec, int32(underlyingDecimals))
		redeem.Fee = redeem.Fee.DivRound(tenPowDec, int32(underlyingDecimals))
		redeem.Gain = redeem.Gain.DivRound(tenPowDec, int32(underlyingDecimals))

		seniorBondRedeems = append(seniorBondRedeems, redeem)
	}

	query, params = builder.Run(`
		select count(*)
		from smart_yield.senior_redeem_events as r
				 inner join smart_yield.senior_entry_events as b
							on r.senior_bond_address = b.senior_bond_address and r.senior_bond_id = b.senior_bond_id
		$filters$
	`)
	var count int
	err = h.db.Connection().QueryRow(ctx, query, params...).Scan(&count)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.OKWithBlock(ctx, h.db, seniorBondRedeems, response.Meta().Set("count", count))
}
