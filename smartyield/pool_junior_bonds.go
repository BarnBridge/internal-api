package smartyield

import (
	"fmt"
	"strings"

	"github.com/barnbridge/internal-api/query"
	"github.com/barnbridge/internal-api/smartyield/types"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/utils"
)

func (h *SmartYield) PoolJuniorBonds(ctx *gin.Context) {
	builder := query.New()

	pool := ctx.Param("address")

	poolAddress, err := utils.ValidateAccount(pool)
	if err != nil {
		response.BadRequest(ctx, errors.New("invalid pool address"))
		return
	}

	underlyingDecimals, err := h.PoolUnderlyingDecimals(ctx, poolAddress)
	if err != nil {
		response.Error(ctx, errors.Wrap(err, "could not find smartyield pool"))
		return
	}
	tenPowDec := decimal.NewFromInt(10).Pow(decimal.NewFromInt(underlyingDecimals))

	builder.Filters.Add("b.pool_address", poolAddress)

	sortDirection := strings.ToLower(ctx.DefaultQuery("sortDirection", "desc"))
	if sortDirection != "desc" && sortDirection != "asc" {
		response.Error(ctx, errors.New("invalid sort direction"))
		return
	}

	sort, err := getSortForJuniorBonds(ctx, sortDirection)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	redeemed := strings.ToLower(ctx.DefaultQuery("redeemed", ""))
	if redeemed == "true" {
		builder.Filters.Add(
			"(select count(*) from smart_yield.junior_2step_redeem_events r where r.junior_bond_id = b.junior_bond_id and r.junior_bond_address = b.junior_bond_address)",
			"0",
			">",
		)
	} else if redeemed == "false" {
		builder.Filters.Add(
			"(select count(*) from smart_yield.junior_2step_redeem_events r where r.junior_bond_id = b.junior_bond_id and r.junior_bond_address = b.junior_bond_address)",
			"0",
		)
	} else if redeemed != "" {
		response.Error(ctx, errors.New("invalid redeem option"))
		return
	}

	q := `select b.buyer_address,
			   b.junior_bond_id,
			   tokens_in as depositedAmount,
			   matures_at as maturityDate,
			   (select count(*)
				from smart_yield.junior_2step_redeem_events r
				where r.junior_bond_id = b.junior_bond_id
				  and r.junior_bond_address = b.junior_bond_address) > 0 as redeemed,
			   p.underlying_address,
			   p.underlying_symbol,
			   p.underlying_decimals,
			   b.tx_hash,
			   b.block_timestamp
			from smart_yield.junior_2step_withdraw_events b
				inner join smart_yield.pools p on b.pool_address = p.pool_address
			$filters$
		order by %s b.included_in_block desc, b.tx_index desc, b.log_index desc
		$offset$ $limit$
	`

	q = fmt.Sprintf(q, sort)

	query, params := builder.WithPaginationFromCtx(ctx).Run(q)

	rows, err := h.db.Connection().Query(ctx, query, params...)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	var users []types.JuniorBondUser
	for rows.Next() {
		var u types.JuniorBondUser
		err := rows.Scan(&u.AccountAddress, &u.JuniorBondId, &u.DepositedAmount, &u.MaturityDate, &u.Redeemed,
			&u.UnderlyingTokenAddress, &u.UnderlyingTokenSymbol, &u.UnderlyingTokenDecimals,
			&u.TxHash, &u.BlockTimestamp,
		)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		u.DepositedAmount = u.DepositedAmount.DivRound(tenPowDec, int32(u.UnderlyingTokenDecimals))

		users = append(users, u)
	}

	query, params = builder.Run(`select count(b.buyer_address) from smart_yield.junior_2step_withdraw_events b $filters$`)
	var count int
	err = h.db.Connection().QueryRow(ctx, query, params...).Scan(&count)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.OKWithBlock(ctx, h.db, users, response.Meta().Set("count", count))
}
