package smartyield

func checkRewardPoolTxType(action string) bool {
	txType := [2]string{"JUNIOR_STAKE", "JUNIOR_UNSTAKE"}
	for _, tx := range txType {
		if action == tx {
			return true
		}
	}

	return false
}
