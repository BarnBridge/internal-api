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
	governance.GET("/proposals", g.AllProposalsHandler)
	// governance.GET("/proposals/:proposalID", a.ProposalDetailsHandler)
	// governance.GET("/proposals/:proposalID/votes", a.VotesHandler)
	// governance.GET("/proposals/:proposalID/events", a.handleProposalEvents)
	// governance.GET("/overview", a.BondOverview)
	// governance.GET("/voters", a.handleVoters)
	// governance.GET("/abrogation-proposals", a.AllAbrogationProposals)
	// governance.GET("/abrogation-proposals/:proposalID", a.AbrogationProposalDetailsHandler)
	// governance.GET("/abrogation-proposals/:proposalID/votes", a.AbrogationVotesHandler)
	// governance.GET("/treasury/transactions", a.handleTreasuryTxs)
	// governance.GET("/treasury/tokens", a.handleTreasuryTokens)
}
