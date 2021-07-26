package yieldfarming

func checkTxType(action string) bool {
	txType := [2]string{"DEPOSIT", "WITHDRAW"}
	for _, tx := range txType {
		if action == tx {
			return true
		}
	}

	return false
}
