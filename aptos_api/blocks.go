package main

import (
	"net/http"

	"github.com/Port3-Network/AptosParser/models"
	oo "github.com/Port3-Network/liboo"
	"github.com/gin-gonic/gin"
)

type BlocksReq struct {
	Height   int64 `form:"height" validate:"gte=0"`
	Offset   int64 `form:"offset" json:"offset" validate:"gte=0"`
	PageSize int64 `form:"pageSize" json:"pageSize" validate:"gt=0"`
}

type BlocksRsp struct {
	List  []BlockData `json:"list"`
	Total int64       `json:"total"`
}

type BlockData struct {
	Height       int64  `json:"height"`
	Hash         string `json:"hash"`
	TxTime       int64  `json:"tx_time"`
	FirstVersion string `json:"first_version"`
	LastVersion  string `json:"last_version"`
}

// GetBlocks
//
// @Summary get block list
// @Id blocks
// @Description event = blocks
// @Tags Block
// @Accept json
// @Produce json
// @Param request body BlocksReq true "request"
// @Success 200 {object} BlocksRsp
// @Router /blocks [get]
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
	sqlStr2 := sqler2.Select("COUNT(*) AS total")
	if err = oo.SqlGet(sqlStr2, &rsp.Total); err != nil {
		oo.LogD("%s: oo.SqlSelect err, msg: %v", c.FullPath(), err)
		appC.Response(http.StatusInternalServerError, ERROR_DB_ERROR, nil)
		return
	}

	appC.Response(http.StatusOK, SUCCESS, rsp)
}

type UserTransactionsReq struct {
	Version  string `form:"version" validate:"omitempty"`
	Offset   int64  `form:"offset" json:"offset" validate:"gte=0"`
	PageSize int64  `form:"pageSize" json:"pageSize" validate:"gt=0"`
}

type UserTransactionsRsp struct {
	List  []UserTransactionJson `json:"list"`
	Total int64                 `json:"total"`
}

type UserTransactionJson struct {
	Id       int64  `json:"id"`
	Version  int64  `json:"version"`
	Hash     string `json:"hash"`
	TxTime   int64  `json:"tx_time"`
	Success  bool   `json:"success"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
}

// GetUserTransactions
//
// @Summary get tx detail
// @Id transactions
// @Description event = transactions
// @Tags Transaction
// @Accept json
// @Produce json
// @Param request body UserTransactionsReq true "request"
// @Success 200 {object} UserTransactionsRsp
// @Router /user_transactions [get]
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

	var data []models.Transaction
	sqler := oo.NewSqler().Table(models.TableTransaction).
		Order("version DESC").
		Limit(int(req.PageSize)).
		Offset(int(req.Offset))
	if req.Version != "" {
		sqler.Where("version", req.Version)
	}
	sqlStr := sqler.Select("*")
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
			Receiver: v.Receiver,
		})
	}

	sqler2 := oo.NewSqler().Table(models.TableTransaction)
	if req.Version != "" {
		sqler2.Where("version", req.Version)
	}
	sqlStr2 := sqler2.Select("COUNT(*) AS total")
	if err = oo.SqlGet(sqlStr2, &rsp.Total); err != nil {
		oo.LogD("%s: oo.SqlSelect err, msg: %v", c.FullPath(), err)
		appC.Response(http.StatusInternalServerError, ERROR_DB_ERROR, nil)
		return
	}
	appC.Response(http.StatusOK, SUCCESS, rsp)
}
