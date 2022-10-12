package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Port3-Network/AptosParser/models"
	oo "github.com/Port3-Network/liboo"
)

func FullSync() {
	for {
		syncNum, err := GetSyncBlockNum()
		if err != nil {
			oo.LogW("SyncAllNFTInfo GetSyncBcnum failed: %v", err)
			continue
		}
		oo.LogD("SyncAllNFTInfo GetSyncBcnum got number: %d", syncNum)

		start := syncNum
		saver := NewDbSaver(uint64(start)+uint64(GDatabase.BlockCount), 0)

		txs, err := GetTransactions(strconv.FormatInt(start, 10), int(GDatabase.BlockCount))
		if err != nil {
			oo.LogD("GetTransactions err, msg: %v", err)
			continue
		}

		for _, tx := range *txs {
			switch tx.Type {
			case UserTransaction:
				err := handlerUserTransaction(saver, tx)
				if err != nil {
					oo.LogD("handlerUserTransaction err, msg: %v", err)
					continue
				}
			}
		}
		// return
		if err := saver.Commit(); err != nil {
			oo.LogW("saver.Commit err %v", err)
			continue
		}
	}
}

func GetTransactions(start string, limit int) (*[]models.TransactionRsp, error) {
	r := &[]models.TransactionRsp{}
	url := fmt.Sprintf("%s/transactions?start=%s&limit=%d", GDatabase.TxRpcUrl, start, limit)
	buf, err := models.HttpGet(url, 2)
	if err != nil {
		return r, fmt.Errorf("tx HttpGet msg: %v", err)
	}

	err = json.Unmarshal(buf, &r)
	if err != nil {
		return r, fmt.Errorf("tx jsonUnmarshal msg: %v", err)
	}
	return r, nil
}

func GetBlocks(blockNum int64) (*models.BlockRsp, error) {
	r := &models.BlockRsp{}
	url := fmt.Sprintf("%s/blocks/by_height/%d?with_transactions=true", GDatabase.TxRpcUrl, blockNum)
	buf, err := models.HttpGet(url, 2)
	if err != nil {
		return r, fmt.Errorf("blocks HttpGet msg: %v", err)
	}

	err = json.Unmarshal(buf, &r)
	if err != nil {
		return r, fmt.Errorf("blocks Unmarshal msg: %v", err)
	}
	return r, nil
}
