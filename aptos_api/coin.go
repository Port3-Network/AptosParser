package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Port3-Network/AptosParser/models"
	oo "github.com/Port3-Network/liboo"
	"github.com/gin-gonic/gin"
)

type CoinInventoryReq struct {
	Offset   int64  `form:"offset" json:"offset" validate:"gte=0"`
	PageSize int64  `form:"pageSize" json:"pageSize" validate:"gt=0"`
	Resource string `form:"resource" json:"resource" validate:"omitempty"`
	Address  string `form:"address" json:"address" validate:"omitempty"`
}

type CoinInventoryRsp struct {
	List  []CoinInventoryJson `json:"list"`
	Total int64               `json:"total"`
}

type CoinInventoryJson struct {
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	ModuleName   string `json:"moduleName"`
	ContractName string `json:"contractName"`
	Resource     string `json:"resource"`
	Owner        string `json:"owner"`
}

// @Tags Coin
// @Summary get token inventory list
// @Description event = TokenInventory
// @Param body query CoinInventoryReq true "request"
// @Success 200 {object} CoinInventoryRsp
// @Router /v1/coin_inventory [get]
func GetCoinInventory(c *gin.Context) {
	appC := Context{C: c}
	req, rsp := &CoinInventoryReq{}, CoinInventoryRsp{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		oo.LogD("%s: ShouldBindQuery err, msg: %v", c.FullPath(), err)
		appC.ResponseInvalidParam()
		return
	}
	if err = oo.ValidateStruct(req); err != nil {
		oo.LogD("%s: Check para err %v", c.FullPath(), err)
		appC.ResponseInvalidParam()
		return
	}

	var data []struct {
		Name         string `db:"name"`
		Symbol       string `db:"symbol"`
		ModuleName   string `db:"module_name"`
		ContractName string `db:"contract_name"`
		Resource     string `db:"resource"`
		Owner        string `db:"sender"`
	}

	sqler := oo.NewSqler().Table(models.TableRecordCoin).
		Order("id DESC").
		Limit(int(req.PageSize)).
		Offset(int(req.Offset))
	if req.Resource != "" {
		sqler.Where("resource", req.Resource)
	}
	if req.Address != "" {
		sqler.Where("resource", req.Address)
	}
	sqlStr := sqler.Select("name,symbol,module_name,contract_name,resource,sender")
	if err = oo.SqlSelect(sqlStr, &data); err != nil {
		oo.LogD("%s: oo.SqlSelect err, msg: %v", c.FullPath(), err)
		appC.Response(http.StatusInternalServerError, ERROR_DB_ERROR, nil)
		return
	}

	for _, v := range data {
		rsp.List = append(rsp.List, CoinInventoryJson{
			Name:         v.Name,
			Symbol:       v.Symbol,
			ModuleName:   v.ModuleName,
			ContractName: v.ContractName,
			Resource:     v.Resource,
			Owner:        v.Owner,
		})
	}

	sqler2 := oo.NewSqler().Table(models.TableRecordCoin)
	if req.Resource != "" {
		sqler2.Where("resource", req.Resource)
	}
	if req.Address != "" {
		sqler2.Where("resource", req.Address)
	}
	sqlStr2 := sqler2.Select("COUNT(*) AS total")

	if err := oo.SqlGet(sqlStr2, &rsp.Total); err != nil {
		oo.LogD("%s: oo.SqlGet err, msg: %v", c.FullPath(), err)
		appC.Response(http.StatusInternalServerError, ERROR_DB_ERROR, nil)
		return
	}
	appC.Response(http.StatusOK, SUCCESS, rsp)
}

type HistoryCoinReq struct {
	Offset   int64  `form:"offset" json:"offset" validate:"gte=0"`
	PageSize int64  `form:"pageSize" json:"pageSize" validate:"gt=0"`
	Resource string `form:"resource" json:"contract" validate:"omitempty"`
	Address  string `form:"address" json:"address" validate:"omitempty"`
}

