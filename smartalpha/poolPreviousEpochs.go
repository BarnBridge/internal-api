package smartalpha

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"

	"github.com/barnbridge/internal-api/query"
	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/smartalpha/types"
	"github.com/barnbridge/internal-api/utils"
)

func (s *SmartAlpha) poolPreviousEpochs(ctx *gin.Context) {
	builder := query.New()
	poolAddress := ctx.Param("poolAddress")
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
	}
	builder.Filters.Add("p.pool_address", poolAddress)
	var poolEpochs types.PoolPreviousEpoch
	q, params := builder.WithPaginationFromCtx(ctx).Run(`
		select p.pool_address,
			   p.pool_name,
			   p.pool_token_address,
			   p.pool_token_symbol,
			   p.pool_token_decimals,
			   p.oracle_asset_symbol
		from smart_alpha.pools p
	$filters$
	$offset$ $limit$ `)
	err := s.db.Connection().QueryRow(ctx, q, params...).Scan(&poolEpochs.PoolAddress, &poolEpochs.PoolName, &poolEpochs.PoolToken.Address, &poolEpochs.PoolToken.Symbol,
		&poolEpochs.PoolToken.Decimals, &poolEpochs.OracleAssetSymbol)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	q, params = builder.WithPaginationFromCtx(ctx).Run(`
	select p.epoch_id,
		   p.senior_liquidity,
		   p.junior_liquidity,
		   p.upside_exposure_rate,
		   p.downside_protection_rate,
		   p.epoch_entry_price,
		   (select block_timestamp from smart_alpha.pool_epoch_info where pool_address = p.pool_address and epoch_id < p.epoch_id order by epoch_id desc limit 1),
		   e.block_timestamp as end_date
	from smart_alpha.pool_epoch_info p
			 left join smart_alpha.epoch_end_events e on e.pool_address = p.pool_address and e.epoch_id = p.epoch_id
	$filters$
	order by p.epoch_id desc
	$offset$ $limit$`)

	rows, err := s.db.Connection().Query(ctx, q, params...)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}
	defer rows.Close()
	var epochs []types.Epoch

	for rows.Next() {
		var e types.Epoch
		err := rows.Scan(&e.Id, &e.SeniorLiquidity, &e.JuniorLiquidity, &e.UpsideExposureRate, &e.DownsideProtectionRate, &e.EntryPrice, &e.StartDate, &e.EndDate)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		epochs = append(epochs, e)
	}

	poolEpochs.Epochs = epochs

	response.OKWithBlock(ctx, s.db, poolEpochs, response.Meta().Set("count", len(poolEpochs.Epochs)))
}
