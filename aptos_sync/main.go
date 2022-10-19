package main

import (
	"errors"
	"flag"
	"os"
	"path"
	"path/filepath"
	"strings"

	oo "github.com/Port3-Network/liboo"
)

type DataBase struct {
	TxRpcUrl      []string `toml:"TX_RPC_URL,omitzero"`
	User          string   `toml:"USER,omitzero"`
	Password      string   `toml:"PASSWORD,omitzero"`
	Host          string   `toml:"HOST,omitzero"`
	Port          int32    `toml:"PORT,omitzero"`
	Name          string   `toml:"NAME,omitzero"`
	BlockCount    int64    `toml:"BLOCK_COUNT,omitzero"`
	RedisHost     string   `toml:"REDIS_HOST,omitzero"`
	RedisPort     int32    `toml:"REDIS_PORT,omitzero"`
	RedisPassword string   `toml:"REDIS_PASSWORD,omitzero"`
}

type rpcStatus struct {
	Url       string
	CoolDown  bool
	FailCount int64
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
	GRpc        string
	RpcMap      map[string]*rpcStatus = make(map[string]*rpcStatus)
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

	GRedis, err = oo.InitRedisPool(GDatabase.RedisHost, GDatabase.RedisPort, GDatabase.RedisPassword)
	if err != nil {
		oo.LogW("Failed to init redis. %v", err)
		return
	}

	err = initRpc()
	if err != nil {
		oo.LogW("Failed to init rpc. %v", err)
		return
	}

	FullSync()
}

func initRpc() (err error) {
	if len(GDatabase.TxRpcUrl) < 1 {
		return errors.New("rpc not found")
	}

	for _, v := range GDatabase.TxRpcUrl {
		RpcMap[v] = &rpcStatus{
			Url:       v,
			CoolDown:  false,
			FailCount: 0,
		}
	}
	GRpc = GDatabase.TxRpcUrl[0]
	return nil
}

func updateRpc() {
	RpcMap[GRpc].CoolDown = true
	RpcMap[GRpc].FailCount++

	minUsed := RpcMap[GDatabase.TxRpcUrl[0]]

	for _, v := range RpcMap {
		if minUsed.FailCount > v.FailCount {
			minUsed = v
		}
	}

	GRpc = minUsed.Url
	RpcMap[GRpc].CoolDown = false
}
