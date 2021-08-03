package smartyield

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/barnbridge/internal-api/query"
	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/smartyield/types"
	"github.com/barnbridge/internal-api/utils"
)

func (h *SmartYield) Pools(ctx *gin.Context) {
	builder := query.New()

	protocols := strings.ToLower(ctx.DefaultQuery("originator", "all"))
	underlyingSymbol := strings.ToUpper(ctx.DefaultQuery("underlyingSymbol", "all"))

	if protocols != "all" {
		protocolsArray := strings.Split(protocols, ",")
		builder.Filters.Add("protocol_id", protocolsArray)
	}

	if underlyingSymbol != "ALL" {
		builder.Filters.Add("upper(underlying_symbol)", underlyingSymbol)
	}

	query, params := builder.Run(`
		select protocol_id,
			   pool_address,
			   controller_address,
			   model_address,
			   provider_address,
			   oracle_address,
			   junior_bond_address,
			   senior_bond_address,
			   receipt_token_address,
			   underlying_address,
			   underlying_symbol,
			   underlying_decimals,
			   coalesce((select pool_address from smart_yield.reward_pools where pool_token_address = p.pool_address ), '') as reward_pool_address
		from smart_yield.pools p
		$filters$
		`)

	rows, err := h.db.Connection().Query(ctx, query, params...)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	tenPow18 := decimal.NewFromInt(10).Pow(decimal.NewFromInt(18))

	var pools []types.Pool
	for rows.Next() {
		var p types.Pool

		err := rows.Scan(&p.ProtocolId, &p.PoolAddress, &p.ControllerAddress, &p.ModelAddress, &p.ProviderAddress, &p.OracleAddress, &p.JuniorBondAddress, &p.SeniorBondAddress, &p.CTokenAddress, &p.UnderlyingAddress, &p.UnderlyingSymbol, &p.UnderlyingDecimals, &p.RewardPoolAddress)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		var state types.PoolState
		err = h.db.Connection().QueryRow(ctx, `
			select included_in_block,
				   block_timestamp,
				   senior_liquidity,
				   junior_liquidity,
				   jtoken_price,
				   senior_apy,
				   junior_apy,
				   originator_apy,
				   originator_net_apy,
				   smart_yield.number_of_seniors(pool_address)                      as number_of_seniors,
				   smart_yield.number_of_active_juniors(pool_address)               as number_of_juniors,
			       smart_yield.number_of_juniors_locked(pool_address)               as number_of_juniors_locked,
				   (abond_matures_at - (select extract (epoch from now())))::double precision / (60*60*24)     as avg_senior_buy,
				   coalesce(smart_yield.junior_liquidity_locked(pool_address), 0)   as junior_liquidity_locked
			from smart_yield.pool_state
			where pool_address = $1
			order by block_timestamp desc
			limit 1;
		`, p.PoolAddress).Scan(
			&state.BlockNumber, &state.BlockTimestamp,
			&state.SeniorLiquidity, &state.JuniorLiquidity, &state.JTokenPrice, &state.SeniorAPY, &state.JuniorAPY,
			&state.OriginatorApy, &state.OriginatorNetApy,
			&state.NumberOfSeniors, &state.NumberOfJuniors, &state.NumberOfJuniorsLocked, &state.AvgSeniorMaturityDays, &state.JuniorLiquidityLocked,
		)
		if err != nil && err != pgx.ErrNoRows {
			response.Error(ctx, err)
			return
		}

		tenPowDec := decimal.NewFromInt(10).Pow(decimal.NewFromInt(p.UnderlyingDecimals))

		state.JuniorLiquidityLocked = state.JuniorLiquidityLocked.Div(tenPowDec)
		state.JTokenPrice = state.JTokenPrice.DivRound(tenPow18, 18)
		state.SeniorLiquidity = state.SeniorLiquidity.Div(tenPowDec)
		state.JuniorLiquidity = state.JuniorLiquidity.Div(tenPowDec)

		p.State = state

		pools = append(pools, p)
	}

	response.OK(ctx, pools)
}

func (h *SmartYield) PoolDetails(ctx *gin.Context) {
	pool := ctx.Param("address")

	poolAddress, err := utils.ValidateAccount(pool)
	if err != nil {
		response.BadRequest(ctx, errors.New("invalid pool address"))
		return
	}

	var p types.Pool

	err = h.db.Connection().QueryRow(ctx, `
		select protocol_id,
			   pool_address,
			   controller_address,
			   model_address,
			   provider_address,
			   oracle_address,
			   junior_bond_address,
			   senior_bond_address,
			   receipt_token_address,
			   underlying_address,
			   underlying_symbol,
			   underlying_decimals,
			   coalesce((select pool_address from smart_yield.reward_pools where pool_token_address = p.pool_address ), '') as reward_pool_address
		from smart_yield.pools p
		where pool_address = $1
	`, poolAddress).Scan(
		&p.ProtocolId,
		&p.PoolAddress, &p.ControllerAddress, &p.ModelAddress, &p.ProviderAddress, &p.OracleAddress, &p.JuniorBondAddress, &p.SeniorBondAddress, &p.CTokenAddress,
		&p.UnderlyingAddress, &p.UnderlyingSymbol, &p.UnderlyingDecimals, &p.RewardPoolAddress,
	)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}
	if err == pgx.ErrNoRows {
		response.NotFound(ctx)
		return
	}

	tenPow18 := decimal.NewFromInt(10).Pow(decimal.NewFromInt(18))

	var state types.PoolState
	err = h.db.Connection().QueryRow(ctx, `
			select included_in_block,
				   block_timestamp,
				   senior_liquidity,
				   junior_liquidity,
				   jtoken_price,
				   senior_apy,
				   junior_apy,
				   originator_apy,
				   originator_net_apy,
				   smart_yield.number_of_seniors(pool_address)                      as number_of_seniors,
				   smart_yield.number_of_active_juniors(pool_address)               as number_of_juniors,
			       smart_yield.number_of_juniors_locked(pool_address)               as number_of_juniors_locked,
				    (abond_matures_at - (select extract (epoch from now())))::double precision / (60*60*24)     as avg_senior_buy,
				   coalesce(smart_yield.junior_liquidity_locked(pool_address), 0)   as junior_liquidity_locked
			from smart_yield.pool_state
			where pool_address = $1
			order by block_timestamp desc
			limit 1;
		`, p.PoolAddress).Scan(
		&state.BlockNumber, &state.BlockTimestamp,
		&state.SeniorLiquidity, &state.JuniorLiquidity, &state.JTokenPrice, &state.SeniorAPY, &state.JuniorAPY,
		&state.OriginatorApy, &state.OriginatorNetApy,
		&state.NumberOfSeniors, &state.NumberOfJuniors, &state.NumberOfJuniorsLocked, &state.AvgSeniorMaturityDays, &state.JuniorLiquidityLocked,
	)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	tenPowDec := decimal.NewFromInt(10).Pow(decimal.NewFromInt(p.UnderlyingDecimals))

	state.JuniorLiquidityLocked = state.JuniorLiquidityLocked.Div(tenPowDec)
	state.JTokenPrice = state.JTokenPrice.DivRound(tenPow18, 18)
	state.SeniorLiquidity = state.SeniorLiquidity.Div(tenPowDec)
	state.JuniorLiquidity = state.JuniorLiquidity.Div(tenPowDec)

	p.State = state

	response.OK(ctx, p)
}
