package models

const TableSysconfig = "sysconfig"

type Sysconfig struct {
	Id         int64  `db:"id,omitempty"`
	CfgName    string `db:"cfg_name,omitempt"`
	CfgValue   string `db:"cfg_val,omitempt"`
	CfgType    string `db:"cfg_type,omitempt"`
	CfgComment string `db:"cfg_comment,omitempt"`
	CreateAt   string `db:"create_at,omitempty"`
	UpdateAt   string `db:"update_at,omitempty"`
}

const TableBlock = "block"

type Block struct {
	Id           int64  `db:"id,omitempty"`
	Height       int64  `db:"height,omitempty"`
	Hash         string `db:"hash,omitempty"`
	BlockTime    int64  `db:"block_time,omitempty"`
	FirstVersion string `db:"first_version,omitempty"`
	LastVersion  string `db:"last_version,omitempty"`
	CreateAt     string `db:"create_at,omitempty"`
	UpdateAt     string `db:"update_at,omitempty"`
}

const TableTransaction = "transaction"

type Transaction struct {
	Id             int64  `db:"id,omitempty"`
	Version        int64  `db:"version,omitempty"`
	Hash           string `db:"hash,omitempty"`
	TxTime         int64  `db:"tx_time,omitempty"`
	Success        bool   `db:"success,omitempty"`
	SequenceNumber int64  `db:"sequence_number,omitempty"`
	GasUsed        string `db:"gas_used,omitempty"`
	GasPrice       string `db:"gas_price,omitempty"`
	Gas            string `db:"gas,omitempty"`
	Type           string `db:"type,omitempty"`
	Sender         string `db:"sender,omitempty"`
	Receiver       string `db:"receiver,omitempty"`
	TxValue        string `db:"tx_value,omitempty"`
	CreateAt       string `db:"create_at,omitempty"`
	UpdateAt       string `db:"update_at,omitempty"`
}

const TablePayload = "payload"

type Payload struct {
	Id             int64  `db:"id,omitempty"`
	Version        int64  `db:"version,omitempty"`
	Hash           string `db:"hash,omitempty"`
	TxTime         int64  `db:"tx_time,omitempty"`
	SequenceNumber int64  `db:"sequence_number,omitempty"`
	Sender         string `db:"sender,omitempty"`
	PayloadFunc    string `db:"payload_func,omitempty"`
	PayloadType    string `db:"payload_type,omitempty"`
	CreateAt       string `db:"create_at,omitempty"`
	UpdateAt       string `db:"update_at,omitempty"`
}

const TableRecordCoin = "record_coin"

type RecordCoin struct {
	Id           int64  `db:"id,omitempty"`
	Version      int64  `db:"version,omitempty"`
	Hash         string `db:"hash,omitempty"`
	TxTime       int64  `db:"tx_time,omitempty"`
	Sender       string `db:"sender,omitempty"`
	ModuleName   string `db:"module_name,omitempty"`
	ContractName string `db:"contract_name,omitempty"`
	Resource     string `db:"resource,omitempty"`
	Name         string `db:"name,omitempty"`
	Symbol       string `db:"symbol,omitempty"`
	CreateAt     string `db:"create_at,omitempty"`
	UpdateAt     string `db:"update_at,omitempty"`
}

const TableHistoryCoin = "history_coin"

type HistoryCoin struct {
	Id       int64  `db:"id,omitempty"`
	Version  int64  `db:"version,omitempty"`
	Hash     string `db:"hash,omitempty"`
	TxTime   int64  `db:"tx_time,omitempty"`
	Sender   string `db:"sender,omitempty"`
	Receiver string `db:"receiver,omitempty"`
	Resource string `db:"resource,omitempty"`
	Amount   string `db:"amount,omitempty"`
	Action   int64  `db:"action,omitempty"`
	CreateAt string `db:"create_at,omitempty"`
	UpdateAt string `db:"update_at,omitempty"`
}