package main

import (
	"strconv"

	"github.com/Port3-Network/AptosParser/models"
	oo "github.com/Port3-Network/liboo"
)

func handlerUserTransaction(db *DbSaver, data models.TransactionRsp) error {
	switch data.Payload.Type {
	case TypeCallFunction:
		version, _ := strconv.ParseInt(data.Version, 10, 64)
		txTime, _ := strconv.ParseInt(data.Timestamp, 10, 64)
		sequenceNum, _ := strconv.ParseInt(data.SequenceNumber, 10, 64)

		// payload -> done
		handlerPayload(db, version, txTime, sequenceNum, data)

		// transaction -> done
		handlerTx(db, version, txTime, sequenceNum, data)

		// record ->
		if data.Payload.Function == FunctionPublishPkg {
			handlerRecordCoin(db, version, txTime, sequenceNum, data)
			return nil
		}

		// history -> done
		handlerHistoryCoin(db, version, txTime, sequenceNum, data)

		handlerCollection(db, version, txTime, sequenceNum, data)
	default:
		oo.LogD("payload type [%s] not found", data.Payload.Type)
	}
	return nil
}

func handlerPayload(saver *DbSaver, version, txTime, sequenceNum int64, data models.TransactionRsp) {
	saver.payload = append(saver.payload, &models.Payload{
		Version:        version,
		Hash:           data.Hash,
		TxTime:         txTime,
		SequenceNumber: sequenceNum,
		Sender:         data.Sender,
		PayloadFunc:    data.Payload.Function,
		PayloadType:    data.Payload.Type,
	})
}

func handlerTx(saver *DbSaver, version, txTime, sequenceNum int64, data models.TransactionRsp) {
	saver.transaction = append(saver.transaction, &models.Transaction{
		Version:        version,
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

func handlerHistoryCoin(saver *DbSaver, version, txTime, sequenceNum int64, data models.TransactionRsp) {
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
			Version:  version,
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

func handlerRecordCoin(saver *DbSaver, version, txTime, sequenceNum int64, data models.TransactionRsp) {
	for _, change := range data.Changes {
		contract := ParseType(change.Data.Type)
		if contract == nil {
			continue
		}
		switch contract.Type {
		case ChangeTypeCoinInfo:
			saver.HandlerAddRecordToken(contract.Resource, &models.RecordCoin{
				Version:      version,
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

func handlerCollection(saver *DbSaver, version, txTime, sequenceNum int64, data models.TransactionRsp) {
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
			Version:     version,
			Hash:        data.Hash,
			TxTime:      txTime,
			Sender:      data.Sender,
			Creator:     event.Data.Creator,
			Name:        event.Data.CollectionName,
			Description: event.Data.Description,
			Uri:         event.Data.Uri,
			Maximun:     event.Data.Maximum,
			Type:        cType,
		})

	}
}
