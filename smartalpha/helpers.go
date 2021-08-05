package smartalpha

import (
	"context"
	"errors"
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