package types

type TxType string

const (
	TxJuniorDeposit         TxType = "JUNIOR_DEPOSIT"
	TxJuniorInstantWithdraw TxType = "JUNIOR_INSTANT_WITHDRAW"
	TxJuniorRegularWithdraw TxType = "JUNIOR_REGULAR_WITHDRAW"
	TxJuniorRedeem          TxType = "JUNIOR_REDEEM"
	TxSeniorDeposit         TxType = "SENIOR_DEPOSIT"
	TxSeniorRedeem          TxType = "SENIOR_REDEEM"
	TxJtokenSend            TxType = "JTOKEN_SEND"
	TxJtokenReceive         TxType = "JTOKEN_RECEIVE"
	TxJbondSend             TxType = "JBOND_SEND"
	TxJbondReceive          TxType = "JBOND_RECEIVE"
	TxSbondSend             TxType = "SBOND_SEND"
	TxSbondReceive          TxType = "SBOND_RECEIVE"
	TxJuniorStake           TxType = "JUNIOR_STAKE"
	TxJuniorUnstake         TxType = "JUNIOR_UNSTAKE"
)
