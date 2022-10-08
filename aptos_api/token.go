package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Port3-Network/AptosParser/models"
	oo "github.com/Port3-Network/liboo"
	"github.com/gin-gonic/gin"
)

type GetActionReq struct {
	Address   string `form:"address"`
	Resource  string `form:"resource"`
	StartTime int64  `form:"start_time"`
	EndTime   int64  `form:"end_time"`
	Offset    int64  `form:"offset" json:"offset" validate:"gte=0"`
	PageSize  int64  `form:"pageSize" json:"pageSize" validate:"gt=0"`
}

type GetActionRsp struct {
	List  []ActionData `json:"list"`
	Total int64        `json:"total"`
}

type ActionData struct {
	Version  string `json:"version"`
	Hash     string `json:"hash"`
	TxTime   int64  `json:"tx_time"`
	Sender   string `json:"sender"`
	FuncName string `json:"function_name"`
	Resource string `json:"resource"`
}

// GetTokenAction
//
// @Summary get user action
// @Id GetTokenAction
// @Description event = blocks
// @Tags Block
// @Accept json
// @Produce json
// @Param request body GetActionReq true "request"
// @Success 200 {object} GetActionRsp
// @Router /get_token_action [get]
func GetTokenAction(c *gin.Context) {
	appC := Context{C: c}
	req, rsp := &GetActionReq{}, &GetActionRsp{}
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
		Version  string         `db:"version"`
		Hash     string         `db:"hash"`
		TxTime   int64          `db:"tx_time"`
		Sender   string         `db:"sender"`
		FuncName string         `db:"function_name"`
		Resource sql.NullString `db:"resource"`
	}

	sqler := oo.NewSqler().Table(models.TablePayload+" AS p").
		LeftJoin(models.TableHistoryToken+" AS h", "p.version=h.version").
		Limit(int(req.PageSize)).
		Offset(int(req.Offset)).
		Order("p.tx_time DESC")

	if req.Address != "" {
		sqler.Where("p.sender", req.Address)
	}
	if req.Resource != "" {
		sqler.Where("h.resource", req.Resource)
	}
	if req.StartTime != 0 {
		startWhere := fmt.Sprintf("p.tx_time >= %d", req.StartTime*1000)
		sqler.Where(startWhere)
	}
	if req.EndTime != 0 {
		endWhere := fmt.Sprintf("p.tx_time <= %d", req.EndTime*1000)
		sqler.Where(endWhere)
	}
	sqlStr1 := sqler.Select("p.version,p.hash,p.tx_time,p.sender,p.payload_func AS function_name,h.resource")
	if err = oo.SqlSelect(sqlStr1, &data); err != nil {
		oo.LogD("%s: oo.SqlSelect err, msg: %v", c.FullPath(), err)
		appC.Response(http.StatusInternalServerError, ERROR_DB_ERROR, nil)
		return
	}
	for _, v := range data {
		rsp.List = append(rsp.List, ActionData{
			Version:  v.Version,
			Hash:     v.Hash,
			TxTime:   v.TxTime,
			Sender:   v.Sender,
			FuncName: v.FuncName,
			Resource: v.Resource.String,
		})
	}

	sqler2 := oo.NewSqler().Table(models.TablePayload+" AS p").
		LeftJoin(models.TableHistoryToken+" AS h", "p.version=h.version")

	if req.Address != "" {
		sqler.Where("p.sender", req.Address)
	}
	if req.Resource != "" {
		sqler.Where("h.resource", req.Resource)
	}
	if req.StartTime != 0 {
		startWhere := fmt.Sprintf("p.tx_time >= %d", req.StartTime*1000)
		sqler.Where(startWhere)
	}
	if req.EndTime != 0 {
		endWhere := fmt.Sprintf("p.tx_time <= %d", req.EndTime*1000)
		sqler.Where(endWhere)
	}
	sqlStr2 := sqler2.Select("COUNT(*) AS total")

	if err = oo.SqlGet(sqlStr2, &rsp.Total); err != nil {
		oo.LogD("%s: oo.SqlGet err, msg: %v", c.FullPath(), err)
		appC.Response(http.StatusInternalServerError, ERROR_DB_ERROR, nil)
		return
	}
	appC.Response(http.StatusOK, SUCCESS, rsp)
}

type TokenInventoryReq struct {
	Offset   int64  `form:"offset" json:"offset" validate:"gte=0"`
	PageSize int64  `form:"pageSize" json:"pageSize" validate:"gt=0"`
	Resource string `form:"resource" json:"resource" validate:"omitempty"`
	Address  string `form:"address" json:"address" validate:"omitempty"`
}

type TokenInventoryRsp struct {
	List  []TokenInventoryJson `json:"list"`
	Total int64                `json:"total"`
}

type TokenInventoryJson struct {
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	ModuleName   string `json:"moduleName"`
	ContractName string `json:"contractName"`
	Resource     string `json:"resource"`
	Owner        string `json:"owner"`
}

// GetTokenInventory
// @Summary get token inventory list
// @Id TokenInventory
// @Description event = TokenInventory
// @Tags Token
// @Accept json
// @Produce json
// @Param request body TokenInventoryReq true "request"
// @Success 200 {object} TokenInventoryRsp
// @Router /token_inventory [get]
func GetTokenInventory(c *gin.Context) {
	appC := Context{C: c}
	req, rsp := &TokenInventoryReq{}, TokenInventoryRsp{}
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

	sqler := oo.NewSqler().Table(models.TableRecordToken).
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
		rsp.List = append(rsp.List, TokenInventoryJson{
			Name:         v.Name,
			Symbol:       v.Symbol,
			ModuleName:   v.ModuleName,
			ContractName: v.ContractName,
			Resource:     v.Resource,
			Owner:        v.Owner,
		})
	}

	sqler2 := oo.NewSqler().Table(models.TableRecordToken)
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

type HistoryTokenReq struct {
	Offset   int64  `form:"offset" json:"offset" validate:"gte=0"`
	PageSize int64  `form:"pageSize" json:"pageSize" validate:"gt=0"`
	Resource string `form:"resource" json:"contract" validate:"omitempty"`
	Address  string `form:"address" json:"address" validate:"omitempty"`
}

type HistoryTokenRsp struct {
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

// GetTokenTransactions
// @Summary get token tx list
// @Id tokenTransactions
// @Description event = tokenTransactions
// @Tags Token
// @Accept json
// @Produce json
// @Param request body HistoryTokenReq true "request"
// @Success 200 {object} HistoryTokenRsp
// @Router /token_transactions [get]
func GetTokenTransactions(c *gin.Context) {
	appC := Context{C: c}
	req, rsp := &HistoryTokenReq{}, &HistoryTokenRsp{}
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

	sql := oo.NewSqler().Table(models.TableHistoryToken+" AS h").
		LeftJoin(models.TableRecordToken+" AS r", "h.resource=r.resource").
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

	sql2 := oo.NewSqler().Table(models.TableHistoryToken)
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
