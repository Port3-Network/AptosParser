package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/Port3-Network/AptosParser/aptos_api/docs"

	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
)

const (
	SUCCESS        = 200
	INVALID_PARAMS = 400
	ERROR_DB_ERROR = 500
)

var MsgReturn = map[int]string{
	SUCCESS:        "ok",
	INVALID_PARAMS: "Invalid params",
	ERROR_DB_ERROR: "There is an error with service",
}

type Context struct {
	C *gin.Context
}

func (c *Context) ResponseInvalidParam() {
	c.Response(http.StatusBadRequest, INVALID_PARAMS, nil)
}

func (c *Context) Response(httpCode int, code int, data interface{}) {
	c.C.JSON(httpCode, gin.H{
		"code": code,
		"msg":  MsgReturn[code],
		"data": data,
	})
}

func InitAPIRouter() *gin.Engine {
	if !GDatabase.EnableDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	if GDatabase.EnableDebug {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
		r.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))
	}

	r.Use(Cros())
	v1Group := r.Group("/v1/")
	{
		v1Group.GET("/blocks", GetBlocks)
		v1Group.GET("/user_transactions", GetTransactions)

		v1Group.GET("/coin_inventory", GetCoinInventory)
		v1Group.GET("/coin_transactions", GetCoinTransactions)

		v1Group.GET("/get_address_action", GetAddressAction)
		v1Group.GET("/get_payload_detail", GetPayloadDetail)
		v1Group.GET("/get_address_amount", GetAddressAmount)
		v1Group.GET("/get_asset_token", GetAssetToken)

		v1Group.GET("/stats", GetStats)
	}

	return r
}

func Cros() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Apitoken, Authorization, Token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Headers")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
