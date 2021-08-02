package smartyield

import (
	"time"

	"github.com/barnbridge/internal-api/response"
	"github.com/barnbridge/internal-api/smartyield/types"
	"github.com/barnbridge/internal-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

func (h *SmartYield) UserPortfolioValue(ctx *gin.Context) {
	address := ctx.Param("address")

	userAddress, err := utils.ValidateAccount(address)
	if err != nil {
		response.BadRequest(ctx, errors.New("invalid user address"))
		return
	}

	rows, err := h.db.Connection().Query(ctx, `
		select ts,
		   smart_yield.junior_portfolio_value_at_ts($1, ts),
		   smart_yield.senior_portfolio_value_at_ts($1, ts)
		from 
			generate_series(( select extract(epoch from now() - interval '30 days')::bigint ),
			(select extract(epoch from now()))::bigint, 12 * 60 * 60) as ts
	`,
		userAddress,
	)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	var points []types.PortfolioValuePoint

	for rows.Next() {
		var ts int64
		var junior, senior float64

		err := rows.Scan(&ts, &junior, &senior)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		points = append(points, types.PortfolioValuePoint{
			Timestamp:   time.Unix(ts, 0),
			SeniorValue: &senior,
			JuniorValue: &junior,
		})
	}

	response.OK(ctx, points)
}

func (h *SmartYield) UserSeniorPortfolioValue(ctx *gin.Context) {
	address := ctx.Param("address")

	userAddress, err := utils.ValidateAccount(address)
	if err != nil {
		response.BadRequest(ctx, errors.New("invalid user address"))
		return
	}

	rows, err := h.db.Connection().Query(ctx, `
		select ts,
		   smart_yield.senior_portfolio_value_at_ts($1, ts)
		from 
			generate_series(( select extract(epoch from now() - interval '30 days')::bigint ),
			(select extract(epoch from now()))::bigint, 12 * 60 * 60) as ts
	`,
		userAddress,
	)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	var points []types.PortfolioValuePoint

	for rows.Next() {
		var ts int64
		var senior float64

		err := rows.Scan(&ts, &senior)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		points = append(points, types.PortfolioValuePoint{
			Timestamp:   time.Unix(ts, 0),
			SeniorValue: &senior,
		})
	}

	response.OK(ctx, points)
}

func (h *SmartYield) UserJuniorPortfolioValue(ctx *gin.Context) {
	address := ctx.Param("address")

	userAddress, err := utils.ValidateAccount(address)
	if err != nil {
		response.BadRequest(ctx, errors.New("invalid user address"))
		return
	}

	rows, err := h.db.Connection().Query(ctx, `
		select ts,
		   smart_yield.junior_portfolio_value_at_ts($1, ts)
		from 
			generate_series(( select extract(epoch from now() - interval '30 days')::bigint ),
			(select extract(epoch from now()))::bigint, 12 * 60 * 60) as ts
	`,
		userAddress,
	)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	var points []types.PortfolioValuePoint

	for rows.Next() {
		var ts int64
		var junior float64

		err := rows.Scan(&ts, &junior)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		points = append(points, types.PortfolioValuePoint{
			Timestamp:   time.Unix(ts, 0),
			JuniorValue: &junior,
		})
	}

	response.OK(ctx, points)
}
