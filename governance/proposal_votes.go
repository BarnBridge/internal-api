package governance

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/barnbridge/internal-api/governance/types"
	"github.com/barnbridge/internal-api/query"
	"github.com/barnbridge/internal-api/response"
)

func (g *Governance) HandleVotes(ctx *gin.Context) {
	proposalID, err := getProposalId(ctx)
	if err != nil {
		response.Error(ctx, errors.New("invalid proposalID"))
		return
	}

	builder := query.New()

	err = builder.SetLimitFromCtx(ctx)
	if err != nil {
		response.BadRequest(ctx, err)
		return
	}

	err = builder.SetOffsetFromCtx(ctx)
	if err != nil {
		response.BadRequest(ctx, err)
		return
	}

	support := strings.ToLower(ctx.DefaultQuery("support", ""))
	if support != "" {
		if support != "true" && support != "false" {
			response.BadRequest(ctx, errors.New("wrong value for support parameter"))
			return
		}
		builder.Filters.Add("support", support)
	}

	q, params := builder.WithPaginationFromCtx(ctx).Run(`
	select user_id, support, block_timestamp, power from governance.proposal_votes($param_overwrite$)
	$filters$
	order by power desc
	$offset$ $limit$
	`)

	params = append(params, proposalID)
	q = strings.Replace(q, "$param_overwrite$", fmt.Sprintf("$%d", len(params)), 1)

	rows, err := g.db.Connection().Query(ctx, q, params...)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}
	defer rows.Close()

	var votes []types.Vote

	for rows.Next() {
		var v types.Vote

		err := rows.Scan(&v.User, &v.Support, &v.BlockTimestamp, &v.Power)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		votes = append(votes, v)
	}

	q, params = builder.Run(`
	select count(*) from governance.proposal_votes($param_overwrite$)
	$filters$
	`)

	params = append(params, proposalID)
	q = strings.Replace(q, "$param_overwrite$", fmt.Sprintf("$%d", len(params)), 1)

	var count int
	err = g.db.Connection().QueryRow(ctx, q, params...).Scan(&count)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.OKWithBlock(ctx, g.db, votes, response.Meta().Set("count", count))
}
