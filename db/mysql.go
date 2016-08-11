package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tolexo/aero/conf"
	"github.com/tolexo/aero/db/orm"
	"math/rand"
	"net/url"
)

func init() {
	initMySqlMaster()
	initMySqlSlaves()
}

var connMySqlWrite string
var connMySqlRead []string

func initMySqlMaster() {
	lookup := "database.master"
	if conf.Exists(lookup) {
		connMySqlWrite = GetMySqlConnFromConfig(lookup)
	}
}

func initMySqlSlaves() {
	lookup := "database.slaves"
	if conf.Exists(lookup) {
		slaves := conf.StringSlice(lookup, []string{})
		connMySqlRead = make([]string, len(slaves))
		for i, container := range slaves {
			connMySqlRead[i] = GetMySqlConnFromConfig(container)
		}
	}
}
func GetMySqlConnFromConfig(container string) string {
	if !conf.Exists(container) {
		panic("Container for mysql configuration not found")
	}

	username := conf.String(container+".username", "")
	password := conf.String(container+".password", "")
	host := conf.String(container+".host", "")
	port := conf.String(container+".port", "")
	db := conf.String(container+".db", "")
	timezone := conf.String(container+".timezone", "")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=%s",
		username, password,
		host, port, db,
		url.QueryEscape(timezone),
	)
}

func GetDefaultMySqlConn(write bool) string {
	if write {
		return connMySqlWrite
	}

	if connMySqlRead == nil || len(connMySqlRead) == 0 {
		return connMySqlWrite
	}

	return connMySqlRead[rand.Intn(len(connMySqlRead))]
}
func Get(writable bool) gorm.DB {
	connStr := GetDefaultMySqlConn(writable)
	return orm.getOrm(connStr)
}

func GetFromConf(container string) gorm.DB {
	connStr := GetMySqlConnFromConfig(container)
	return orm.getOrm(connStr)
}
