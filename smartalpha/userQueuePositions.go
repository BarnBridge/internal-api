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

	response.OK(ctx, userQueuePositions)
}
