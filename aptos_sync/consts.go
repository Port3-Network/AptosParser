package main

const (
	SyncLatestBlock = "cache.sync_latest_block"

	NativeAptosCoin = "0x1::aptos_coin::AptosCoin"
)

const (
	UserTransaction = "user_transaction"
	ZeroAddress     = "0x0000000000000000000000000000000000000000000000000000000000000000"

	ActionMint     = 1
	ActionBurn     = 2
	ActionTransfer = 3

	TypeCallFunction      = "entry_function_payload"
	EventWithdraw         = "0x1::coin::WithdrawEvent"
	EventDeposit          = "0x1::coin::DepositEvent"
	EventCollectionCreate = "0x3::token::CreateCollectionEvent"

	FunctionPublishPkg      = "0x1::code::publish_package_txn"
	FunctionMint            = "0x1::managed_coin::mint"
	FunctionTransfer        = "0x1::coin::transfer"
	FunctionAccountTransfer = "0x1::aptos_account::transfer"

	ChangeTypeCoinInfo = "0x1::coin::CoinInfo"
)
