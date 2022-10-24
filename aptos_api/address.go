package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Port3-Network/AptosParser/models"
	oo "github.com/Port3-Network/liboo"
	"github.com/gin-gonic/gin"
)

type GetActionReq struct {
	Address   string `form:"address"`                                  // optional, user address
	Resource  string `form:"resource"`                                 // optional, resource name
	FuncName  string `form:"funcName"`                                 // optional, call function
	StartTime int64  `form:"startTime"`                                // optional, begin time
	EndTime   int64  `form:"endTime"`                                  // optional, end time
	Offset    int64  `form:"offset" json:"offset" validate:"gte=0"`    // required, data offset
	PageSize  int64  `form:"pageSize" json:"pageSize" validate:"gt=0"` // required, number of data a time
}

type GetActionRsp struct {
	List  []ActionData `json:"list"`  // data list
	Total int64        `json:"total"` // total num
}

type ActionData struct {
	Version  string `json:"version"`       // tx version
	Hash     string `json:"hash"`          // tx hash
	TxTime   int64  `json:"tx_time"`       // tx timestamp
	Sender   string `json:"sender"`        // tx sender
	FuncName string `json:"function_name"` // call function
	Resource string `json:"resource"`      // which resource
}

// @Tags Address
// @Summary get user action
// @Description get address action
// @Param body query GetActionReq true "request"
// @Success 200 {object} GetActionRsp
// @Router /v1/get_address_action [get]
func GetAddressAction(c *gin.Context) {
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
	fmt.Printf("req: %v\n", req)
	sqler := oo.NewSqler().Table(models.TablePayload+" AS p").
		LeftJoin(models.TableHistoryCoin+" AS h", "p.version=h.version").
		LeftJoin(models.TableHistoryToken+" AS ht", "p.version=ht.version").
		Limit(int(req.PageSize)).
		Offset(int(req.Offset)).
		Order("p.id DESC")
		// Order("p.tx_time DESC")

	if req.Address != "" {
		sqler.Where("p.sender", req.Address)
		history_addr := fmt.Sprintf("h.sender='%s' or h.receiver='%s' or ht.sender='%s' or ht.receiver='%s'", req.Address, req.Address, req.Address, req.Address)
		sqler.Where(history_addr)
	}
	if req.Resource != "" {
		resLike := fmt.Sprintf("h.resource like '%s%%'", req.Resource)
		sqler.Where(resLike)
		// sqler.Where("h.resource", req.Resource)
	}
	if req.FuncName != "" {
		funcLike := fmt.Sprintf("p.payload_func like '%s%%'", req.FuncName)
		sqler.Where(funcLike)
		// sqler2.Where("p.payload_func", req.FuncName)
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
	oo.LogD("%s: sqlStr1: %v\n", c.FullPath(), sqlStr1)
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
		LeftJoin(models.TableHistoryCoin+" AS h", "p.version=h.version").
		LeftJoin(models.TableHistoryToken+" AS ht", "p.version=ht.version")

	if req.Address != "" {
		sqler2.Where("p.sender", req.Address)
		history_addr := fmt.Sprintf("h.sender='%s' or h.receiver='%s' or ht.sender='%s' or ht.receiver='%s'", req.Address, req.Address, req.Address, req.Address)
		sqler2.Where(history_addr)
	}
	if req.Resource != "" {
		resLike := fmt.Sprintf("h.resource like '%s%%'", req.Resource)
		sqler2.Where(resLike)
		// sqler2.Where("h.resource", req.Resource)
	}
	if req.FuncName != "" {
		funcLike := fmt.Sprintf("p.payload_func like '%s%%'", req.FuncName)
		sqler2.Where(funcLike)
		// sqler2.Where("p.payload_func", req.FuncName)
	}
	if req.StartTime != 0 {
		startWhere := fmt.Sprintf("p.tx_time >= %d", req.StartTime*1000)
		sqler2.Where(startWhere)
	}
	if req.EndTime != 0 {
		endWhere := fmt.Sprintf("p.tx_time <= %d", req.EndTime*1000)
		sqler2.Where(endWhere)
	}
	sqlStr2 := sqler2.Select("COUNT(*) AS total")
	oo.LogD("%s: sqlStr2: %v\n", c.FullPath(), sqlStr2)

	if err = oo.SqlGet(sqlStr2, &rsp.Total); err != nil {
		oo.LogD("%s: oo.SqlGet err, msg: %v", c.FullPath(), err)
		appC.Response(http.StatusInternalServerError, ERROR_DB_ERROR, nil)
		return
	}
	appC.Response(http.StatusOK, SUCCESS, rsp)
}

type GetAmountReq struct {
	Address  string `form:"address"`  // user address
	Resource string `form:"resource"` // which resource, default 0x1::aptos_coin::AptosCoin
}

type GetAmountRsp struct {
	Amount string `json:"amount"` // amount
}

// @Tags Address
// @Summary get user action
// @Description event = blocks
// @Param body query GetAmountReq true "request"
// @Success 200 {object} GetAmountRsp
// @Router /v1/get_address_amount [get]
func GetAddressAmount(c *gin.Context) {
	appC := Context{C: c}
	req, rsp := &GetAmountReq{}, &GetAmountRsp{}
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

	rsp.Amount = "0"
	// call api
	res, err := GetAccountResource(req.Address)
	if err != nil {
		oo.LogD("%s: GetAccountResource err %v", c.FullPath(), err)
		appC.Response(http.StatusOK, SUCCESS, rsp)
		return
	}

	for _, r := range *res {
		c := ParseType(r.Type)
		if c == nil || c.Type != TypeCoin {
			continue
		}

		if c.Resource == req.Resource {
			rsp.Amount = r.Data.Coin.Value
			break
		}

		if req.Resource == "" && c.Resource == NativeAptosCoin {
			rsp.Amount = r.Data.Coin.Value
			break
		}
	}

	appC.Response(http.StatusOK, SUCCESS, rsp)
}

func GetAccountResource(address string) (*[]models.ResourceRsp, error) {
	r := &[]models.ResourceRsp{}
	url := fmt.Sprintf("%s/accounts/%s/resources", GRpc, address)
	buf, _, err := models.HttpGet(url, 2)
	if err != nil {
		updateRpc()
		return r, fmt.Errorf("resource HttpGet msg: %v", err)
	}

	err = json.Unmarshal(buf, &r)
	if err != nil {
		return r, fmt.Errorf("resource Unmarshal msg: %v", err)
	}
	return r, nil
}

type Contract struct {
	Type     string
	Address  string
	Module   string
	Name     string
	Resource string
}

func ParseType(e string) *Contract {
	// parse string to struct
	// 0x1::coin::CoinStore<0xa1dffeb39031fbab2ae3cbbd2c59fcfcca1cd2eb3b80f521974d38e3de0c96e6::moon_coin::MoonCoin>
	/*
		cc := &Contract{
			Type:     "0x1::coin::CoinStore",
			Address:  "0xa1dffeb39031fbab2ae3cbbd2c59fcfcca1cd2eb3b80f521974d38e3de0c96e6",
			Module:   "moon_coin",
			Name:     "MoonCoin",
			Resource: "0xa1dffeb39031fbab2ae3cbbd2c59fcfcca1cd2eb3b80f521974d38e3de0c96e6::moon_coin::MoonCoin",
		}
	*/

	c := &Contract{}
	es := strings.Replace(strings.Replace(e, "<", "-", 1), ">", "", 1)
	psType := strings.Split(es, "-")
	if len(psType) < 2 {
		return nil
	}
	c.Type = psType[0]
	c.Resource = psType[1]
	eArr := strings.Split(psType[1], "::")
	if len(eArr) != 3 {
		return nil
	}
	c.Address = eArr[0]
	c.Module = eArr[1]
	c.Name = eArr[2]
	return c
}
