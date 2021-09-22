package smartexposure

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/shopspring/decimal"

	"github.com/barnbridge/internal-api/query"
	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/smartexposure/types"
	globalTypes "github.com/barnbridge/internal-api/types"
	"github.com/barnbridge/internal-api/utils"
)

func (s *SmartExposure) allTranches(ctx *gin.Context) {
	builder := query.New()
	poolAddress := ctx.DefaultQuery("poolAddress", "")
	if poolAddress != "" {
		poolAddress, err := utils.ValidateAccount(poolAddress)
		if err != nil {
			response.Error(ctx, err)
			return
		}
		err, exists := s.checkPoolExists(ctx, poolAddress)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		if !exists {
			response.NotFound(ctx)
			return
		}

		builder.Filters.Add("t.pool_address", poolAddress)
	}

	q, params := builder.Run(`
		select t.pool_address,
			   t.etoken_address,
			   t.etoken_symbol,
			   t.s_factor_e,
			   t.target_ratio,
			   t.token_a_ratio,
			   t.token_b_ratio,
			   p.token_a_address,
			   p.token_a_decimals,
			   p.token_a_symbol,
			   p.token_b_address,
			   p.token_b_decimals,
			   p.token_b_symbol,
			   coalesce((select price from token_prices where token_address = p.token_a_address and quote_asset = 'USD' order by block_timestamp desc limit 1),0) as token_a_price,
			   coalesce((select price from token_prices where token_address = p.token_b_address and quote_asset = 'USD'  order by block_timestamp desc limit 1),0) as token_b_price
		from smart_exposure.tranches t
				 inner join smart_exposure.pools p on p.pool_address = t.pool_address
		$filters$
	`)

	rows, err := s.db.Connection().Query(ctx, q, params...)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}
	defer rows.Close()

	var tranches []types.Tranche

	for rows.Next() {
		var t types.Tranche
		var tokenAState, tokenBState globalTypes.TokenState
		err = rows.Scan(&t.PoolAddress, &t.ETokenAddress, &t.ETokenSymbol, &t.SFactorE, &t.TargetRatio, &t.TokenARatio, &t.TokenBRatio, &t.TokenA.Address,
			&t.TokenA.Decimals, &t.TokenA.Symbol, &t.TokenB.Address, &t.TokenB.Decimals, &t.TokenB.Symbol, &tokenAState.Price, &tokenBState.Price,
		)
		if err != nil {
			response.Error(ctx, err)
			return
		}
		t.ETokenAddress = utils.NormalizeAddress(t.ETokenAddress)
		t.PoolAddress = utils.NormalizeAddress(t.PoolAddress)

		builder.Filters.Add("etoken_address", t.ETokenAddress)
		q, params = builder.Run(`select token_a_liquidity,
					   token_b_liquidity,
					   etoken_price,
                       token_a_current_ratio,
                       token_b_current_ratio,
					   included_in_block,
					   to_timestamp(block_timestamp)
				from smart_exposure.tranche_state t
				$filters$
				order by block_timestamp desc limit 1`)

		err = s.db.Connection().QueryRow(ctx, q, params...).Scan(&t.State.TokenALiquidity, &t.State.TokenBLiquidity, &t.State.ETokenPrice, &t.State.TokenACurrentRatio, &t.State.TokenBCurrentRatio, &t.State.BlockNumber, &t.State.BlockTimestamp)
		if err != nil && err != pgx.ErrNoRows {
			response.Error(ctx, err)
			return
		}
		tokenAState.BlockNumber = t.State.BlockNumber
		tokenAState.BlockTimestamp = t.State.BlockTimestamp.Unix()
		tokenBState.BlockNumber = t.State.BlockNumber
		tokenBState.BlockTimestamp = t.State.BlockTimestamp.Unix()

		t.TokenA.State = &tokenAState
		t.TokenB.State = &tokenBState

		tranches = append(tranches, t)
	}

	response.OKWithBlock(ctx, s.db, tranches)
}

func (s *SmartExposure) trancheDetails(ctx *gin.Context) {
	eTokenAddress := ctx.Param("eTokenAddress")
	eTokenAddress, err := utils.ValidateAccount(eTokenAddress)
	if err != nil {
		response.BadRequest(ctx, err)
		return
	}

	err, exists := s.checkTrancheExists(ctx, eTokenAddress)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	if !exists {
		response.NotFound(ctx)
		return
	}

	var t types.Tranche
	var tokenAState, tokenBState globalTypes.TokenState
	var rebalancingCondition decimal.Decimal
	err = s.db.Connection().QueryRow(ctx, `select s_factor_e,
					   target_ratio,
					   token_a_ratio,
					   token_a_address,
					   token_a_symbol,
					   token_a_decimals,
					   coalesce(token_a_price_usd,0),
					   token_a_included_in_block,
					   token_a_block_timestamp,
					   token_b_address,
					   coalesce(token_b_price_usd,0),
					   token_b_included_in_block,
					   token_b_block_timestamp,
					   token_b_ratio,
					   token_b_symbol,
					   token_b_decimals,
					   pool_state_rebalancing_interval,
					   pool_state_rebalancing_condition,
					   pool_state_last_rebalance,
					   tranche_state_token_a_liquidity,
					   tranche_state_token_b_liquidity,
					   tranche_state_e_token_price,
					   tranche_state_current_ratio,
					   tranche_state_token_a_current_ratio,
					   tranche_state_token_b_current_ratio,
					   tranche_state_included_in_block,
					   to_timestamp(tranche_state_block_timestamp) from smart_exposure.get_tranche_details($1)`, eTokenAddress).Scan(&t.SFactorE, &t.TargetRatio, &t.TokenARatio, &t.TokenA.Address, &t.TokenA.Symbol,
		&t.TokenA.Decimals, &tokenAState.Price, &tokenAState.BlockNumber, &tokenAState.BlockTimestamp, &t.TokenB.Address, &tokenBState.Price, &tokenBState.BlockNumber,
		&tokenBState.BlockTimestamp, &t.TokenBRatio, &t.TokenB.Symbol, &t.TokenB.Decimals, &t.RebalancingInterval, &rebalancingCondition, &t.State.LastRebalance, &t.State.TokenALiquidity,
		&t.State.TokenBLiquidity, &t.State.ETokenPrice, &t.State.CurrentRatio, &t.State.TokenACurrentRatio, &t.State.TokenBCurrentRatio, &t.State.BlockNumber, &t.State.BlockTimestamp)

	if err != nil {
		response.Error(ctx, err)
		return
	}
	t.TokenA.State = &tokenAState
	t.TokenB.State = &tokenBState
	t.RebalancingCondition = rebalancingCondition.String()
	t.ETokenAddress = utils.NormalizeAddress(eTokenAddress)
	response.OKWithBlock(ctx, s.db, t)
}
