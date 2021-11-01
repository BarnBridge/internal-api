package smartalpha

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

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
			   p.epoch_duration,
			   coalesce(tvl.epoch_junior_tvl,0),
			   coalesce(tvl.epoch_senior_tvl,0),
			   coalesce(tvl.junior_entry_queue_tvl,0),
			   coalesce(tvl.senior_entry_queue_tvl,0),
               coalesce(tvl.junior_exit_queue_tvl,0),
               coalesce(tvl.senior_exit_queue_tvl,0),
			   coalesce(tvl.junior_exited_tvl,0),
			   coalesce(tvl.senior_exited_tvl,0)
		from smart_alpha.pools p
			left join smart_alpha.pool_tvl_v2(p.pool_address) tvl on 1=1
		$filters$
		order by p.start_at_block, p.pool_name
	`)

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
			&p.SeniorTokenAddress, &p.OracleAddress, &p.OracleAssetSymbol, &p.SeniorRateModelAddress, &p.AccountingModelAddress, &p.Epoch1Start, &p.EpochDuration,
			&p.TVL.EpochJuniorTVL, &p.TVL.EpochSeniorTVL, &p.TVL.JuniorEntryQueueTVL, &p.TVL.SeniorEntryQueueTVL, &p.TVL.JuniorExitQueueTVL, &p.TVL.SeniorExitQueueTVL, &p.TVL.JuniorExitedTVL, &p.TVL.SeniorExitedTVL)
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

	userAddress := ctx.DefaultQuery("userAddress", "")
	if userAddress != "" {
		userAddress, err := utils.ValidateAccount(userAddress)
		if err != nil {
			response.Error(ctx, errors.Wrap(err, "invalid user address"))
			return
		}

		for i := range pools {
			x := false
			pools[i].UserHasActivePosition = &x
		}

		rows, err := s.db.Connection().Query(ctx, `
			select distinct on (pool_address) coalesce(x1.pool_address, x2.pool_address) as pool_address
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
			var poolAddr string

			err := rows.Scan(&poolAddr)
			if err != nil {
				response.Error(ctx, err)
				return
			}

			for i := range pools {
				if utils.NormalizeAddress(pools[i].PoolAddress) == utils.NormalizeAddress(poolAddr) {
					x := true
					pools[i].UserHasActivePosition = &x
					break
				}
			}
		}
	}

	response.OK(ctx, pools)
}
