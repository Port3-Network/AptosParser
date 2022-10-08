package main

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Port3-Network/AptosParser/models"
	oo "github.com/Port3-Network/liboo"
)

type DbSaver struct {
	height uint64

	block        []*models.Block
	transaction  []*models.Transaction
	payload      []*models.Payload
	recordToken  []*models.RecordToken
	historyToken []*models.HistoryToken
	assetToken   []*models.AssetToken
	// recordToken  map[string]*models.RecordToken
}

func NewDbSaver(height, block_ts uint64) *DbSaver {
	return &DbSaver{
		height: height,

		block:       make([]*models.Block, 0),
		transaction: make([]*models.Transaction, 0),
		payload:     make([]*models.Payload, 0),
		recordToken: make([]*models.RecordToken, 0),
		// recordToken:  make(map[string]*models.RecordToken),
		historyToken: make([]*models.HistoryToken, 0),
		assetToken:   make([]*models.AssetToken, 0),
	}
}

func (db *DbSaver) Commit() error {
	oo.LogD("db(height: %d) Commit", db.height)

	dbConn, dbTx, err := oo.NewSqlTxConn()
	if err != nil {
		oo.LogD("models.NewSqlTxConn %v", err)
		return err
	}
	defer oo.CloseSqlTxConn(dbConn, dbTx, &err)

	if err := db.doCommitSyncHeight(dbTx); err != nil {
		oo.LogD("doCommitSyncHeight %v", err)
		return err
	}
	if err := db.doCommitBlock(dbTx); err != nil {
		oo.LogD("doCommitBlock %v", err)
		return err
	}
	if err := db.doCommitTransaction(dbTx); err != nil {
		oo.LogD("doCommitTransaction %v", err)
		return err
	}
	if err := db.doCommitPayload(dbTx); err != nil {
		oo.LogD("doCommitPayload %v", err)
		return err
	}
	if err := db.doCommitRecordToken(dbTx); err != nil {
		oo.LogD("doCommitRecordToken %v", err)
		return err
	}
	if err := db.doCommitHistoryToken(dbTx); err != nil {
		oo.LogD("doCommitHistoryToken %v", err)
		return err
	}
	return nil
}

func (db *DbSaver) doCommitSyncHeight(tx *sql.Tx) (err error) {
	sqlStr := fmt.Sprintf(`UPDATE sysconfig SET cfg_val="%d" WHERE cfg_name="%s"`, db.height, SyncLatestBlock)
	return oo.SqlTxExec(tx, sqlStr)
}

func (db *DbSaver) doCommitBlock(tx *sql.Tx) (err error) {
	if len(db.block) == 0 {
		return nil
	}
	var vals []string
	for _, block := range db.block {
		v := fmt.Sprintf("(%d,'%s',%d,'%s','%s')",
			block.Height, block.Hash, block.BlockTime, block.FirstVersion, block.LastVersion)
		vals = append(vals, v)
	}
	sqlStr := fmt.Sprintf(`INSERT INTO %s(height,hash,block_time,first_version,last_version) VALUES %s`, models.TableBlock, strings.Join(vals, ","))
	if err := oo.SqlTxExec(tx, sqlStr); err != nil {
		return err
	}
	return nil
}

func (db *DbSaver) doCommitTransaction(tx *sql.Tx) (err error) {
	if len(db.transaction) == 0 {
		return nil
	}
	var vals []string
	for _, t := range db.transaction {
		v := fmt.Sprintf("(%d,'%s',%d,%t,%d,'%s','%s','%s','%s','%s','%s','%s')",
			t.Version, t.Hash, t.TxTime, t.Success, t.SequenceNumber, t.GasUsed, t.GasPrice, t.Gas, t.Type, t.Sender, t.Receiver, t.TxValue)
		vals = append(vals, v)
	}
	sqlStr := fmt.Sprintf(`INSERT INTO %s(version,hash,tx_time,success,sequence_number,gas_used,gas_price,gas,type,sender,receiver,tx_value) VALUES %s`, models.TableTransaction, strings.Join(vals, ","))
	if err := oo.SqlTxExec(tx, sqlStr); err != nil {
		return err
	}
	return nil
}

