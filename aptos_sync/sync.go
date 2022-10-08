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
		// block, err := GetBlocks(start)
		// if err != nil {
		// 	oo.LogD("GetBlocks err, msg: %v", err)
		// 	continue
		// }

		// height, _ := strconv.ParseInt(block.BlockHeight, 10, 64)
		// blockTime, _ := strconv.ParseInt(block.BlockTimestamp, 10, 64)
		// saver.AddBlock(&models.Block{
		// 	Height:       height,
		// 	Hash:         block.BlockHash,
		// 	BlockTime:    blockTime,
		// 	FirstVersion: block.FirstVersion,
		// 	LastVersion:  block.LastVersion,
		// })

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

		if err := saver.Commit(); err != nil {
			oo.LogW("saver.Commit err %v", err)
			continue
		}
	}

}

func GetTransactions(start string, limit int) (*[]models.TransactionRsp, error) {
	r := &[]models.TransactionRsp{}
	url := fmt.Sprintf("https://fullnode.devnet.aptoslabs.com/v1/transactions?start=%s&limit=%d", start, limit)
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
	url := fmt.Sprintf("https://fullnode.devnet.aptoslabs.com/v1/blocks/by_height/%d?with_transactions=true", blockNum)
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
