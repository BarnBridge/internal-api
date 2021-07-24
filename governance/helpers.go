package governance

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/barnbridge/internal-api/governance/types"
)

func getTimeLeft(state types.ProposalState, createTime, warmUpDuration, activeDuration, queueDuration, gracePeriodDuration int64) *int64 {
	now := time.Now().Unix()
	var timeLeft int64

	switch state {
	case types.CANCELED, types.FAILED, types.ACCEPTED, types.EXPIRED, types.EXECUTED, types.ABROGATED:
		return nil
	case types.WARMUP:
		timeLeft = createTime + warmUpDuration - now
	case types.ACTIVE:
		timeLeft = createTime + warmUpDuration + activeDuration - now
	case types.QUEUED:
		timeLeft = createTime + warmUpDuration + activeDuration + queueDuration - now
	case types.GRACE:
		timeLeft = createTime + warmUpDuration + activeDuration + queueDuration + gracePeriodDuration - now
	}

	return &timeLeft
}

func getProposalId(ctx *gin.Context) (int64, error) {
	proposalIDString := ctx.Param("proposalID")
	proposalID, err := strconv.ParseInt(proposalIDString, 10, 64)
	if err != nil {
		return 0, errors.New("invalid proposalID")
	}

	return proposalID, nil
}
