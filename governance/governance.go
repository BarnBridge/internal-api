package governance

import (
	"github.com/gin-gonic/gin"

	"github.com/barnbridge/internal-api/db"
)

type Governance struct {
	db *db.DB
}

func New(db *db.DB) *Governance {
	return &Governance{db: db}
}

func (g *Governance) SetRoutes(engine *gin.Engine) {
	governance := engine.Group("/api/governance")
	governance.GET("/proposals", g.HandleProposals)
	governance.GET("/proposals/:proposalID", g.HandleProposalDetails)
	governance.GET("/proposals/:proposalID/votes", g.HandleVotes)
	governance.GET("/proposals/:proposalID/events", g.HandleProposalEvents)
	governance.GET("/overview", g.HandleOverview)
	governance.GET("/voters", g.HandleVoters)
	governance.GET("/abrogation-proposals", g.HandleAbrogationProposals)
	governance.GET("/abrogation-proposals/:proposalID", g.HandleAbrogationProposalDetails)
	governance.GET("/abrogation-proposals/:proposalID/votes", g.HandleAbrogationProposalVotes)
	governance.GET("/treasury/transactions", g.HandleTreasuryTxs)
	governance.GET("/treasury/tokens", g.HandleTreasuryTokens)
}
