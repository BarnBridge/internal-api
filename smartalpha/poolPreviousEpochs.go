package smartalpha

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

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

	cursorStr := ctx.DefaultQuery("cursor", "-2")
	cursor, err := strconv.ParseInt(cursorStr, 0, 64)
	if err != nil {
		response.Error(ctx, errors.New("invalid parameter 'cursor'"))
		return
	}
	if cursor < 0 {
		cursor = math.MaxInt32
	}

	limit, err := utils.GetQueryLimit(ctx)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	var dirOperator string
	var orderDirection string
	direction := strings.ToLower(ctx.DefaultQuery("direction", "down"))
	switch direction {
	case "up":
		dirOperator = ">="
		orderDirection = "asc"
	case "down":
		dirOperator = "<="
		orderDirection = "desc"
	default:
		response.Error(ctx, errors.New("invalid parameter 'direction'"))
		return
	}

	var pool types.PoolDetails
	err = s.db.Connection().QueryRow(ctx, `
		select p.pool_address,
			   p.pool_name,
			   p.pool_token_address,
			   p.pool_token_symbol,
			   p.pool_token_decimals,
			   p.oracle_asset_symbol
		from smart_alpha.pools p
		where p.pool_address = $1`, poolAddress).Scan(&pool.PoolAddress, &pool.PoolName, &pool.PoolToken.Address, &pool.PoolToken.Symbol,
		&pool.PoolToken.Decimals, &pool.OracleAssetSymbol)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	} else if err == pgx.ErrNoRows {
		response.NotFound(ctx)
		return
	}

	rows, err := s.db.Connection().Query(ctx,
		fmt.Sprintf(`select * from (
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
			where p.pool_address = $1 and p.epoch_id %s $2
			order by p.epoch_id %s
			limit $3) x order by x.epoch_id desc
	`, dirOperator, orderDirection), poolAddress, cursor, limit)
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
		if strings.ToUpper(pool.OracleAssetSymbol) == "ETH" {
			priceDecimals = 18
		} else {
			priceDecimals = 8
		}

		e.SeniorLiquidity = e.SeniorLiquidity.Shift(-int32(pool.PoolToken.Decimals))
		e.JuniorLiquidity = e.JuniorLiquidity.Shift(-int32(pool.PoolToken.Decimals))
		e.JuniorProfits = e.JuniorProfits.Shift(-int32(pool.PoolToken.Decimals))
		e.SeniorProfits = e.SeniorProfits.Shift(-int32(pool.PoolToken.Decimals))
		e.JuniorTokenPriceStart = e.JuniorTokenPriceStart.Shift(-18)
		e.SeniorTokenPriceStart = e.SeniorTokenPriceStart.Shift(-18)
		e.EntryPrice = e.EntryPrice.Shift(-priceDecimals)

		epochs = append(epochs, e)
	}

	if len(epochs) == 0 {
		response.OKWithBlock(ctx, s.db, epochs, response.Meta().Set("hasNewer", false).Set("hasOlder", false))
		return
	}

	var hasNext, hasPrev bool
	err = s.db.Connection().QueryRow(ctx,
		`
	select count(*) > 0
		from smart_alpha.pool_epoch_info p
				 inner join smart_alpha.epoch_end_events e
						   on e.pool_address = p.pool_address and e.epoch_id = p.epoch_id 
			where p.pool_address = $1 and p.epoch_id < $2
	`,
		poolAddress, epochs[len(epochs)-1].Id).Scan(&hasPrev)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	err = s.db.Connection().QueryRow(ctx,
		`
	select count(*) > 0
		from smart_alpha.pool_epoch_info p
				 inner join smart_alpha.epoch_end_events e
						   on e.pool_address = p.pool_address and e.epoch_id = p.epoch_id 
			where p.pool_address = $1 and p.epoch_id > $2
	`,
		poolAddress, epochs[0].Id).Scan(&hasNext)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.OKWithBlock(ctx, s.db, epochs, response.Meta().Set("hasNewer", hasNext).Set("hasOlder", hasPrev))
}
