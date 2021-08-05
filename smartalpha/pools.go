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

func (s *SmartAlpha) Pools(ctx *gin.Context) {
	builder := query.New()
	poolAddress := strings.ToLower(ctx.DefaultQuery("poolAddress", "all"))
	if poolAddress != "all" {
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

		builder.Filters.Add("p.pool_address", poolAddress)
	}

	q, params := builder.Run(`
					select p.pool_address,
						   p.pool_name,
						   p.pool_token_address,
						   p.pool_token_symbol,
						   p.pool_token_decimals,
						   p.junior_token_address,
						   p.senior_token_address,
						   p.oracle_address,
						   p.oracle_asset_symbol,
						   p.epoch1_start,
						   p.epoch_duration,
						   i.epoch_id,
						   i.senior_liquidity,
						   i.junior_liquidity,
						   i.upside_exposure_rate,
						   i.downside_protection_rate
					from smart_alpha.pools p
							 inner join smart_alpha.pool_epoch_info i on i.pool_address = p.pool_address
					$filters$
					order by epoch_id desc
					limit 1`)

	rows, err := s.db.Connection().Query(ctx, q, params...)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	defer rows.Close()
	var pools []types.Pool

	for rows.Next() {
		var p types.Pool
		err = rows.Scan(&p.PoolAddress, &p.PoolName, &p.PoolToken.TokenAddress, &p.PoolToken.TokenSymbol, &p.PoolToken.TokenDecimals, &p.JuniorTokenAddress,
			&p.SeniorTokenAddress, &p.OracleAddress, &p.OracleAssetSymbol, &p.Epoch1Start, &p.EpochDuration, &p.State.Epoch, &p.State.SeniorLiquidity, &p.State.JuniorLiquidity,
			&p.State.UpsideExposureRate, &p.State.DownsideProtectionRate)

		if err != nil {
			response.Error(ctx, err)
			return
		}

		pools = append(pools, p)
	}

	response.OK(ctx, pools)
}
