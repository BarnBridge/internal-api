package smartalpha

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/smartalpha/types"
	"github.com/barnbridge/internal-api/utils"
)

func (s *SmartAlpha) UserQueuePositions(ctx *gin.Context) {
	userAddress := ctx.Param("address")
	userAddress, err := utils.ValidateAccount(userAddress)
	if err != nil {
		response.Error(ctx, errors.Wrap(err, "invalid user address"))
		return
	}

	includeActiveVal := ctx.DefaultQuery("includeActive", "false")
	var includeActive bool
	switch includeActiveVal {
	case "true":
		includeActive = true
	case "false":
		includeActive = false
	default:
		response.Error(ctx, errors.New("invalid parameter `includeActive`"))
		return
	}

	rows, err := s.db.Connection().Query(ctx, `
				select distinct on (e.pool_address, e.tranche) e.pool_address,
					   p.pool_name,
					   p.pool_token_address,
					   p.pool_token_symbol,
					   p.pool_token_decimals,
					   p.oracle_asset_symbol,
					   e.tranche,
					   e.block_timestamp,
					   'entry'
				from smart_alpha.user_join_entry_queue_events e
						 left join smart_alpha.pools p on e.pool_address = p.pool_address
				where e.user_address = $1 and (select count(*) from smart_alpha.user_redeem_tokens_events re where re.pool_address = e.pool_address and re.user_address = e.user_address and e.epoch_id = re.epoch_id) = 0
				union
				select distinct on(x.pool_address, x.tranche) x.pool_address,
					   p2.pool_name,
					   p2.pool_token_address,
					   p2.pool_token_symbol,
					   p2.pool_token_decimals,
					   p2.oracle_asset_symbol,
					   x.tranche,
					   x.block_timestamp,
					   'exit'
				from smart_alpha.user_join_exit_queue_events x
						 left join smart_alpha.pools p2 on x.pool_address = p2.pool_address
				where x.user_address = $1 and (select count(*) from smart_alpha.user_redeem_underlying_events ru where ru.pool_address = x.pool_address and ru.user_address = x.user_address and x.epoch_id = ru.epoch_id) = 0
				order by block_timestamp desc`, userAddress)

	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	defer rows.Close()

	var userQueuePositions = make([]types.UserQueuePosition, 0)
	for rows.Next() {
		var u types.UserQueuePosition
		err := rows.Scan(&u.PoolAddress, &u.PoolName, &u.PoolToken.Address, &u.PoolToken.Symbol, &u.PoolToken.Decimals, &u.OracleAssetSymbol, &u.Tranche, &u.BlockTimestamp, &u.QueueType)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		userQueuePositions = append(userQueuePositions, u)
	}

	if includeActive {
		rows, err := s.db.Connection().Query(ctx, `
			select distinct on (pool_address) coalesce(x1.pool_address, x2.pool_address) as pool_address,
											  coalesce(x1.pool_name, x2.pool_name),
											  coalesce(x1.pool_token_address, x2.pool_token_address),
											  coalesce(x1.pool_token_symbol, x2.pool_token_symbol),
											  coalesce(x1.pool_token_decimals, x2.pool_token_decimals),
											  coalesce(x1.oracle_asset_symbol, x2.oracle_asset_symbol)
			from public.erc20_balances_at_ts($1,
											 ( select array_agg(junior_token_address) || array_agg(senior_token_address)
											   from smart_alpha.pools ), ( select extract(epoch from now()) )::bigint)
					 left join smart_alpha.pools x1 on x1.senior_token_address = token_address
					 left join smart_alpha.pools x2 on x2.junior_token_address = token_address
			where balance > 0;
		`, userAddress)
		if err != nil && err != pgx.ErrNoRows {
			response.Error(ctx, err)
			return
		}

		defer rows.Close()

		for rows.Next() {
			var u types.UserQueuePosition
			err := rows.Scan(&u.PoolAddress, &u.PoolName, &u.PoolToken.Address, &u.PoolToken.Symbol, &u.PoolToken.Decimals, &u.OracleAssetSymbol)
			if err != nil {
				response.Error(ctx, err)
				return
			}

			userQueuePositions = appendQueuePositionIfNotExists(userQueuePositions, u)
		}
	}

	response.OK(ctx, userQueuePositions)
}

func appendQueuePositionIfNotExists(positions []types.UserQueuePosition, newPosition types.UserQueuePosition) []types.UserQueuePosition {
	for _, v := range positions {
		if v.PoolName == newPosition.PoolName {
			return positions
		}
	}

	return append(positions, newPosition)
}
