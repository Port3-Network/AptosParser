package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Port3-Network/AptosParser/models"
	oo "github.com/Port3-Network/liboo"
)

func parseInterface(v interface{}, nft *models.EventTokenDataId) {
	data, ok := v.(map[string]interface{})
	if !ok {
		return
	}
	for k, v := range data {
		switch val := v.(type) {
		// case int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint:
		// 	if strings.EqualFold(k, "collection") && len(nft.Collection) <= 0 {
		// 		nft.Collection = fmt.Sprintf("%d", val)
		// 	} else if strings.EqualFold(k, "creator") && len(nft.Creator) <= 0 {
		// 		nft.Creator = fmt.Sprintf("%d", val)
		// 	} else if strings.EqualFold(k, "name") && len(nft.Name) <= 0 {
		// 		nft.Name = fmt.Sprintf("%d", val)
		// 	}
		// case float32, float64:
		// 	if strings.EqualFold(k, "collection") && len(nft.Collection) <= 0 {
		// 		nft.Collection = fmt.Sprintf("%.0f", val)
		// 	} else if strings.EqualFold(k, "creator") && len(nft.Creator) <= 0 {
		// 		nft.Creator = fmt.Sprintf("%.0f", val)
		// 	} else if strings.EqualFold(k, "name") && len(nft.Name) <= 0 {
		// 		nft.Name = fmt.Sprintf("%.0f", val)
		// 	}
		case string:
			if strings.EqualFold(k, "collection") && len(nft.Collection) <= 0 {
				nft.Collection = val
			} else if strings.EqualFold(k, "creator") && len(nft.Creator) <= 0 {
				nft.Creator = val
			} else if strings.EqualFold(k, "name") && len(nft.Name) <= 0 {
				nft.Name = val
			}
		case interface{}:
			parseInterface(val, nft)
			// case []interface{}:

		}
	}
}
func handlerUserTransaction(db *DbSaver, data models.TransactionRsp) error {
	switch data.Payload.Type {
	case TypeCallFunction:
		txTime, _ := strconv.ParseInt(data.Timestamp, 10, 64)
		sequenceNum, _ := strconv.ParseInt(data.SequenceNumber, 10, 64)

		// payload -> done
		handlerPayload(db, txTime, sequenceNum, data)

		// payloadDetail -> done
		handlerPayloadDetail(db, txTime, data)

		// transaction -> done
		handlerTx(db, txTime, sequenceNum, data)

		// record ->
		// if data.Payload.Function == FunctionPublishPkg {
		handlerRecordCoin(db, txTime, sequenceNum, data)
		// 	return nil
		// }

		// history -> done
		handlerHistoryCoin(db, txTime, sequenceNum, data)

		// collection -> done
		handlerCollection(db, txTime, sequenceNum, data)

		// recordToken -> done
		handlerRecordToken(db, txTime, sequenceNum, data)

		// assetToken -> done
		handlerAssetToken(db, txTime, sequenceNum, data)

		// historyToken ->
		handlerHistoryToken(db, txTime, sequenceNum, data)
	default:
		oo.LogD("payload type [%s] not found", data.Payload.Type)
	}
	return nil
}

func handlerPayload(saver *DbSaver, txTime, sequenceNum int64, data models.TransactionRsp) {
	payloadFunc := data.Payload.Function
	if len(data.Payload.Function) > 128 {
		payloadFunc = data.Payload.Function[:128]
	}
	saver.payload = append(saver.payload, &models.Payload{
		Version:        data.Version,
		Hash:           data.Hash,
		TxTime:         txTime,
		SequenceNumber: sequenceNum,
		Sender:         data.Sender,
		PayloadFunc:    payloadFunc,
		PayloadType:    data.Payload.Type,
	})
}

