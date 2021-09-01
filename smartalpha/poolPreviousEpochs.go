package smartalpha

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"

	"github.com/barnbridge/internal-api/query"
	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/smartalpha/types"
	"github.com/barnbridge/internal-api/utils"
)

func (s *SmartAlpha) poolPreviousEpochs(ctx *gin.Context) {
	poolAddress := ctx.Param("poolAddress")
	poolAddress, err := utils.ValidateAccount(poolAddress)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	var poolEpochs types.PoolPreviousEpoch

	err = s.db.Connection().QueryRow(ctx, `
		select p.pool_address,
			   p.pool_name,
			   p.pool_token_address,
			   p.pool_token_symbol,
			   p.pool_token_decimals,
			   p.oracle_asset_symbol
		from smart_alpha.pools p
		where p.pool_address = $1`, poolAddress).Scan(&poolEpochs.PoolAddress, &poolEpochs.PoolName, &poolEpochs.PoolToken.Address, &poolEpochs.PoolToken.Symbol,
		&poolEpochs.PoolToken.Decimals, &poolEpochs.OracleAssetSymbol)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	if err == pgx.ErrNoRows {
		response.NotFound(ctx)
		return
	}

	builder := query.New()
	builder.Filters.Add("p.pool_address", poolAddress)

	q, params := builder.WithPaginationFromCtx(ctx).Run(`
		select p.epoch_id,
			   p.senior_liquidity,
			   p.junior_liquidity,
			   p.upside_exposure_rate,
			   p.downside_protection_rate,
			   coalesce(p.epoch_entry_price, 0),
			   ( select block_timestamp + 1
				 from smart_alpha.epoch_end_events
				 where pool_address = p.pool_address
				   and epoch_id < p.epoch_id
				 order by epoch_id desc
				 limit 1 )       as start_date,
			   e.block_timestamp as end_date,
	           e.junior_profits,
	           e.senior_profits,
	           p.junior_token_price_start,
	           p.senior_token_price_start
		from smart_alpha.pool_epoch_info p
				 inner join smart_alpha.epoch_end_events e
						   on e.pool_address = p.pool_address and e.epoch_id = p.epoch_id 
			$filters$
			order by p.epoch_id desc
			$offset$ $limit$
	`)

	rows, err := s.db.Connection().Query(ctx, q, params...)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}
	defer rows.Close()

	var epochs = make([]types.Epoch, 0)

	for rows.Next() {
		var e types.Epoch
		err := rows.Scan(&e.Id,
			&e.SeniorLiquidity, &e.JuniorLiquidity,
			&e.UpsideExposureRate, &e.DownsideProtectionRate,
			&e.EntryPrice,
			&e.StartDate, &e.EndDate,
			&e.JuniorProfits, &e.SeniorProfits,
			&e.JuniorTokenPriceStart, &e.SeniorTokenPriceStart,
		)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		var priceDecimals int32
		if strings.ToUpper(poolEpochs.OracleAssetSymbol) == "ETH" {
			priceDecimals = 18
		} else {
			priceDecimals = 8
		}

		e.SeniorLiquidity = e.SeniorLiquidity.Shift(-int32(poolEpochs.PoolToken.Decimals))
		e.JuniorLiquidity = e.JuniorLiquidity.Shift(-int32(poolEpochs.PoolToken.Decimals))
		e.JuniorProfits = e.JuniorProfits.Shift(-int32(poolEpochs.PoolToken.Decimals))
		e.SeniorProfits = e.SeniorProfits.Shift(-int32(poolEpochs.PoolToken.Decimals))
		e.JuniorTokenPriceStart = e.JuniorTokenPriceStart.Shift(-18)
		e.SeniorTokenPriceStart = e.SeniorTokenPriceStart.Shift(-18)
		e.EntryPrice = e.EntryPrice.Shift(-priceDecimals)

		epochs = append(epochs, e)
	}

	poolEpochs.Epochs = epochs

	q, params = builder.Run(`
		select count(*)
		from smart_alpha.pool_epoch_info p
				 inner join smart_alpha.epoch_end_events e
						   on e.pool_address = p.pool_address and e.epoch_id = p.epoch_id 
			$filters$
	`)
	var count int64

	err = s.db.Connection().QueryRow(ctx, q, params...).Scan(&count)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.OKWithBlock(ctx, s.db, poolEpochs, response.Meta().Set("count", count))
}
