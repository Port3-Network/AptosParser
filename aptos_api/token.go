package main

import (
	"database/sql"
	"net/http"
	u "net/url"
	"strings"

	"github.com/Port3-Network/AptosParser/models"
	oo "github.com/Port3-Network/liboo"
	"github.com/gin-gonic/gin"
)

type AssetTokensReq struct {
	Address      string `form:"address" json:"address" validate:"omitempty"`             // optional, user address
	CollectionId string `form:"collection_id" json:"collection_id" validate:"omitempty"` // optional, collection id -> creator+collection name
	TokenId      string `form:"token_id" json:"token_id" validate:"omitempty"`           // optional, token_id -> token name, requires collection id
	Offset       int64  `form:"offset" json:"offset" validate:"gte=0"`                   // required, data offset
	PageSize     int64  `form:"pageSize" json:"pageSize" validate:"gt=0"`                // required, number of data a time
}

type AssetTokensRsp struct {
	List  []AssetTokenJson `json:"list"`  // data list
	Total int64            `json:"total"` // total num
}

type AssetTokenJson struct {
	Id                    int64  `json:"id"`
	Version               string `json:"version"` // tx version
	Hash                  string `json:"hash"`    // tx hash
	TxTime                int64  `json:"tx_time"` // tx timestamp
	Owner                 string `json:"ownder"`  // asset owner
	CollectionCreator     string `json:"collection_creator"`
	CollectionName        string `json:"collection_name"`
	CollectionDescription string `json:"collection_description"`
	CollectionUri         string `json:"collection_uri"`
	TokenName             string `json:"token_name"`
	TokenDescription      string `json:"token_description"`
	Amount                string `json:"amount"`
	TokenUri              string `json:"token_uri"`
}

type nftToken struct {
	Creator    string
	Collection string
	Name       string
}

// @Tags Token
// @Summary get asset token list
// @Description Get the list of tokens attributed to the user
// @Param body query AssetTokensReq true "request"
// @Success 200 {object} AssetTokensRsp
// @Router /v1/get_asset_token [get]
func GetAssetToken(c *gin.Context) {
	appC := Context{C: c}
	req, rsp := &AssetTokensReq{}, &AssetTokensRsp{}
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
	nftId, ok := chkAssetTokenReq(*req)
	if !ok {
		oo.LogD("%s: chkAssetTokenReq err %v", c.FullPath(), err)
		appC.ResponseInvalidParam()
		return
	}

	var data []struct {
		Id                int64  `json:"id"`
		Version           string `json:"version"` // tx version
		Hash              string `json:"hash"`    // tx hash
		TxTime            int64  `json:"tx_time"` // tx timestamp
		Owner             string `json:"ownder"`  // asset owner
		CollectionCreator string `json:"collection_creator"`
		CollectionName    string `json:"collection_name"`
		TokenName         string `json:"token_name"`
		Amount            string `json:"amount"`
	}
	sqler := oo.NewSqler().Table(models.TableAssetToken).
		Offset(int(req.Offset)).
		Limit(int(req.PageSize)).
		Order("id DESC")

	if req.Address != "" {
		sqler.Where("owner", req.Address)
	}
	if nftId.Name != "" {
		sqler.Where("name", u.QueryEscape(nftId.Name))
	}
	if nftId.Creator != "" && nftId.Collection != "" {
		sqler.Where("creator", nftId.Creator).Where("collection", u.QueryEscape(nftId.Collection))
	}
	sqlStr := sqler.Select("id,version,hash,tx_time,owner,creator AS collection_creator,collection AS collection_name,name AS token_name,amount")
	oo.LogD("%s: sqlStr: %v\n", sqlStr)
	if err = oo.SqlSelect(sqlStr, &data); err != nil {
		oo.LogD("%s: oo.SqlSelect err, msg: %v", c.FullPath(), err)
		appC.Response(http.StatusInternalServerError, ERROR_DB_ERROR, nil)
		return
	}
	type collectionInfo struct {
		Creator string
		Name    string
		NFTName string
	}
	type NftInfo struct {
		CollectionName        sql.NullString `json:"collection_name"`
		CollectionDescription sql.NullString `json:"collection_description"`
		CollectionUri         sql.NullString `json:"collection_uri"`
		TokenDescription      sql.NullString `json:"token_description"`
		TokenUri              sql.NullString `json:"token_uri"`
	}
	nft := make(map[collectionInfo]*NftInfo)

	for _, v := range data {
		cInfo := collectionInfo{
			Creator: v.CollectionCreator,
			Name:    v.CollectionName,
			NFTName: v.TokenName,
		}
		if _, ok := nft[cInfo]; !ok {
			resData := NftInfo{}
			innSql := oo.NewSqler().Table(models.TableRecordToken+" AS r").
				LeftJoin(models.TableCollection+" AS c", "r.creator=c.creator AND r.collection=c.name").
				Where("r.creator", v.CollectionCreator).
				Where("r.collection", v.CollectionName).
				Where("r.name", v.TokenName).
				Select("c.name AS collection_name,c.description AS collection_description,c.uri AS collection_uri,r.description AS token_description,r.uri AS token_uri")
			oo.LogD("innsql: %v\n", innSql)
			if err := oo.SqlGet(innSql, &resData); err != nil {
				oo.LogD("%s: oo.SqlGet err, msg: %v", c.FullPath(), err)
			}
			nft[cInfo] = &resData
		}
		colName, _ := u.QueryUnescape(v.CollectionName)
		colDesc, _ := u.QueryUnescape(nft[cInfo].CollectionDescription.String)
		tokenName, _ := u.QueryUnescape(v.TokenName)
		tokenDesc, _ := u.QueryUnescape(nft[cInfo].TokenDescription.String)
		tokenUri, _ := u.QueryUnescape(nft[cInfo].TokenUri.String)
		collectionUri, _ := u.QueryUnescape(nft[cInfo].CollectionUri.String)
		rsp.List = append(rsp.List, AssetTokenJson{
			Id:                    v.Id,
			Version:               v.Version,
			Hash:                  v.Hash,
			TxTime:                v.TxTime,
			Owner:                 v.Owner,
			CollectionCreator:     v.CollectionCreator,
			CollectionName:        colName,
			CollectionDescription: colDesc,
			TokenName:             tokenName,
			TokenDescription:      tokenDesc,
			Amount:                v.Amount,
			TokenUri:              tokenUri,
			CollectionUri:         collectionUri,
		})
	}

	sqler2 := oo.NewSqler().Table(models.TableAssetToken)
	if req.Address != "" {
		sqler2.Where("owner", req.Address)
	}
	if nftId.Name != "" {
		sqler2.Where("name", u.QueryEscape(nftId.Name))
	}
	if nftId.Creator != "" && nftId.Collection != "" {
		sqler2.Where("creator", nftId.Creator).Where("collection", u.QueryEscape(nftId.Collection))
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

func chkAssetTokenReq(req AssetTokensReq) (nft *nftToken, ret bool) {
	nft = &nftToken{}
	ret = false
	if req.TokenId != "" {
		if req.CollectionId == "" {
			return nft, false
		}
		nft.Name = req.TokenId
	}
	if req.CollectionId != "" {
		cols := strings.SplitN(req.CollectionId, "_", 2)
		if len(cols) != 2 {
			return nft, false
		}

		nft.Creator = cols[0]
		nft.Collection = cols[1]
	}
	return nft, true
}
