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
						   p.senior_rate_model_address,
						   p.accounting_model_address,
						   p.epoch1_start,
						   p.epoch_duration
					from smart_alpha.pools p
					$filters$`)

	rows, err := s.db.Connection().Query(ctx, q, params...)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	defer rows.Close()
	var pools []types.Pool

	for rows.Next() {
		var p types.Pool

		err = rows.Scan(&p.PoolAddress, &p.PoolName, &p.PoolToken.Address, &p.PoolToken.Symbol, &p.PoolToken.Decimals, &p.JuniorTokenAddress,
			&p.SeniorTokenAddress, &p.OracleAddress, &p.OracleAssetSymbol, &p.SeniorRateModelAddress, &p.AccountingModelAddress, &p.Epoch1Start, &p.EpochDuration)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		err = s.db.Connection().QueryRow(
			ctx,
			`
			select epoch_id, senior_liquidity, junior_liquidity, upside_exposure_rate, downside_protection_rate
			from smart_alpha.pool_epoch_info
			where pool_address = $1
			order by block_timestamp desc
			limit 1;
			`,
			p.PoolAddress,
		).Scan(
			&p.State.Epoch,
			&p.State.SeniorLiquidity,
			&p.State.JuniorLiquidity,
			&p.State.UpsideExposureRate,
			&p.State.DownsideProtectionRate,
		)

		p.State.SeniorLiquidity = p.State.SeniorLiquidity.Shift(-int32(p.PoolToken.Decimals))
		p.State.JuniorLiquidity = p.State.JuniorLiquidity.Shift(-int32(p.PoolToken.Decimals))

		pools = append(pools, p)
	}

	response.OK(ctx, pools)
}