func handlerPayloadDetail(saver *DbSaver, txTime int64, data models.TransactionRsp) {
	payloadFunc := data.Payload.Function
	if len(data.Payload.Function) > 128 {
		payloadFunc = data.Payload.Function[:128]
	}
	typeArguments, _ := json.Marshal(data.Payload.TypeArguments)
	arguments, _ := json.Marshal(data.Payload.Arguments)

	saver.payloadDetail = append(saver.payloadDetail, &models.PayloadDetail{
		Version:       data.Version,
		Hash:          data.Hash,
		TxTime:        txTime,
		Success:       data.Success,
		Sender:        data.Sender,
		PayloadFunc:   payloadFunc,
		TypeArguments: typeArguments,
		Arguments:     arguments,
	})
}

func handlerTx(saver *DbSaver, txTime, sequenceNum int64, data models.TransactionRsp) {
	saver.transaction = append(saver.transaction, &models.Transaction{
		Version:        data.Version,
		Hash:           data.Hash,
		TxTime:         txTime,
		Success:        data.Success,
		SequenceNumber: sequenceNum,
		GasUsed:        data.GasUsed,
		GasPrice:       data.GasUnitPrice,
		Gas:            data.MaxGasAmount,
		Type:           data.Type,
		Sender:         data.Sender,
		TxValue:        "0",
	})
}

func handlerHistoryCoin(saver *DbSaver, txTime, sequenceNum int64, data models.TransactionRsp) {
	for _, event := range data.Events {
		var resource, amount string
		var action int64
		var sender, receiver string = ZeroAddress, ZeroAddress
		cnum := event.Guid.CreationNumber
		addr := event.Guid.AccountAddress
		amount = event.Data.Amount
		switch event.Type {
		case EventDeposit:
			receiver = event.Guid.AccountAddress
			for _, c := range data.Changes {
				if cnum == c.Data.Data.DepositEvent.Guid.ID.CreationNum && addr == c.Data.Data.DepositEvent.Guid.ID.Addr {
					t := ParseType(c.Data.Type)
					if t == nil {
						continue
					}
					resource = t.Resource
				}
			}
		case EventWithdraw:
			sender = event.Guid.AccountAddress
			for _, c := range data.Changes {
				if cnum == c.Data.Data.WithdrawEvent.Guid.ID.CreationNum && addr == c.Data.Data.WithdrawEvent.Guid.ID.Addr {
					t := ParseType(c.Data.Type)
					if t == nil {
						continue
					}
					resource = t.Resource
				}
			}
		}
		if sender == receiver && sender == ZeroAddress {
			continue
		} else if sender == ZeroAddress {
			action = ActionMint
		} else if receiver == ZeroAddress {
			action = ActionBurn
		} else {
			action = ActionTransfer
		}
		saver.historyCoin = append(saver.historyCoin, &models.HistoryCoin{
			Version:  data.Version,
			Hash:     data.Hash,
			TxTime:   txTime,
			Sender:   sender,
			Receiver: receiver,
			Resource: resource,
			Amount:   amount,
			Action:   action,
		})
	}
}

func handlerRecordCoin(saver *DbSaver, txTime, sequenceNum int64, data models.TransactionRsp) {
	for _, event := range data.Events {
		if event.Type != EventCoinRegister {
			continue
		}
		coinKey := coinInfo{
			Owner:      event.Data.TypeInfo.AccountAddress,
			ModuleName: "",
			StructName: "",
		}
		mdName := event.Data.TypeInfo.ModuleName
		if strings.HasPrefix(mdName, "0x") && len(mdName) > 2 {
			n, _ := hex.DecodeString(mdName[2:])
			coinKey.ModuleName = string(n)
		}

		sName := event.Data.TypeInfo.StructName
		if strings.HasPrefix(sName, "0x") && len(sName) > 2 {
			s, _ := hex.DecodeString(sName[2:])
			coinKey.StructName = string(s)
		}

		var name, symbol string
		var decimals int64
		resource := fmt.Sprintf("%s::%s::%s", coinKey.Owner, coinKey.ModuleName, coinKey.StructName)
		dataType := fmt.Sprintf("0x1::coin::CoinInfo<%s>", resource)
		for _, change := range data.Changes {
			if change.Data.Type == dataType {
				decimals = change.Data.Data.Decimals
				name = change.Data.Data.Name
				symbol = change.Data.Data.Symbol
			}
		}

		record, ok := saver.recordCoin[coinKey]
		if !ok {
			record = &models.RecordCoin{
				Version:      data.Version,
				Hash:         data.Hash,
				TxTime:       txTime,
				Sender:       data.Sender,
				ModuleName:   coinKey.ModuleName,
				ContractName: coinKey.StructName,
				Name:         name,
				Symbol:       symbol,
				Decimals:     decimals,
				Resource:     resource,
			}
			saver.recordCoin[coinKey] = record
		} else {
			record.Version = data.Version
			record.Hash = data.Hash
		}
	}
}

