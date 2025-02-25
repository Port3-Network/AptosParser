package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/Port3-Network/AptosParser/models"
	oo "github.com/Port3-Network/liboo"
	"github.com/gogf/gf/v2/util/gconv"
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

	for _, tx := range *r {
		for _, event := range tx.Events {
			if event.Type != EventString {
				//data, _ := json.Marshal(event.RawData)
				//if err = json.Unmarshal(data, &event.Data); err != nil {
				//	return r, fmt.Errorf("event data unmarshal: %v", err)
				//}
				valueType := reflect.TypeOf(event.RawData)
				if valueType.Kind() == reflect.String {
					oo.LogD("event.RawData type of string. continue, data: %v", event.RawData)
					continue
				}

				rawData, ok := event.RawData.(map[string]interface{})
				if !ok {
					oo.LogD("event.RawData.(map[string]interface{}) !ok, data: %v", event.RawData)
					continue
				}
				event.Data.Amount = gconv.String(rawData["amount"])
				event.Data.CollectionName = gconv.String(rawData["collection_name"])
				event.Data.Name = gconv.String(rawData["name"])
				event.Data.Creator = gconv.String(rawData["creator"])
				event.Data.Description = gconv.String(rawData["description"])
				event.Data.Maximum = gconv.String(rawData["maximum"])
				event.Data.Uri = gconv.String(rawData["uri"])
				event.Data.Id = rawData["uri"]
				if typeInfo, found := rawData["type_info"]; found {
					ti, _ := json.Marshal(typeInfo)
					if err = json.Unmarshal(ti, &event.Data.TypeInfo); err != nil {
						return r, fmt.Errorf("event data unmarshal: %v", err)
					}
				}
			}
		}
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
