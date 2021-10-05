package smartalpha

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/barnbridge/internal-api/query"
	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/smartalpha/types"
	globalTypes "github.com/barnbridge/internal-api/types"
	"github.com/barnbridge/internal-api/utils"
)

func (s *SmartAlpha) RewardPools(ctx *gin.Context) {
	builder := query.New()

	pool := ctx.DefaultQuery("poolAddress", "")
	if pool != "" {
		rewardPoolAddress, err := utils.ValidateAccount(pool)
		if err != nil {
			response.BadRequest(ctx, errors.New("invalid pool address"))
			return
		}
		builder.Filters.Add("r.pool_address", rewardPoolAddress)
	}

	query, params := builder.Run(`
		select
			   r.pool_type,
			   r.pool_address,
			   r.pool_token_address,
			   r.reward_token_addresses,
			   t.decimals,
			   t.symbol
		from smart_alpha.reward_pools as r
		inner join public.tokens as t
		on t.address = r.pool_token_address 
		$filters$
	`)
	rows, err := s.db.Connection().Query(ctx, query, params...)

	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	var pools []types.RewardPool
	for rows.Next() {
		var p types.RewardPool
		var rewardTokens []string
		err := rows.Scan(&p.PoolType, &p.PoolAddress, &p.PoolToken.Address, &rewardTokens, &p.PoolToken.Decimals, &p.PoolToken.Symbol)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		for _, t := range rewardTokens {
			var symbol string
			var decimals int64

			err := s.db.Connection().QueryRow(
				ctx,
				`select symbol, decimals from public.tokens where lower(address) = $1`,
				utils.NormalizeAddress(t),
			).Scan(&symbol, &decimals)
			if err != nil && err != pgx.ErrNoRows {
				response.Error(ctx, err)
				return
			}

			p.RewardTokens = append(p.RewardTokens, globalTypes.Token{
				Address:  t,
				Symbol:   symbol,
				Decimals: decimals,
			})
		}

		pools = append(pools, p)
	}

	response.OK(ctx, pools)
}

func (s *SmartAlpha) RewardPoolTransactions(ctx *gin.Context) {
	builder := query.New()

	pool := ctx.Param("poolAddress")

	rewardPoolAddress, err := utils.ValidateAccount(pool)
	if err != nil {
		response.BadRequest(ctx, errors.New("invalid pool address"))
		return
	}
	builder.Filters.Add("pool_address", rewardPoolAddress)

	userAddress := ctx.DefaultQuery("userAddress", "all")
	if userAddress != "all" {
		builder.Filters.Add("user_address", utils.NormalizeAddress(userAddress))
	}

	transactionType := strings.ToUpper(ctx.DefaultQuery("transactionType", "all"))
	if transactionType != "ALL" {
		if !checkRewardPoolTxType(transactionType) {
			response.Error(ctx, errors.New("transaction type does not exist"))
			return
		}

		builder.Filters.Add("action_type", transactionType)
	}

	query, params := builder.WithPaginationFromCtx(ctx).Run(`
		select
			user_address,
			amount,
			action_type,
			block_timestamp,
			included_in_block,
			tx_hash
		from smart_alpha.rewards_staking_actions
		$filters$
		order by included_in_block desc ,
				 tx_index desc,
				 log_index desc
		$offset$ $limit$
	`)

	rows, err := s.db.Connection().Query(ctx, query, params...)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	var underlyingDecimals int64
	err = s.db.Connection().QueryRow(ctx, `select decimals from public.tokens where address = (select pool_token_address from smart_alpha.reward_pools where pool_address = $1)`, rewardPoolAddress).Scan(&underlyingDecimals)

	tenPowDec := decimal.NewFromInt(10).Pow(decimal.NewFromInt(underlyingDecimals))

	var transactions []types.RewardPoolTransaction
	for rows.Next() {
		var t types.RewardPoolTransaction
		err := rows.Scan(&t.UserAddress, &t.Amount, &t.TransactionType, &t.BlockTimestamp, &t.BlockNumber, &t.TxHash)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		t.Amount = t.Amount.DivRound(tenPowDec, int32(underlyingDecimals))

		transactions = append(transactions, t)
	}

	query, params = builder.Run(`select count(*) from smart_alpha.rewards_staking_actions t $filters$`)
	var count int
	err = s.db.Connection().QueryRow(ctx, query, params...).Scan(&count)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.OKWithBlock(ctx, s.db, transactions, response.Meta().Set("count", count))
}

func checkRewardPoolTxType(action string) bool {
	txType := [2]string{"DEPOSIT", "WITHDRAW"}
	for _, tx := range txType {
		if action == tx {
			return true
		}
	}

	return false
}
