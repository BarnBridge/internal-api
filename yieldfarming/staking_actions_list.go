package yieldfarming

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/barnbridge/internal-api/query"
	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/utils"
	"github.com/barnbridge/internal-api/yieldfarming/types"
)

func (h *YieldFarming) StakingActionsList(ctx *gin.Context) {
	builder := query.New()

	userAddress := ctx.DefaultQuery("userAddress", "all")
	if userAddress != "all" {
		ua, err := utils.ValidateAccount(userAddress)
		if err != nil {
			response.BadRequest(ctx, err)
			return
		}
		builder.Filters.Add("user_address", ua)
	}

	actionType := strings.ToUpper(ctx.DefaultQuery("actionType", "all"))
	if actionType != "ALL" {
		if !checkTxType(actionType) {
			response.Error(ctx, errors.New("action type does not exist"))
			return
		}

		builder.Filters.Add("action_type", actionType)
	}

	tokenAddress := strings.ToLower(ctx.DefaultQuery("tokenAddress", "all"))
	if tokenAddress != "all" {
		builder.Filters.Add("token_address", tokenAddress)
	}

	q, params := builder.WithPaginationFromCtx(ctx).Run(`
		select
			tx_hash,
			user_address,
			token_address,
			amount ,
			action_type,
			block_timestamp
		from yield_farming.transactions
		$filters$
		order by included_in_block desc, tx_index desc, log_index desc
		$offset$ $limit$
	`)

	rows, err := h.db.Connection().Query(ctx, q, params...)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	var stakingActions []types.StakingAction
	for rows.Next() {
		var sa types.StakingAction
		err := rows.Scan(&sa.TransactionHash, &sa.UserAddress, &sa.TokenAddress, &sa.Amount, &sa.ActionType, &sa.BlockTimestamp)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		stakingActions = append(stakingActions, sa)
	}

	q, params = builder.Run(`select count(*) from yield_farming.transactions t $filters$`)
	var count int
	err = h.db.Connection().QueryRow(ctx, q, params...).Scan(&count)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.OKWithBlock(ctx, h.db, stakingActions, response.Meta().Set("count", count))
}
