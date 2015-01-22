package db

import (
	"aero/conf"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"net/url"
)

var writeConn string
var readConn []string

func init() {
	initializeMasterConnection()
	initializeSlaveConnections()
}

func initializeMasterConnection() {
	lookup := "database.master"
	if conf.Exists(lookup) {
		writeConn = ParseConnStringFromConfig(lookup)
	}
}

func initializeSlaveConnections() {
	lookup := "database.slaves"
	if conf.Exists(lookup) {
		slaves := conf.StringSlice(lookup, []string{})
		readConn = make([]string, len(slaves))
		for i, container := range slaves {
			readConn[i] = ParseConnStringFromConfig(container)
		}
	}
}

// Expected attributes for MySQL connection
// container:
// 	 username: user
//   password: pass
//   host:     localhost
//   port:     3306
//   db:       db-name
//   timezone:
func ParseConnStringFromConfig(container string) string {
	if !conf.Exists(container) {
		panic("Container not found")
	}

	username := conf.Get(container+".username", "")
	password := conf.Get(container+".password", "")
	host := conf.Get(container+".host", "")
	port := conf.Get(container+".port", "")
	db := conf.Get(container+".db", "")
	timezone := conf.Get(container+".timezone", "")

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=%s",
		username, password,
		host, port, db,
		url.QueryEscape(timezone.(string)),
	)
}

func GetConnString(write bool) string {
	if write {
		return writeConn
	}

	if readConn == nil || len(readConn) == 0 {
		return writeConn
	}

	return readConn[rand.Intn(len(readConn))]
}