type HistoryCoinRsp struct {
	List  []HistoryTokenJson `json:"list"`
	Total int64              `json:"total"`
}

type HistoryTokenJson struct {
	Version  int64  `json:"version"`
	Hash     string `json:"hash"`
	TxTime   int64  `json:"tx_time"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Resource string `json:"resource"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Amount   string `json:"amount"`
	FuncName string `json:"func_name"`
}

// @Tags Coin
// @Summary get token tx list
// @Description event = tokenTransactions
// @Param body query HistoryCoinReq true "request"
// @Success 200 {object} HistoryCoinRsp
// @Router /v1/coin_transactions [get]
func GetCoinTransactions(c *gin.Context) {
	appC := Context{C: c}
	req, rsp := &HistoryCoinReq{}, &HistoryCoinRsp{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		oo.LogD("%s: ShouldBindQuery err, msg: %v", c.FullPath(), err)
		appC.ResponseInvalidParam()
		return
	}
	if err = oo.ValidateStruct(req); err != nil {
		oo.LogD("%s: Check para err %v", c.FullPath(), err)
		appC.ResponseInvalidParam()
		return
	}
	var data []struct {
		Version  int64          `db:"version"`
		Hash     string         `db:"hash"`
		TxTime   int64          `db:"tx_time"`
		Sender   string         `db:"sender"`
		Receiver string         `db:"receiver"`
		Resource string         `db:"resource"`
		Name     sql.NullString `db:"name"`
		Symbol   sql.NullString `db:"symbol"`
		Amount   string         `db:"amount"`
		FuncName sql.NullString `db:"payload_func"`
	}

	sql := oo.NewSqler().Table(models.TableHistoryCoin+" AS h").
		LeftJoin(models.TableRecordCoin+" AS r", "h.resource=r.resource").
		LeftJoin(models.TablePayload+" AS p", "h.version=p.version").
		Order("h.id DESC").
		Offset(int(req.Offset)).
		Limit(int(req.PageSize))
	if len(req.Resource) > 0 {
		sql.Where("h.resource", req.Resource)
	}
	if len(req.Address) > 0 {
		sqlOr := fmt.Sprintf("h.sender='%s' OR h.receiver='%s'", req.Address, req.Address)
		sql.Where(sqlOr)
	}

	sqlStr := sql.Select("h.version,h.hash,h.tx_time,h.sender,h.receiver,h.resource,h.amount,r.name,r.symbol,p.payload_func")
	if err = oo.SqlSelect(sqlStr, &data); err != nil {
		oo.LogD("%s: oo.SqlSelect1 err, msg: %v", c.FullPath(), err)
		appC.Response(http.StatusInternalServerError, ERROR_DB_ERROR, nil)
		return
	}
	for _, v := range data {
		rsp.List = append(rsp.List, HistoryTokenJson{
			Version:  v.Version,
			Hash:     v.Hash,
			TxTime:   v.TxTime,
			Sender:   v.Sender,
			Receiver: v.Receiver,
			Resource: v.Resource,
			Name:     v.Name.String,
			Symbol:   v.Symbol.String,
			Amount:   v.Amount,
			FuncName: v.FuncName.String,
		})
	}

	sql2 := oo.NewSqler().Table(models.TableHistoryCoin)
	if len(req.Resource) > 0 {
		sql2.Where("resource", req.Resource)
	}
	if len(req.Address) > 0 {
		sqlOr := fmt.Sprintf("sender='%s' OR receiver='%s'", req.Address, req.Address)
		sql2.Where(sqlOr)
	}

	sqlStr2 := sql2.Select("COUNT(*) AS total")
	if err := oo.SqlGet(sqlStr2, &rsp.Total); err != nil {
		oo.LogD("%s: oo.SqlGet err, msg: %v", c.FullPath(), err)
		appC.Response(http.StatusInternalServerError, ERROR_DB_ERROR, nil)
		return
	}
	appC.Response(http.StatusOK, SUCCESS, &rsp)
}
