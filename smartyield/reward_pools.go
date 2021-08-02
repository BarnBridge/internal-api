package smartyield

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"

	"github.com/barnbridge/internal-api/query"
	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/smartyield/types"
	"github.com/barnbridge/internal-api/utils"
)

func (h *SmartYield) RewardPools(ctx *gin.Context) {
	builder := query.New()

	protocols := strings.ToLower(ctx.DefaultQuery("originator", "all"))
	underlyingSymbol := strings.ToUpper(ctx.DefaultQuery("underlyingSymbol", "all"))
	underlyingAddress := strings.ToLower(ctx.DefaultQuery("underlyingAddress", "all"))

	builder.Filters.Add("pool_type", types.PoolTypeSingle)
	if protocols != "all" {
		protocolsArray := strings.Split(protocols, ",")
		builder.Filters.Add("p.protocol_id", protocolsArray)
	}

	if underlyingSymbol != "ALL" {
		builder.Filters.Add("upper(p.underlying_symbol)", underlyingSymbol)
	}

	if underlyingAddress != "all" {
		builder.Filters.Add("p.underlying_address", utils.NormalizeAddress(underlyingAddress))
	}

	query, params := builder.Run(`
		select
			   r.pool_address,
			   r.pool_token_address,
			   r.reward_token_addresses,
			   p.underlying_decimals,
			   p.protocol_id,
			   p.underlying_symbol,
			   p.underlying_address
		from smart_yield.reward_pools as r
		inner join smart_yield.pools as p
		on p.pool_address = r.pool_token_address 
		$filters$
		`)

	var pools []types.RewardPool
	rows, err := h.db.Connection().Query(ctx, query, params...)

	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	for rows.Next() {
		var p types.RewardPool
		var rewardTokens []string
		err := rows.Scan(&p.PoolAddress, &p.PoolTokenAddress, &rewardTokens, &p.PoolTokenDecimals, &p.ProtocolID, &p.UnderlyingSymbol, &p.UnderlyingAddress)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		p.RewardTokenAddress = rewardTokens[0]

		pools = append(pools, p)
	}

	response.OK(ctx, pools)
}

func (h *SmartYield) RewardPoolsV2(ctx *gin.Context) {
	builder := query.New()

	protocols := strings.ToLower(ctx.DefaultQuery("originator", "all"))
	underlyingSymbol := strings.ToUpper(ctx.DefaultQuery("underlyingSymbol", "all"))
	underlyingAddress := strings.ToLower(ctx.DefaultQuery("underlyingAddress", "all"))

	if protocols != "all" {
		protocolsArray := strings.Split(protocols, ",")
		builder.Filters.Add("p.protocol_id", protocolsArray)
	}

	if underlyingSymbol != "ALL" {
		builder.Filters.Add("upper(p.underlying_symbol)", underlyingSymbol)
	}

	if underlyingAddress != "all" {
		builder.Filters.Add("p.underlying_address", utils.NormalizeAddress(underlyingAddress))
	}

	query, params := builder.Run(`
		select
			   r.pool_type,
			   r.pool_address,
			   r.pool_token_address,
			   r.reward_token_addresses,
			   p.underlying_decimals,
			   p.protocol_id,
			   p.underlying_symbol,
			   p.underlying_address,
			   p.controller_address
		from smart_yield.reward_pools as r
		inner join smart_yield.pools as p
		on p.pool_address = r.pool_token_address 
		$filters$
	`)
	rows, err := h.db.Connection().Query(ctx, query, params...)

	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	var pools []types.RewardPoolV2
	for rows.Next() {
		var p types.RewardPoolV2
		var rewardTokens []string
		err := rows.Scan(&p.PoolType, &p.PoolAddress, &p.PoolTokenAddress, &rewardTokens, &p.PoolTokenDecimals, &p.ProtocolID, &p.UnderlyingSymbol, &p.UnderlyingAddress, &p.PoolControllerAddress)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		for _, t := range rewardTokens {
			var symbol string
			var decimals int64

			err := h.db.Connection().QueryRow(
				ctx,
				`select symbol, decimals from public.tokens where lower(address) = $1`,
				utils.NormalizeAddress(t),
			).Scan(&symbol, &decimals)
			if err != nil && err != pgx.ErrNoRows {
				response.Error(ctx, err)
				return
			}

			p.RewardTokens = append(p.RewardTokens, types.RewardPoolV2RewardToken{
				Address:  t,
				Symbol:   symbol,
				Decimals: decimals,
			})
		}

		pools = append(pools, p)
	}

	response.OK(ctx, pools)
}
