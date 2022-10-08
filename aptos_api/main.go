package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/DeanThompson/ginpprof"
	oo "github.com/Port3-Network/liboo"
)

type DataBase struct {
	TxRpcUrl     []string `toml:"TX_RPC_URL,omitzero"`
	TxGraphqlUrl []string `toml:"TX_GRAPHQL_URL,omitzero"`
	User         string   `toml:"USER,omitzero"`
	Password     string   `toml:"PASSWORD,omitzero"`
	Host         string   `toml:"HOST,omitzero"`
	Port         int32    `toml:"PORT,omitzero"`
	Name         string   `toml:"NAME,omitzero"`
	ApiPort      int64    `toml:"API_PORT,omitzero"`
	EnableDebug  bool     `toml:"ENABLE_DEBUG,omitempry"`
}

var (
	GitVersion  string = "unknown"
	GWorkDir    string = ""
	GServerName string = ""
	GServerMark string = ""
	GConfig     *oo.Config
	GDatabase   *DataBase
	GNetwork    string
	GMysql      *oo.MysqlPool
	GRedis      *oo.RedisPool
)

func main() {
	defer func() {
		if err := recover(); nil != err {
			oo.LogW("panic err %v", err)
		}
	}()

	var err error

	flag.StringVar(&GNetwork, "n", "main", "main test dev")
	flag.Parse()

	GWorkDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	GServerName = strings.Split(filepath.Base(os.Args[0]), ".")[0]
	GServerMark = oo.GetSvrmark(GServerName)

	svrTag := GServerMark + "." + GitVersion
	oo.InitLog("./", GServerName, svrTag, func(str string) {})

	// config
	GConfig, err = oo.InitConfig(path.Join(GWorkDir, "../etc/config.conf"), nil)
	if err != nil {
		oo.LogW("Failed to load config. %v", err)
		return
	}
	if err = GConfig.SessDecode(GNetwork, &GDatabase); err != nil {
		oo.LogW("Decode config error. err=%v", err)
		return
	}

	// mysql
	GMysql, err = oo.InitMysqlPool(GDatabase.Host, GDatabase.Port, GDatabase.User, GDatabase.Password, GDatabase.Name)
	if err != nil {
		oo.LogW("Failed to init mysql. %v", err)
		return
	}

	router := InitAPIRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf("0.0.0.0:%d", GDatabase.ApiPort),
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
	}

	ginpprof.Wrap(router)

	oo.LogD("service run at %d", GDatabase.ApiPort)

	err = s.ListenAndServe()
	if err != nil {
		oo.LogW("ListenAndServe err %v", err)
		return
	}

}
