package main

import (
	"net/http"

	"github.com/Port3-Network/AptosParser/models"
	oo "github.com/Port3-Network/liboo"
	"github.com/gin-gonic/gin"
)

type BlocksReq struct {
	Height   int64 `form:"height" validate:"gte=0"`                  // optional, block num
	Offset   int64 `form:"offset" json:"offset" validate:"gte=0"`    // required, data offset
	PageSize int64 `form:"pageSize" json:"pageSize" validate:"gt=0"` // required, number of data a time
}

type BlocksRsp struct {
	List  []BlockData `json:"list"`  // data list
	Total int64       `json:"total"` // total num
}

type BlockData struct {
	Height       int64  `json:"height"`        // block height
	Hash         string `json:"hash"`          // block hash
	TxTime       int64  `json:"tx_time"`       // block timestamp
	FirstVersion string `json:"first_version"` // the first version of this block
	LastVersion  string `json:"last_version"`  // the last version of this block
}

// @Tags Tx
// @Summary get block list
// @Description event = blocks
// @Param body query BlocksReq true "request"
// @Success 200 {object} BlocksRsp
// @Router /v1/blocks [get]
func GetBlocks(c *gin.Context) {
	appC := Context{C: c}
	req, rsp := &BlocksReq{}, &BlocksRsp{}
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

	// Get data
	var data []models.Block
	sqler := oo.NewSqler().Table(models.TableBlock).
		Order("height DESC").
		Limit(int(req.PageSize)).
		Offset(int(req.Offset))
	if req.Height > 0 {
		sqler.Where("height", req.Height)
	}
	sqlStr := sqler.Select("*")

	// call mysql -> oo.SqlSelect use *sqlx.DB.Select
	if err = oo.SqlSelect(sqlStr, &data); err != nil {
		oo.LogD("%s: oo.SqlSelect err, msg: %v", c.FullPath(), err)
		appC.Response(http.StatusInternalServerError, ERROR_DB_ERROR, nil)
		return
	}

	for _, v := range data {
		rsp.List = append(rsp.List, BlockData{
			Height:       v.Height,
			Hash:         v.Hash,
			TxTime:       v.BlockTime,
			FirstVersion: v.FirstVersion,
			LastVersion:  v.LastVersion,
		})
	}

	// Get Count
	sqler2 := oo.NewSqler().Table(models.TableBlock)
	if req.Height > 0 {
		sqler2.Where("height", req.Height)
	}
	// call mysql -> oo.sqlGet use *sqlx.DB.Get
	sqlStr2 := sqler2.Select("COUNT(*) AS total")
	if err = oo.SqlGet(sqlStr2, &rsp.Total); err != nil {
		oo.LogD("%s: oo.SqlSelect err, msg: %v", c.FullPath(), err)
		appC.Response(http.StatusInternalServerError, ERROR_DB_ERROR, nil)
		return
	}

	appC.Response(http.StatusOK, SUCCESS, rsp)
}

type UserTransactionsReq struct {
	Version  string `form:"version" validate:"omitempty"`             // optional, version num
	Offset   int64  `form:"offset" json:"offset" validate:"gte=0"`    // required, data offset
	PageSize int64  `form:"pageSize" json:"pageSize" validate:"gt=0"` // required, number of data a time
}

type UserTransactionsRsp struct {
	List  []UserTransactionJson `json:"list"`  // data list
	Total int64                 `json:"total"` // total num
}

type UserTransactionJson struct {
	Id       int64  `json:"id"`
	Version  string `json:"version"`  // tx version
	Hash     string `json:"hash"`     // tx hash
	TxTime   int64  `json:"tx_time"`  // tx timestamp
	Success  bool   `json:"success"`  // tx status, just success
	Sender   string `json:"sender"`   // tx sender
	Function string `json:"function"` // call function
}

// @Tags Tx
// @Summary get tx list
// @Description event = transactions
// @Param body query UserTransactionsReq true "request"
// @Success 200 {object} UserTransactionsRsp
// @Router /v1/user_transactions [get]
func GetTransactions(c *gin.Context) {
	appC := Context{C: c}
	req, rsp := &UserTransactionsReq{}, &UserTransactionsRsp{}
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
		Id       int64  `json:"id"`
		Version  string `json:"version"`
		Hash     string `json:"hash"`
		TxTime   int64  `json:"tx_time"`
		Success  bool   `json:"success"`
		Sender   string `json:"sender"`
		CallFunc string `json:"call_func"`
	}

	sqler := oo.NewSqler().Table(models.TableTransaction+" AS t").
		LeftJoin(models.TablePayload+" AS p", "p.hash=t.hash").
		Order("t.tx_time DESC").
		Limit(int(req.PageSize)).
		Offset(int(req.Offset))
	if req.Version != "" {
		sqler.Where("t.version", req.Version)
	}
	sqlStr := sqler.Select("t.id,t.version,t.hash,t.tx_time,t.success,t.sender,p.payload_func AS call_func")
	oo.LogD("%s: sql: %v\n", c.FullPath(), sqlStr)
	if err = oo.SqlSelect(sqlStr, &data); err != nil {
		oo.LogD("%s: oo.SqlSelect err, msg: %v", c.FullPath(), err)
		appC.Response(http.StatusInternalServerError, ERROR_DB_ERROR, nil)
		return
	}
	for _, v := range data {
		rsp.List = append(rsp.List, UserTransactionJson{
			Id:       v.Id,
			Version:  v.Version,
			Hash:     v.Hash,
			TxTime:   v.TxTime,
			Success:  v.Success,
			Sender:   v.Sender,
			Function: v.CallFunc,
		})
	}

	sqler2 := oo.NewSqler().Table(models.TableTransaction)
	if req.Version != "" {
		sqler2.Where("version", req.Version)
	}
	sqlStr2 := sqler2.Select("COUNT(*) AS total")
	oo.LogD("%s: sql2: %v\n", c.FullPath(), sqlStr2)
	if err = oo.SqlGet(sqlStr2, &rsp.Total); err != nil {
		oo.LogD("%s: oo.SqlSelect err, msg: %v", c.FullPath(), err)
		appC.Response(http.StatusInternalServerError, ERROR_DB_ERROR, nil)
		return
	}
	appC.Response(http.StatusOK, SUCCESS, rsp)
}

type SyncStatsRsp struct {
	CurrentVersion string `json:"current_version"`
}

// @Tags Tx
// @Summary get system stats
// @Description event = current block stats of sync
// @Success 200 {object} SyncStatsRsp
// @Router /v1/stats [get]
func GetStats(c *gin.Context) {
	appC := Context{C: c}
	rsp := &SyncStatsRsp{}
	var err error
	rsp.CurrentVersion, err = GetSysConfig(SyncLatestBlock)
	if err != nil {
		oo.LogD("%s: GetSyncBlockNum err, msg: %v", c.FullPath(), err)
	}

	appC.Response(http.StatusOK, SUCCESS, rsp)
}

func GetSysConfig(key string) (str string, err error) {
	conn := GMysql.GetConn()
	defer GMysql.UnGetConn(conn)

	sqlstr := oo.NewSqler().Table(models.TableSysconfig).
		Where("cfg_name", key).
		Select("cfg_val")

	err = conn.Get(&str, sqlstr)
	if nil != err && oo.ErrNoRows != err {
		err = oo.NewError("failed to sql[%s] err[%v]", sqlstr, err)
		return
	}
	return
}
