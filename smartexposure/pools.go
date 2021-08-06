package smartexposure

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"

	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/smartexposure/types"
	"github.com/barnbridge/internal-api/utils"
)

func (s *SmartExposure) sePools(ctx *gin.Context) {
	rows, err := s.db.Connection().Query(ctx, `
			select pool_address,
				   pool_name,
				   token_a_address,
				   token_a_symbol,
				   token_a_decimals,
				   token_b_address,
				   token_b_symbol,
				   token_b_decimals
				   from smart_exposure.pools`)

	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	var pools []types.Pool

	for rows.Next() {
		var p types.Pool
		err := rows.Scan(&p.PoolAddress, &p.PoolName, &p.TokenA.Address, &p.TokenA.Symbol, &p.TokenA.Decimals,
			&p.TokenB.Address, &p.TokenB.Symbol, &p.TokenB.Decimals)

		if err != nil {
			response.Error(ctx, err)
			return
		}
		p.PoolAddress = utils.NormalizeAddress(p.PoolAddress)
		p.TokenA.Address = utils.NormalizeAddress(p.TokenA.Address)
		p.TokenB.Address = utils.NormalizeAddress(p.TokenB.Address)

		var state types.PoolState
		err = s.db.Connection().QueryRow(ctx,
			`select pool_liquidity, 
						   last_rebalance,
						   rebalancing_interval,
						   rebalancing_condition,
						   included_in_block,
						   to_timestamp(block_timestamp)
					from smart_exposure.pool_state
					where pool_address = $1
					order by block_timestamp desc limit 1
					`, p.PoolAddress).Scan(&state.PoolLiquidity, &state.LastRebalance, &state.RebalancingInterval,
			&state.RebalancingCondition, &state.BlockNumber, &state.BlockTimestamp)

		if err != nil {
			response.Error(ctx, err)
			return
		}

		p.State = state

		pools = append(pools, p)
	}

	response.OK(ctx, pools)
}
