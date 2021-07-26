package governance

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"

	"github.com/barnbridge/internal-api/governance/types"
	"github.com/barnbridge/internal-api/query"
	"github.com/barnbridge/internal-api/response"
)

func (g *Governance) HandleAbrogationProposals(ctx *gin.Context) {
	builder := query.New()

	q, params := builder.WithPaginationFromCtx(ctx).Run(`
	select proposal_id, creator, create_time
	from governance.abrogation_proposals
	order by proposal_id desc
	$offset$ $limit$
	`)

	rows, err := g.db.Connection().Query(ctx, q, params...)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}
	defer rows.Close()

	var list []types.AbrogationProposal
	for rows.Next() {
		var a types.AbrogationProposal

		err := rows.Scan(&a.ProposalID, &a.Creator, &a.CreateTime)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		list = append(list, a)
	}

	var count int
	err = g.db.Connection().QueryRow(ctx, `select count(*) from governance.abrogation_proposals`).Scan(&count)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.OKWithBlock(ctx, g.db, list, response.Meta().Set("count", count))
}