func (db *DbSaver) doCommitPayload(tx *sql.Tx) (err error) {
	if len(db.payload) == 0 {
		return nil
	}
	var vals []string
	for _, payload := range db.payload {
		v := fmt.Sprintf("(%d,'%s',%d,%d,'%s','%s','%s')",
			payload.Version, payload.Hash, payload.TxTime, payload.SequenceNumber, payload.Sender, payload.PayloadFunc, payload.PayloadType)
		vals = append(vals, v)
	}
	sqlStr := fmt.Sprintf(`INSERT INTO %s(version,hash,tx_time,sequence_number,sender,payload_func,payload_type) VALUES %s`, models.TablePayload, strings.Join(vals, ","))
	if err := oo.SqlTxExec(tx, sqlStr); err != nil {
		return err
	}
	return nil
}

func (db *DbSaver) doCommitRecordToken(tx *sql.Tx) (err error) {
	if len(db.recordToken) == 0 {
		return nil
	}
	var vals []string
	for _, record := range db.recordToken {
		v := fmt.Sprintf("(%d,'%s',%d,'%s','%s','%s','%s','%s','%s')",
			record.Version, record.Hash, record.TxTime, record.Sender, record.ModuleName, record.ContractName, record.Resource, record.Name, record.Symbol)
		vals = append(vals, v)
	}
	sqlStr := fmt.Sprintf(`INSERT INTO %s(version,hash,tx_time,sender,module_name,contract_name,resource,name,symbol) VALUES %s`, models.TableRecordToken, strings.Join(vals, ","))
	if err := oo.SqlTxExec(tx, sqlStr); err != nil {
		return err
	}
	return nil
}

func (db *DbSaver) doCommitHistoryToken(tx *sql.Tx) (err error) {
	if len(db.historyToken) == 0 {
		return nil
	}

	var vals []string
	for _, history := range db.historyToken {
		v := fmt.Sprintf("(%d,'%s',%d,'%s','%s','%s','%s',%d)",
			history.Version, history.Hash, history.TxTime, history.Sender, history.Receiver, history.Resource, history.Amount, history.Action)
		vals = append(vals, v)
	}
	sqlStr := fmt.Sprintf(`INSERT INTO %s(version,hash,tx_time,sender,receiver,resource,amount,action) VALUES %s`, models.TableHistoryToken, strings.Join(vals, ","))
	if err := oo.SqlTxExec(tx, sqlStr); err != nil {
		return err
	}
	return nil
}

func (db *DbSaver) AddBlock(block *models.Block) {
	db.block = append(db.block, block)
}

func (db *DbSaver) AddTransaction(transaction *models.Transaction) {
	db.transaction = append(db.transaction, transaction)
}

func (db *DbSaver) AddPayload(payload *models.Payload) {
	db.payload = append(db.payload, payload)
}

func (db *DbSaver) HandlerAddRecordToken(resource string, data *models.RecordToken) {
	// if _, ok := db.recordToken[resource]; ok {
	// 	return
	// }
	// db.recordToken[resource] = &models.RecordToken{
	// 	Version:      data.Version,
	// 	Hash:         data.Hash,
	// 	TxTime:       data.TxTime,
	// 	Sender:       data.Sender,
	// 	ModuleName:   data.ModuleName,
	// 	ContractName: data.ContractName,
	// 	Resource:     data.Resource,
	// 	Name:         data.Name,
	// 	Symbol:       data.Symbol,
	// }
	db.recordToken = append(db.recordToken, &models.RecordToken{
		Version:      data.Version,
		Hash:         data.Hash,
		TxTime:       data.TxTime,
		Sender:       data.Sender,
		ModuleName:   data.ModuleName,
		ContractName: data.ContractName,
		Resource:     data.Resource,
		Name:         data.Name,
		Symbol:       data.Symbol,
	})
}
