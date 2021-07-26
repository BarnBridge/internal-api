package yieldfarming

import (
	"strconv"
	"time"

	"github.com/barnbridge/internal-api/yieldfarming/types"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

func checkTxType(action string) bool {
	txType := [2]string{"DEPOSIT", "WITHDRAW"}
	for _, tx := range txType {
		if action == tx {
			return true
		}
	}

	return false
}

func validateTs(ts string) (*time.Time, error) {
	timestamp, err := strconv.ParseInt(ts, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "invalid timestamp")
	}

	if timestamp == -1 {
		return nil, errors.New("timestamp required")
	}

	tm := time.Unix(timestamp, 0)

	return &tm, nil
}

func getChart(rows pgx.Rows) (*types.Chart, error) {
	x := make(types.Chart)

	for rows.Next() {
		var t time.Time
		var a types.Aggregate
		err := rows.Scan(&t, &a.SumDeposits, &a.SumWithdrawals)
		if err != nil {
			return nil, err
		}

		x[t] = a
	}
	return &x, nil
}
