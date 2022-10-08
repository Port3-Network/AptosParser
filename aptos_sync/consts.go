package main

const (
	SyncLatestBlock = "cache.sync_latest_block"

	NativeAptosCoin = "0x1::aptos_coin::AptosCoin"
	IsCoin          = "0x1::coin::CoinStore"
	NewCoinInfo     = "0x1::coin::CoinInfo"
)

const (
	UserTransaction = "user_transaction"
	EventWithdraw   = "0x1::coin::WithdrawEvent"
	EventDeposit    = "0x1::coin::DepositEvent"
	EventTransfer   = "0x1::aptos_account::transfer"

	ZeroAddress = "0x0000000000000000000000000000000000000000000000000000000000000000"

	ActionMint     = 1
	ActionBurn     = 2
	ActionTransfer = 3
)
