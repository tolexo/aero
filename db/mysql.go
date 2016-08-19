package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/tolexo/aero/conf"
	"math/rand"
	"net/url"
)

var engines map[string]gorm.DB
var ormInit []func(*gorm.DB)

func init() {
	engines = make(map[string]gorm.DB)
	initMySqlConn()
}
func initMySqlConn() {
	initMaster()
	initSlaves()
}
func initMaster() {
	lookup := "database.master"
	if conf.Exists(lookup) {
		connMySqlWrite = GetMySqlConnString(lookup)
	}
}
func initSlaves() {
	lookup := "database.slaves"
	if conf.Exists(lookup) {
		slaves := conf.StringSlice(lookup, []string{})
		connMySqlRead = make([]string, len(slaves))
		for i, container := range slaves {
			connMySqlRead[i] = GetMySqlConnString(container)
		}
	}
}
func GetMySqlConnString(container string) string {
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

func GetDefaultConn(write bool) string {
	if write {
		return connMySqlWrite
	}

	if connMySqlRead == nil || len(connMySqlRead) == 0 {
		return connMySqlWrite
	}

	return connMySqlRead[rand.Intn(len(connMySqlRead))]
}

func GetMySqlConn(writable bool) (gorm.DB, error) {
	connStr := GetDefaultConn(writable)
	return GetOrm(connStr)
}

func GetOrm(connStr string) (ormObj gorm.DB, err error) {
	var ok bool

	if ormObj, ok = engines[connStr]; ok {
		return
	}
	// http://go-database-sql.org/accessing.html
	// the sql.DB object is designed to be long-lived
	if ormObj, err = gorm.Open("mysql", connStr); err == nil {
		if ormInit != nil {
			for _, fn := range ormInit {
				fn(&ormObj)
			}
		}
		engines[connStr] = ormObj
		return

	} else {
		return
	}
}

func DoOrmInit(fn func(*gorm.DB)) {
	// TODO: use mutex
	if ormInit == nil {
		ormInit = make([]func(*gorm.DB), 0)
	}
	ormInit = append(ormInit, fn)
}
