package governance

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/barnbridge/internal-api/governance/types"
	"github.com/barnbridge/internal-api/response"
)

func (g *Governance) HandleAbrogationProposalDetails(ctx *gin.Context) {
	proposalID, err := getProposalId(ctx)
	if err != nil {
		response.Error(ctx, errors.New("invalid proposalID"))
		return
	}

	var ap types.AbrogationProposal
	err = g.db.Connection().QueryRow(ctx, `
	select proposal_id, creator, create_time, description ,
		   coalesce(( select sum(power) from governance.abrogation_proposal_votes(proposal_id) where support = true ), 0) as for_votes,
		   coalesce(( select sum(power) from governance.abrogation_proposal_votes(proposal_id) where support = false ), 0) as against_votes
	from governance.abrogation_proposals 
	where proposal_id = $1
	`, proposalID).Scan(&ap.ProposalID, &ap.Creator, &ap.CreateTime, &ap.Description, &ap.ForVotes, &ap.AgainstVotes)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}

	if err == pgx.ErrNoRows {
		response.NotFound(ctx)
		return
	}

	response.OKWithBlock(ctx, g.db, ap)
}
