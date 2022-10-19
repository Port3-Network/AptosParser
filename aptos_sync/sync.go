package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Port3-Network/AptosParser/models"
	oo "github.com/Port3-Network/liboo"
)

func FullSync() {
	var (
		limit    int64 = GDatabase.BlockCount
		minCount int64 = 16
	)
	for {
		syncNum, err := GetSyncBlockNum()
		if err != nil {
			oo.LogW("SyncAllNFTInfo GetSyncBcnum failed: %v", err)
			continue
		}
		oo.LogD("SyncAllNFTInfo GetSyncBcnum got number: %d", syncNum)

		start := syncNum
		var end int64
		txs, err := GetTransactions(strconv.FormatInt(start, 10), limit)
		if err != nil {
			switch err.Error() {
			case "getBuf err":
				if limit >= minCount {
					limit = limit / 2
				} else {
					time.Sleep(time.Second * 3)
				}
			case "statusCode err":
				updateRpc()
			default:
				oo.LogD("GetTransactions err, msg: %v", err)
			}
			continue
		}
		if int(limit) != len(*txs) {
			end = start + int64(len(*txs))
		} else {
			end = start + limit
		}
		saver := NewDbSaver(uint64(end), 0)
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

func GetTransactions(start string, limit int64) (r *[]models.TransactionRsp, err error) {
	r = &[]models.TransactionRsp{}
	sTime := time.Now().UnixMilli()
	url := fmt.Sprintf("%s/transactions?start=%s&limit=%d", GRpc, start, limit)
	buf, _, err := models.HttpGet(url, 2)
	if err != nil {
		return r, err
	}
	if buf == nil {
		return nil, oo.NewError("getBuf err")
	}

	err = json.Unmarshal(buf, &r)
	if err != nil {
		return r, fmt.Errorf("tx jsonUnmarshal msg: %v", err)
	}
	eTime := time.Now().UnixMilli()
	oo.LogD("http due: %vms\n", eTime-sTime)
	return r, nil
}

func GetBlocks(blockNum int64) (*models.BlockRsp, string, error) {
	r := &models.BlockRsp{}
	url := fmt.Sprintf("%s/blocks/by_height/%d?with_transactions=false", GDatabase.TxRpcUrl, blockNum)
	buf, blockHeight, err := models.HttpGet(url, 2)
	if err != nil {
		return r, "", fmt.Errorf("blocks HttpGet msg: %v", err)
	}
	err = json.Unmarshal(buf, &r)
	if err != nil {
		return r, "", fmt.Errorf("blocks Unmarshal msg: %v", err)
	}

	return r, blockHeight, nil
}