func handlerCollection(saver *DbSaver, txTime, sequenceNum int64, data models.TransactionRsp) {
	for _, event := range data.Events {
		if event.Type != EventCollectionCreate {
			continue
		}
		var cType string
		cnum := event.Guid.CreationNumber
		addr := event.Guid.AccountAddress
		for _, c := range data.Changes {
			if cnum == c.Data.Data.CreateCollectionEvent.Guid.ID.CreationNum && addr == c.Data.Data.CreateCollectionEvent.Guid.ID.Addr {
				cType = c.Data.Type
			}
		}
		saver.collection = append(saver.collection, &models.Collection{
			Version:     data.Version,
			Hash:        data.Hash,
			TxTime:      txTime,
			Sender:      data.Sender,
			Creator:     event.Data.Creator,
			Name:        event.Data.CollectionName,
			Description: event.Data.Description,
			Uri:         event.Data.Uri,
			Maximum:     event.Data.Maximum,
			Type:        cType,
		})

	}
}

func handlerRecordToken(saver *DbSaver, txTime, sequenceNum int64, data models.TransactionRsp) {
	for _, event := range data.Events {
		if event.Type != EventTokenCreate {
			continue
		}

		var cType string
		cnum := event.Guid.CreationNumber
		addr := event.Guid.AccountAddress
		for _, c := range data.Changes {
			if cnum == c.Data.Data.CreateTokenDataEvent.Guid.ID.CreationNum && addr == c.Data.Data.CreateTokenDataEvent.Guid.ID.Addr {
				cType = c.Data.Type
			}
		}
		var nftTokenData models.EventTokenDataId
		parseInterface(event.Data.Id, &nftTokenData)
		saver.recordToken = append(saver.recordToken, &models.RecordToken{
			Version:    data.Version,
			Hash:       data.Hash,
			TxTime:     txTime,
			Sender:     data.Sender,
			Creator:    nftTokenData.Creator,
			Collection: nftTokenData.Collection,
			// Creator:     event.Data.Id.Creator,
			// Collection:  event.Data.Id.Collection,
			Name:        event.Data.Name,
			Description: event.Data.Description,
			Uri:         event.Data.Uri,
			Maximum:     event.Data.Maximum,
			Type:        cType,
		})
	}
}

