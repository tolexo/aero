package tmysql

import (
	"fmt"
	"math/rand"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/tolexo/aero/conf"
	"github.com/tolexo/aero/db/orm"
)

var (
	engines        map[string]gorm.DB
	connContainer  map[string]string
	connMySqlWrite string
	connMySqlRead  []string
)

func init() {
	engines = make(map[string]gorm.DB)
	connContainer = make(map[string]string)
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
func getMySqlConnString(container string) (connStr string) {
	if !conf.Exists(container) {
		panic("Container for mysql configuration not found")
	}
	username := conf.String(container+".username", "")
	password := conf.String(container+".password", "")
	host := conf.String(container+".host", "")
	port := conf.String(container+".port", "")
	db := conf.String(container+".db", "")
	timezone := conf.String(container+".timezone", "")

	connStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=%s",
		username, password,
		host, port, db,
		url.QueryEscape(timezone),
	)
	connContainer[connStr] = container
	return
}

func getDefaultConn(write bool) string {
	if write {
		initMaster()
		return connMySqlWrite
	} else {
		initSlaves()
		if len(connMySqlRead) == 0 {
			initMaster()
			return connMySqlWrite
		}
		return connMySqlRead[rand.Intn(len(connMySqlRead))]
	}
}

//newConn will create new database connection and initialize all database setting
func newConn(connStr string) (dbConn gorm.DB, err error) {
	if dbConn, err = gorm.Open("mysql", connStr); err == nil {
		engines[connStr] = dbConn
		container := connContainer[connStr]
		connMaxLifetime := conf.Int(container+".maxlifetime", 10)
		maxIdleConns := conf.Int(container+".maxidleconn", 10)
		maxOpenConns := conf.Int(container+".maxopenconn", 200)
		dbConn.DB().SetConnMaxLifetime(time.Second * time.Duration(connMaxLifetime))
		dbConn.DB().SetMaxIdleConns(maxIdleConns)
		dbConn.DB().SetMaxOpenConns(maxOpenConns)
	}
	return
}

//Get MySql connection
func GetMySqlTmpConn(writable bool) (dbConn gorm.DB, err error) {
	dbConn = orm.Get(writable)
	return
}

//Get MySql connection
func GetMySqlConn(writable bool) (dbConn gorm.DB, err error) {
	var connExists bool
	rand.Seed(time.Now().UnixNano())
	connStr := getDefaultConn(writable)
	if dbConn, connExists = engines[connStr]; connExists == true {
		err = dbConn.DB().Ping()
	}
	if connExists == false || err != nil {
		dbConn, err = newConn(connStr)
	}
	return
}
