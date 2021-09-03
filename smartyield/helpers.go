package smartyield

import (
	"strings"

	"github.com/barnbridge/internal-api/smartyield/types"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func checkRewardPoolTxType(action string) bool {
	txType := [2]string{"JUNIOR_STAKE", "JUNIOR_UNSTAKE"}
	for _, tx := range txType {
		if action == tx {
			return true
		}
	}

	return false
}

func validateWindow(window string) (string, string, error) {
	if window == "24h" {
		return "'24 hours'", "'minute'", nil
	}

	if window == "1w" {
		return "'7 days'", "'hour'", nil
	}

	if window == "30d" {
		return "'30 days'", "'day'", nil
	}

	return "", "", errors.New("invalid window")
}

func isSupportedTxType(t string) bool {
	switch types.TxType(strings.ToUpper(t)) {
	case types.TxJuniorDeposit, types.TxJuniorInstantWithdraw, types.TxJuniorRegularWithdraw,
		types.TxJuniorRedeem, types.TxSeniorDeposit, types.TxSeniorRedeem, types.TxJtokenSend, types.TxJtokenReceive,
		types.TxJbondSend, types.TxJbondReceive, types.TxSbondSend, types.TxSbondReceive,
		types.TxJuniorStake, types.TxJuniorUnstake:
		return true
	}

	return false
}

func getSortForSeniorBonds(ctx *gin.Context, direction string) (string, error) {
	sort := ctx.DefaultQuery("sort", "")

	if sort != "maturityDate" && sort != "depositedAmount" && sort != "redeemableAmount" && sort != "" {
		return "", errors.New("invalid sort")
	} else if sort != "" {
		return sort + " " + direction + ", ", nil
	}

	return "", nil
}

func getSortForJuniorBonds(ctx *gin.Context, direction string) (string, error) {
	sort := ctx.DefaultQuery("sort", "")

	if sort != "maturityDate" && sort != "depositedAmount" && sort != "" {
		return "", errors.New("invalid sort")
	} else if sort != "" {
		return sort + " " + direction + ", ", nil
	}

	return "", nil
}

func isSupportedOriginator(originator string) bool {
	switch strings.ToLower(originator) {

	// TODO this should be dynamic
	case "compound/v2":
		return true
	case "cream/v2":
		return true
	case "aave/v2":
		return true
	}

	return false
}
