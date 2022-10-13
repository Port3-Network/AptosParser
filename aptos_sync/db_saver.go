package main

import (
	"database/sql"
	"fmt"
	u "net/url"
	"strconv"
	"strings"

	"github.com/Port3-Network/AptosParser/models"
	oo "github.com/Port3-Network/liboo"
)

type nftToken struct {
	Owner      string
	Creator    string
	Collection string
	Name       string
}

type coinInfo struct {
	Owner      string
	ModuleName string
	StructName string
}

type DbSaver struct {
	height uint64

	block        []*models.Block
	transaction  []*models.Transaction
	payload      []*models.Payload
	recordCoin   map[coinInfo]*models.RecordCoin
	historyCoin  []*models.HistoryCoin
	collection   []*models.Collection
	recordToken  []*models.RecordToken
	assetToken   map[nftToken]*models.AssetToken
	historyToken []*models.HistoryToken
}

func NewDbSaver(height, block_ts uint64) *DbSaver {
	return &DbSaver{
		height: height,

		block:        make([]*models.Block, 0),
		transaction:  make([]*models.Transaction, 0),
		payload:      make([]*models.Payload, 0),
		recordCoin:   make(map[coinInfo]*models.RecordCoin),
		historyCoin:  make([]*models.HistoryCoin, 0),
		collection:   make([]*models.Collection, 0),
		recordToken:  make([]*models.RecordToken, 0),
		assetToken:   make(map[nftToken]*models.AssetToken),
		historyToken: make([]*models.HistoryToken, 0),
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
	if err := db.doCommitRecordCoin(dbTx); err != nil {
		oo.LogD("doCommitRecordToken %v", err)
		return err
	}
	if err := db.doCommitHistoryCoin(dbTx); err != nil {
		oo.LogD("doCommitHistoryToken %v", err)
		return err
	}
	if err := db.doCommitCollection(dbTx); err != nil {
		oo.LogD("doCommitCollection %v", err)
		return err
	}
	if err := db.doCommitRecordToken(dbTx); err != nil {
		oo.LogD("doCommitRecordToken %v", err)
		return err
	}
	if err := db.doCommitAssetToken(dbTx); err != nil {
		oo.LogD("doCommitAssetToken %v", err)
		return err
	}
	if err := db.doCommitHistoryToken(dbTx); err != nil {
		oo.LogD("doCommitHistoryToken %v", err)
		return err
	}
	if err := db.doCommitSyncHeight(dbTx); err != nil {
		oo.LogD("doCommitSyncHeight %v", err)
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
		v := fmt.Sprintf("('%s','%s',%d,%t,%d,'%s','%s','%s','%s','%s','%s','%s')",
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
		v := fmt.Sprintf("('%s','%s',%d,%d,'%s','%s','%s')",
			payload.Version, payload.Hash, payload.TxTime, payload.SequenceNumber, payload.Sender, payload.PayloadFunc, payload.PayloadType)
		vals = append(vals, v)
	}
	sqlStr := fmt.Sprintf(`INSERT INTO %s(version,hash,tx_time,sequence_number,sender,payload_func,payload_type) VALUES %s`, models.TablePayload, strings.Join(vals, ","))
	if err := oo.SqlTxExec(tx, sqlStr); err != nil {
		return err
	}
	return nil
}

func (db *DbSaver) doCommitRecordCoin(tx *sql.Tx) (err error) {
	if len(db.recordCoin) == 0 {
		return nil
	}
	for _, record := range db.recordCoin {
		var recordCoin models.RecordCoin
		sqlStr := oo.NewSqler().Table(models.TableRecordCoin).
			Where("resource", record.Resource).Select("*")
		err := oo.SqlGet(sqlStr, &recordCoin)
		if err != nil && err != oo.ErrNoRows {
			return err
		}
		if err == oo.ErrNoRows {
			sqlStr := oo.NewSqler().Table(models.TableRecordCoin).
				Insert(map[string]interface{}{
					"version":       record.Version,
					"hash":          record.Hash,
					"tx_time":       record.TxTime,
					"sender":        record.Sender,
					"module_name":   record.ModuleName,
					"contract_name": record.ContractName,
					"resource":      record.Resource,
					"name":          record.Name,
					"symbol":        record.Symbol,
				})
			if err := oo.SqlTxExec(tx, sqlStr); err != nil {
				return err
			}
		} else {
			sqlStr := oo.NewSqler().Table(models.TableRecordCoin).
				Where("resource", record.Resource).
				Update(map[string]interface{}{
					"version": record.Version,
					"hash":    record.Hash,
					"tx_time": record.TxTime,
					"name":    record.Name,
					"symbol":  record.Symbol,
				})
			if err := oo.SqlTxExec(tx, sqlStr); err != nil {
				return err
			}
		}
	}
	return nil
}

func (db *DbSaver) doCommitHistoryCoin(tx *sql.Tx) (err error) {
	if len(db.historyCoin) == 0 {
		return nil
	}

	var vals []string
	for _, history := range db.historyCoin {
		v := fmt.Sprintf("('%s','%s',%d,'%s','%s','%s','%s',%d)",
			history.Version, history.Hash, history.TxTime, history.Sender, history.Receiver, history.Resource, history.Amount, history.Action)
		vals = append(vals, v)
	}
	sqlStr := fmt.Sprintf(`INSERT INTO %s(version,hash,tx_time,sender,receiver,resource,amount,action) VALUES %s`, models.TableHistoryCoin, strings.Join(vals, ","))
	if err := oo.SqlTxExec(tx, sqlStr); err != nil {
		return err
	}
	return nil
}

func (db *DbSaver) doCommitCollection(tx *sql.Tx) (err error) {
	if len(db.collection) == 0 {
		return nil
	}

	var vals []string
	for _, collection := range db.collection {
		v := fmt.Sprintf("('%s','%s',%d,'%s','%s','%s','%s','%s','%s','%s')",
			collection.Version, collection.Hash, collection.TxTime, collection.Sender, collection.Creator, u.QueryEscape(collection.Name), u.QueryEscape(collection.Description), u.QueryEscape(collection.Uri), collection.Maximum, collection.Type)
		vals = append(vals, v)
	}
	sqlStr := fmt.Sprintf(`INSERT INTO %s(version,hash,tx_time,sender,creator,name,description,uri,maximum,type) VALUES %s`, models.TableCollection, strings.Join(vals, ","))
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
		v := fmt.Sprintf("('%s','%s',%d,'%s','%s','%s','%s','%s','%s','%s','%s')",
			record.Version, record.Hash, record.TxTime, record.Sender, record.Creator, u.QueryEscape(record.Collection), u.QueryEscape(record.Name), u.QueryEscape(record.Description), u.QueryEscape(record.Uri), record.Maximum, record.Type)
		vals = append(vals, v)
	}
	sqlStr := fmt.Sprintf(`INSERT INTO %s(version,hash,tx_time,sender,creator,collection,name,description,uri,maximum,type) VALUES %s`, models.TableRecordToken, strings.Join(vals, ","))
	if err := oo.SqlTxExec(tx, sqlStr); err != nil {
		return err
	}
	return nil
}

func (db *DbSaver) doCommitAssetToken(tx *sql.Tx) (err error) {
	if len(db.assetToken) == 0 {
		return nil
	}

	for _, asset := range db.assetToken {
		var assetToken models.AssetToken
		sqlStr := oo.NewSqler().Table(models.TableAssetToken).
			Where("owner", asset.Owner).
			Where("creator", asset.Creator).
			Where("collection", u.QueryEscape(asset.Collection)).
			Where("name", u.QueryEscape(asset.Name)).
			Select("*")
		err := oo.SqlGet(sqlStr, &assetToken)
		if err != nil && err != oo.ErrNoRows {
			return err
		}
		if err == oo.ErrNoRows {
			sqlStr := oo.NewSqler().Table(models.TableAssetToken).
				Insert(map[string]interface{}{
					"version":    asset.Version,
					"hash":       asset.Hash,
					"tx_time":    asset.TxTime,
					"owner":      asset.Owner,
					"creator":    asset.Creator,
					"collection": u.QueryEscape(asset.Collection),
					"name":       u.QueryEscape(asset.Name),
					"amount":     asset.Amount,
				})
			if err := oo.SqlTxExec(tx, sqlStr); err != nil {
				return err
			}
		} else {
			amountToken := models.ParseInt64(assetToken.Amount)
			amount := string(strconv.FormatInt(models.ParseInt64(assetToken.Amount)+amountToken, 10))
			sqlStr := oo.NewSqler().Table(models.TableAssetToken).
				Where("owner", asset.Owner).
				Where("creator", asset.Creator).
				Where("collection", u.QueryEscape(asset.Collection)).
				Where("name", u.QueryEscape(asset.Name)).
				Update(map[string]interface{}{
					"version": asset.Version,
					"hash":    asset.Hash,
					"tx_time": asset.TxTime,
					"amount":  amount,
				})
			if err := oo.SqlTxExec(tx, sqlStr); err != nil {
				return err
			}
		}
	}
	return nil
}

func (db *DbSaver) doCommitHistoryToken(tx *sql.Tx) (err error) {
	if len(db.historyToken) == 0 {
		return nil
	}

	var vals []string
	for _, history := range db.historyToken {
		v := fmt.Sprintf("('%s','%s',%d,'%s','%s','%s','%s','%s','%s',%d)",
			history.Version, history.Hash, history.TxTime, history.Sender, history.Receiver, history.Creator, u.QueryEscape(history.Collection), u.QueryEscape(history.Name), history.Amount, history.Action)
		vals = append(vals, v)
	}
	sqlStr := fmt.Sprintf(`INSERT INTO %s(version,hash,tx_time,sender,receiver,creator,collection,name,amount,action) VALUES %s`, models.TableHistoryToken, strings.Join(vals, ","))
	if err := oo.SqlTxExec(tx, sqlStr); err != nil {
		return err
	}
	return nil
}
