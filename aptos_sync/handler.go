package main

import (
	"strconv"
	"strings"

	"github.com/Port3-Network/AptosParser/models"
	oo "github.com/Port3-Network/liboo"
)

func GetSyncBlockNum() (bcnum int64, err error) {
	bcnum, err = GetSysConfigInt64(SyncLatestBlock)
	if nil != err {
		err = oo.NewError("GetSysConfigInt64 err %v", err)
		return
	}

	return
}

func GetSysConfigInt64(key string) (num int64, err error) {
	str, err := GetSysConfig(key)
	if nil != err {
		return
	}

	num, err = strconv.ParseInt(str, 10, 64)
	if nil != err {
		err = oo.NewError("failed to ParseInt[%s] err[%v]", str, err)
		return
	}

	return
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

func ParseType(e string) *Contract {
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

func ParseResource(r string) *ResourceStruct {
	res := &ResourceStruct{}
	rArr := strings.Split(r, "::")
	if len(rArr) != 3 {
		return nil
	}
	res.Address = rArr[0]
	res.Module = rArr[1]
	res.Name = rArr[2]
	return res
}

type Contract struct {
	Type     string
	Address  string
	Module   string
	Name     string
	Resource string
}

type ResourceStruct struct {
	Address string
	Module  string
	Name    string
}
