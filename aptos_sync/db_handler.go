package main

import (
	"strconv"

	"github.com/Port3-Network/AptosParser/models"
)

func handlerUserTransaction(db *DbSaver, data models.TransactionRsp) error {
	var resource string
	version, _ := strconv.ParseInt(data.Version, 10, 64)
	txTime, _ := strconv.ParseInt(data.Timestamp, 10, 64)
	sequenceNum, _ := strconv.ParseInt(data.SequenceNumber, 10, 64)
	db.AddTransaction(&models.Transaction{
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

	db.AddPayload(&models.Payload{
		Version:        version,
		Hash:           data.Hash,
		SequenceNumber: sequenceNum,
		TxTime:         txTime,
		Sender:         data.Sender,
		PayloadFunc:    data.Payload.Function,
		PayloadType:    data.Payload.Type,
	})

	if len(data.Changes) > 0 {
		for _, c := range data.Changes {
			contract := ParseType(c.Data.Type)
			if contract == nil {
				continue
			}
			resource = contract.Resource
			if contract.Type == NewCoinInfo {
				db.HandlerAddRecordToken(contract.Resource, &models.RecordToken{
					Version:      version,
					Hash:         data.Hash,
					TxTime:       txTime,
					Sender:       data.Sender,
					ModuleName:   contract.Module,
					ContractName: contract.Name,
					Resource:     contract.Resource,
					Name:         c.Data.Data.Name,
					Symbol:       c.Data.Data.Symbol,
				})
			}
		}
	}

	// historyToken
	if len(data.Events) > 0 {
		var action int64 = ActionTransfer
		var sender string = ZeroAddress
		var receiver string = ZeroAddress
		amount := data.Events[0].Data.Amount

		// for _, e := range data.Events {

		// 	fmt.Printf("event: %v\n", e)
		// 	switch e.Type {
		// 	case EventWithdraw:
		// 		sender = data.Sender
		// 		// sender = "0x" + e.Key[18:]
		// 	case EventDeposit:
		// 		// receiver = "0x" + e.Key[18:]
		// 		receiver = data.Changes
		// 	}
		// }
		if data.Payload.Type == EventTransfer {
			sender = data.Sender
			receiver = data.Payload.Arguments[0].(string)
		}

		if receiver == ZeroAddress && sender == ZeroAddress {
			return nil
		}
		if sender == ZeroAddress {
			action = ActionMint
		}
		if receiver == ZeroAddress {
			action = ActionBurn
		}

		db.historyToken = append(db.historyToken, &models.HistoryToken{
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

	return nil
}
