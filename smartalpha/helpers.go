package smartalpha

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"
)

func (s *SmartAlpha) checkPoolExists(ctx context.Context, addr string) (error, bool) {
	var exists bool

	err := s.db.Connection().QueryRow(ctx, `select (select count(*) from smart_alpha.pools where pool_address = $1) > 0`, addr).Scan(&exists)

	return err, exists
}

func checkTxType(action string) bool {
	txType := []string{"JUNIOR_ENTRY", "JUNIOR_REDEEM_TOKENS", "JUNIOR_EXIT", "JUNIOR_REDEEM_UNDERLYING", "SENIOR_ENTRY",
		"SENIOR_REDEEM_TOKENS", "SENIOR_EXIT", "SENIOR_REDEEM_UNDERLYING", "JTOKEN_SEND", "JTOKEN_RECEIVE", "STOKEN_SEND", "STOKEN_RECEIVE"}
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

func getTotalPoints(window string) string {
	if window == "24h" {
		return "30 * 60"
	}
	if window == "1w" {
		return "3*60*60"
	}
	if window == "30d" {
		return "12*60*60"
	}
	return ""
}

func getTxTokenSymbol(txType, poolTokenSymbol, juniorTokenSymbol, seniorTokenSymbol string) string {
	tokenActions := map[string]string{
		"JUNIOR_ENTRY":             poolTokenSymbol,
		"SENIOR_ENTRY":             poolTokenSymbol,
		"JUNIOR_REDEEM_UNDERLYING": poolTokenSymbol,
		"SENIOR_REDEEM_UNDERLYING": poolTokenSymbol,
		"JUNIOR_EXIT":              juniorTokenSymbol,
		"JUNIOR_REDEEM_TOKENS":     juniorTokenSymbol,
		"JTOKEN_SEND":              juniorTokenSymbol,
		"JTOKEN_RECEIVE":           juniorTokenSymbol,
		"SENIOR_EXIT":              seniorTokenSymbol,
		"SENIOR_REDEEM_TOKENS":     seniorTokenSymbol,
		"STOKEN_SEND":              seniorTokenSymbol,
		"STOKEN_RECEIVE":           seniorTokenSymbol,
	}
	return tokenActions[txType]
}

func getAmountInAsset(txType string, amount, juniorTokenPrice, seniorTokenPrice decimal.Decimal) decimal.Decimal {
	switch txType {
	case "JUNIOR_ENTRY", "SENIOR_ENTRY", "JUNIOR_REDEEM_UNDERLYING", "SENIOR_REDEEM_UNDERLYING":
		return amount
	case "JUNIOR_EXIT", "JUNIOR_REDEEM_TOKENS", "JTOKEN_SEND", "JTOKEN_RECEIVE":
		return amount.Mul(juniorTokenPrice)
	case "SENIOR_EXIT", "SENIOR_REDEEM_TOKENS", "STOKEN_SEND", "STOKEN_RECEIVE":
		return amount.Mul(seniorTokenPrice)
	default:
		return decimal.Zero
	}
}
