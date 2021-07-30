package smartyield

import (
	"strings"

	"github.com/barnbridge/internal-api/query"
	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/smartyield/types"
	"github.com/barnbridge/internal-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func (h *SmartYield) RewardPoolsStakingActions(ctx *gin.Context) {
	builder := query.New()

	pool := ctx.Param("poolAddress")

	rewardPoolAddress, err := utils.ValidateAccount(pool)
	if err != nil {
		response.BadRequest(ctx, errors.New("invalid pool address"))
		return
	}

	var underlyingDecimals int64
	err = h.db.Connection().QueryRow(
		ctx,
		`	select underlying_decimals 
					from smart_yield.pools p inner join smart_yield.reward_pools rp
						on p.pool_address = rp.pool_token_address
				where rp.pool_address = $1`,
		rewardPoolAddress,
	).Scan(&underlyingDecimals)
	if err != nil {
		response.Error(ctx, errors.Wrap(err, "could not find smartyield pool"))
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
		from smart_yield.rewards_staking_actions
		$filters$
		order by included_in_block desc ,
				 tx_index desc,
				 log_index desc
		$offset$ $limit$
	`)

	rows, err := h.db.Connection().Query(ctx, query, params...)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	tenPowDec := decimal.NewFromInt(10).Pow(decimal.NewFromInt(underlyingDecimals))

	var transactions []types.StakingAction
	for rows.Next() {
		var t types.StakingAction
		err := rows.Scan(&t.UserAddress, &t.Amount, &t.TransactionType, &t.BlockTimestamp, &t.BlockNumber, &t.TxHash)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		t.Amount = t.Amount.DivRound(tenPowDec, int32(underlyingDecimals))

		transactions = append(transactions, t)
	}

	query, params = builder.Run(`select count(*) from smart_yield.rewards_staking_actions t $filters$`)
	var count int
	err = h.db.Connection().QueryRow(ctx, query, params...).Scan(&count)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.OKWithBlock(ctx, h.db, transactions, response.Meta().Set("count", count))
}
