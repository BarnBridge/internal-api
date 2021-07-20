package governance

import (
	"time"

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
