package tmysql

import (
	"fmt"
	"math/rand"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/tolexo/aero/conf"
)

var (
	engines        map[string]bool
	connMySqlWrite string
	connMySqlRead  []string
)

func init() {
	engines = make(map[string]bool)
}

func initMaster() {
	lookup := "database.master"
	if conf.Exists(lookup) {
		connMySqlWrite = getMySqlConnString(lookup)
	}
}
func initSlaves() {
	lookup := "database.slaves"
	if conf.Exists(lookup) {
		slaves := conf.StringSlice(lookup, []string{})
		connMySqlRead = make([]string, len(slaves))
		for i, container := range slaves {
			connMySqlRead[i] = getMySqlConnString(container)
		}
	}
}
func getMySqlConnString(container string) string {
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

func getDefaultConn(write bool) string {
	if write {
		initMaster()
		return connMySqlWrite
	} else {
		initSlaves()
		if connMySqlRead == nil || len(connMySqlRead) == 0 {
			initMaster()
			return connMySqlWrite
		}
		return connMySqlRead[rand.Intn(len(connMySqlRead))]
	}
}

//Get MySql connection
func GetMySqlConn(writable bool) (dbConn gorm.DB, err error) {
	connStr := getDefaultConn(writable)
	if dbConn, err = gorm.Open("mysql", connStr); err == nil && engines[connStr] == false {
		engines[connStr] = true
		dbConn.DB().SetConnMaxLifetime(time.Second * 30)
		dbConn.DB().SetMaxIdleConns(10)
		dbConn.DB().SetMaxOpenConns(200)
	}
	return
}
