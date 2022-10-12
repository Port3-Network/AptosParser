package models

type BlockRsp struct {
	BlockHeight    string           `json:"block_height"`
	BlockHash      string           `json:"block_hash"`
	BlockTimestamp string           `json:"block_timestamp"`
	FirstVersion   string           `json:"first_version"`
	LastVersion    string           `json:"last_version"`
	Transactions   []TransactionRsp `json:"transactions"`
}

type TransactionRsp struct {
	Type                    string     `json:"type"`
	Version                 string     `json:"version"`
	Hash                    string     `json:"hash"`
	StateChangeHash         string     `json:"state_change_hash"`
	EventRootHash           string     `json:"event_root_hash"`
	StateCheckpointHash     string     `json:"state_checkpoint_hash"`
	GasUsed                 string     `json:"gas_used"`
	Success                 bool       `json:"success"`
	VmStatus                string     `json:"vm_status"`
	AccumulatorRootHash     string     `json:"accumulator_root_hash"`
	Changes                 []TxChange `json:"changes"` //
	Sender                  string     `json:"sender"`
	SequenceNumber          string     `json:"sequence_number"`
	MaxGasAmount            string     `json:"max_gas_amount"`
	GasUnitPrice            string     `json:"gas_unit_price"`
	ExpirationTimestampSecs string     `json:"expiration_timestamp_secs"`
	Payload                 TxPayload  `json:"payload"` //
	Events                  []TxEvent  `json:"events"`  //
	Timestamp               string     `json:"timestamp"`
	// Signature               BlockTxSignature `json:"signature"` //
}

type TxChange struct {
	Type         string     `json:"type"`
	Address      string     `json:"address"`
	StateKeyHash string     `json:"state_key_hash"`
	Module       string     `json:"module"`   // ex: 0x1::aptos_coin
	Resource     string     `json:"resource"` // ex: 0x1::coin::CoinStore<0x1::aptos_coin::AptosCoin>
	Data         ChangeData `json:"data"`
}

type ChangeData struct {
	Type string         `json:"type"`
	Data ChangeDataSubD `json:"data"`
}

type ChangeDataSubD struct {
	Decimals              int64           `json:"decimals"`
	Name                  string          `json:"name"`
	Symbol                string          `json:"symbol"`
	WithdrawEvent         ChangeDataEvent `json:"withdraw_events"`
	DepositEvent          ChangeDataEvent `json:"deposit_events"`
	CreateCollectionEvent ChangeDataEvent `json:"create_collection_events"`
	CreateTokenDataEvent  ChangeDataEvent `json:"create_token_data_events"`
}

type ChangeDataEvent struct {
	Counter string `json:"counter"`
	Guid    Guid   `json:"guid"`
}
type Guid struct {
	ID GuidID `json:"id"`
}
type GuidID struct {
	Addr        string `json:"addr"`
	CreationNum string `json:"creation_num"`
}

type TxPayload struct {
	Type          string        `json:"type"`
	Function      string        `json:"function"`
	TypeArguments []string      `json:"type_arguments"`
	Arguments     []interface{} `json:"arguments"`
}

type TxEvent struct {
	Key            string    `json:"key"`
	SequenceNumber string    `json:"sequence_number"` //
	Type           string    `json:"type"`            //
	Guid           EventGuid `json:"guid"`            //
	Data           EventData `json:"data"`            //
}

type EventData struct {
	Amount         string      `json:"amount"`
	CollectionName string      `json:"collection_name"`
	Name           string      `json:"name"`
	Creator        string      `json:"creator"`
	Description    string      `json:"description"`
	Maximum        string      `json:"maximum"`
	Uri            string      `json:"uri"`
	Id             interface{} `json:"id,omitempty"`
}

type EventDataId struct {
	EventTokenDataId
	TokenDataId EventTokenDataId `json:"token_data_id"`
}

type EventTokenDataId struct {
	Collection string `json:"collection"`
	Creator    string `json:"creator"`
	Name       string `json:"name"`
}

type EventGuid struct {
	CreationNumber string `json:"creation_number"`
	AccountAddress string `json:"account_address"`
}

type ResourceRsp struct {
	Type string       `json:"type"`
	Data ResourceData `json:"data"`
}

type ResourceData struct {
	Coin CoinValue `json:"coin"`
}

type CoinValue struct {
	Value string `json:"value"`
}
