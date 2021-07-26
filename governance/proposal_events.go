package governance

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/barnbridge/internal-api/governance/types"
	"github.com/barnbridge/internal-api/response"
)

func (g *Governance) HandleProposalEvents(ctx *gin.Context) {
	proposalID, err := getProposalId(ctx)
	if err != nil {
		response.Error(ctx, errors.New("invalid proposalID"))
		return
	}

	rows, err := g.db.Connection().Query(ctx, `
	select proposal_id, caller, event_type, event_data, block_timestamp, tx_hash
	from governance.proposal_events
	where proposal_id = $1
	`, proposalID)
	if err != nil && err != pgx.ErrNoRows {
		response.Error(ctx, err)
		return
	}
	defer rows.Close()

	var events []types.Event
	for rows.Next() {
		var e types.Event

		err := rows.Scan(&e.ProposalID, &e.Caller, &e.EventType, &e.Eta, &e.CreateTime, &e.TxHash)
		if err != nil {
			response.Error(ctx, err)
			return
		}

		events = append(events, e)
	}

	response.OKWithBlock(ctx, g.db, events)
}
