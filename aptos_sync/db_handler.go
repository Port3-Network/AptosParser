package main

import (
	"strconv"

	"github.com/Port3-Network/AptosParser/models"
	oo "github.com/Port3-Network/liboo"
)

func handlerUserTransaction(db *DbSaver, data models.TransactionRsp) error {
	switch data.Payload.Type {
	case TypeCallFunction:
		txTime, _ := strconv.ParseInt(data.Timestamp, 10, 64)
		sequenceNum, _ := strconv.ParseInt(data.SequenceNumber, 10, 64)

		// payload -> done
		handlerPayload(db, txTime, sequenceNum, data)

		// transaction -> done
		handlerTx(db, txTime, sequenceNum, data)

		// record ->
		if data.Payload.Function == FunctionPublishPkg {
			handlerRecordCoin(db, txTime, sequenceNum, data)
			return nil
		}

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
	saver.payload = append(saver.payload, &models.Payload{
		Version:        data.Version,
		Hash:           data.Hash,
		TxTime:         txTime,
		SequenceNumber: sequenceNum,
		Sender:         data.Sender,
		PayloadFunc:    data.Payload.Function,
		PayloadType:    data.Payload.Type,
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
		switch event.Type {
		case EventDeposit:
			receiver = event.Guid.AccountAddress
			amount = event.Data.Amount
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
			amount = event.Data.Amount
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
	for _, change := range data.Changes {
		contract := ParseType(change.Data.Type)
		if contract == nil {
			continue
		}
		switch contract.Type {
		case ChangeTypeCoinInfo:
			saver.recordCoin = append(saver.recordCoin, &models.RecordCoin{
				Version:      data.Version,
				Hash:         data.Hash,
				TxTime:       txTime,
				Sender:       data.Sender,
				ModuleName:   contract.Module,
				ContractName: contract.Name,
				Resource:     contract.Resource,
				Name:         change.Data.Data.Name,
				Symbol:       change.Data.Data.Symbol,
			})
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

		saver.recordToken = append(saver.recordToken, &models.RecordToken{
			Version:     data.Version,
			Hash:        data.Hash,
			TxTime:      txTime,
			Sender:      data.Sender,
			Creator:     event.Data.Id.Creator,
			Collection:  event.Data.Id.Collection,
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
		switch event.Type {
		case EventTokenDeposit:
			ownerKey := nftToken{
				Owner:      event.Guid.AccountAddress,
				Creator:    event.Data.Id.TokenDataId.Creator,
				Collection: event.Data.Id.TokenDataId.Collection,
				Name:       event.Data.Id.TokenDataId.Name,
			}
			asset, ok := saver.assetToken[ownerKey]
			if !ok {
				asset = &models.AssetToken{
					Version:    data.Version,
					Hash:       data.Hash,
					TxTime:     txTime,
					Owner:      event.Guid.AccountAddress,
					Creator:    event.Data.Id.TokenDataId.Creator,
					Collection: event.Data.Id.TokenDataId.Collection,
					Name:       event.Data.Id.TokenDataId.Name,
					Amount:     event.Data.Amount,
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
				Creator:    event.Data.Id.TokenDataId.Creator,
				Collection: event.Data.Id.TokenDataId.Collection,
				Name:       event.Data.Id.TokenDataId.Name,
			}
			asset, ok := saver.assetToken[ownerKey]
			if !ok {
				asset = &models.AssetToken{
					Version:    data.Version,
					Hash:       data.Hash,
					TxTime:     txTime,
					Owner:      event.Guid.AccountAddress,
					Creator:    event.Data.Id.TokenDataId.Creator,
					Collection: event.Data.Id.TokenDataId.Collection,
					Name:       event.Data.Id.TokenDataId.Name,
					Amount:     event.Data.Amount,
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

		switch event.Type {
		case EventTokenDeposit:
			creator = event.Data.Id.TokenDataId.Creator
			collection = event.Data.Id.TokenDataId.Collection
			name = event.Data.Id.TokenDataId.Name
			receiver = event.Guid.AccountAddress
			amount = event.Data.Amount
		case EventTokenWithdraw:
			creator = event.Data.Id.TokenDataId.Creator
			collection = event.Data.Id.TokenDataId.Collection
			name = event.Data.Id.TokenDataId.Name
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