func handlerAssetToken(saver *DbSaver, txTime, sequenceNum int64, data models.TransactionRsp) {
	for _, event := range data.Events {
		var nftTokenData models.EventTokenDataId
		parseInterface(event.Data.Id, &nftTokenData)

		switch event.Type {
		case EventTokenDeposit:
			ownerKey := nftToken{
				Owner:      event.Guid.AccountAddress,
				Creator:    nftTokenData.Creator,
				Collection: nftTokenData.Collection,
				Name:       nftTokenData.Name,
				// Creator:    event.Data.Id.TokenDataId.Creator,
				// Collection: event.Data.Id.TokenDataId.Collection,
				// Name:       event.Data.Id.TokenDataId.Name,
			}
			asset, ok := saver.assetToken[ownerKey]
			if !ok {
				asset = &models.AssetToken{
					Version:    data.Version,
					Hash:       data.Hash,
					TxTime:     txTime,
					Owner:      event.Guid.AccountAddress,
					Creator:    nftTokenData.Creator,
					Collection: nftTokenData.Collection,
					Name:       nftTokenData.Name,
					// Creator:    event.Data.Id.TokenDataId.Creator,
					// Collection: event.Data.Id.TokenDataId.Collection,
					// Name:       event.Data.Id.TokenDataId.Name,
					Amount: event.Data.Amount,
				}
				saver.assetToken[ownerKey] = asset
			} else {
				asset.Version = data.Version
				asset.Hash = data.Hash
				asset.TxTime = txTime
				amount := models.ParseInt64(event.Data.Amount)
				asset.Amount = string(strconv.FormatInt(models.ParseInt64(asset.Amount)+amount, 10))
			}
		case EventTokenWithdraw:
			ownerKey := nftToken{
				Owner:      event.Guid.AccountAddress,
				Creator:    nftTokenData.Creator,
				Collection: nftTokenData.Collection,
				Name:       nftTokenData.Name,
				// Creator:    event.Data.Id.TokenDataId.Creator,
				// Collection: event.Data.Id.TokenDataId.Collection,
				// Name:       event.Data.Id.TokenDataId.Name,
			}
			asset, ok := saver.assetToken[ownerKey]
			if !ok {
				asset = &models.AssetToken{
					Version:    data.Version,
					Hash:       data.Hash,
					TxTime:     txTime,
					Owner:      event.Guid.AccountAddress,
					Creator:    nftTokenData.Creator,
					Collection: nftTokenData.Collection,
					Name:       nftTokenData.Name,
					// Creator:    event.Data.Id.TokenDataId.Creator,
					// Collection: event.Data.Id.TokenDataId.Collection,
					// Name:       event.Data.Id.TokenDataId.Name,
					Amount: event.Data.Amount,
				}
				saver.assetToken[ownerKey] = asset
			} else {
				asset.Version = data.Version
				asset.Hash = data.Hash
				asset.TxTime = txTime
				amount := models.ParseInt64(event.Data.Amount)
				asset.Amount = string(strconv.FormatInt(models.ParseInt64(asset.Amount)-amount, 10))
			}

		}
	}
}

func handlerHistoryToken(saver *DbSaver, txTime, sequenceNum int64, data models.TransactionRsp) {
	for _, event := range data.Events {
		var amount string
		var action int64
		var sender, receiver string = ZeroAddress, ZeroAddress
		var creator, collection, name string = "", "", ""

		var nftTokenData models.EventTokenDataId
		parseInterface(event.Data.Id, &nftTokenData)

		switch event.Type {
		case EventTokenDeposit:
			creator = nftTokenData.Creator
			collection = nftTokenData.Collection
			name = nftTokenData.Name
			// creator = event.Data.Id.TokenDataId.Creator
			// collection = event.Data.Id.TokenDataId.Collection
			// name = event.Data.Id.TokenDataId.Name
			receiver = event.Guid.AccountAddress
			amount = event.Data.Amount
		case EventTokenWithdraw:
			creator = nftTokenData.Creator
			collection = nftTokenData.Collection
			name = nftTokenData.Name
			// creator = event.Data.Id.TokenDataId.Creator
			// collection = event.Data.Id.TokenDataId.Collection
			// name = event.Data.Id.TokenDataId.Name
			sender = event.Guid.AccountAddress
			amount = event.Data.Amount
		}

		if sender == receiver && sender == ZeroAddress {
			continue
		} else if sender == ZeroAddress {
			action = ActionMint
		} else if receiver == ZeroAddress {
			action = ActionBurn
		} else {
			action = ActionTransfer
		}
		saver.historyToken = append(saver.historyToken, &models.HistoryToken{
			Version:    data.Version,
			Hash:       data.Hash,
			TxTime:     txTime,
			Sender:     sender,
			Receiver:   receiver,
			Creator:    creator,
			Collection: collection,
			Name:       name,
			Amount:     amount,
			Action:     action,
		})
	}
}
