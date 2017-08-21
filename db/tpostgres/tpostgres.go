package tpostgres

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/tolexo/aero/conf"
	"math/rand"
)

var connPostgresWrite string
var connPostgresRead []string
var validSSLMode = map[string]int{"disable": 1, "allow": 1, "prefer": 1, "require": 1, "verify-ca": 1, "verify-full": 1}

func initMaster(sslMode string) {
	lookup := "database.master"
	if conf.Exists(lookup) {
		connPostgresWrite = getPostgresConnString(lookup, sslMode)
	}
}
func initSlaves(sslMode string) {
	lookup := "database.slaves"
	if conf.Exists(lookup) {
		slaves := conf.StringSlice(lookup, []string{})
		connPostgresRead = make([]string, len(slaves))
		for i, container := range slaves {
			connPostgresRead[i] = getPostgresConnString(container, sslMode)
		}
	}
}
func getPostgresConnString(container string, sslmode string) string {
	if !conf.Exists(container) {
		panic("Container for postgres configuration not found")
	}

	username := conf.String(container+".username", "")
	password := conf.String(container+".password", "")
	host := conf.String(container+".host", "")
	port := conf.String(container+".port", "")
	db := conf.String(container+".db", "")
	return fmt.Sprintf("host=%s port=%s sslmode=%s user=%s dbname=%s password=%s", host, port, sslmode, username, db, password)
}

func getDefaultConn(write bool, sslModeVal ...string) string {
	sslMode := "disable"
	if len(sslModeVal) > 0 {
		sslMode = sslModeVal[0]
	}
	if write {
		initMaster(sslMode)
		return connPostgresWrite
	} else {
		initSlaves(sslMode)
		if connPostgresRead == nil || len(connPostgresRead) == 0 {
			initMaster(sslMode)
			return connPostgresWrite
		}
		return connPostgresRead[rand.Intn(len(connPostgresRead))]
	}
}

//Get MySql connection
func GetPostgresConn(writable bool, sslMode ...string) (dbConn gorm.DB, err error) {
	if isValidSSLMode(sslMode...) == true {
		connStr := getDefaultConn(writable, sslMode...)
		fmt.Println(connStr)
		if dbConn, err = gorm.Open("postgres", connStr); err != nil {
			return
		}
	} else {
		err = errors.New("Invalid sslMode given")
	}
	return
}

//check if sslmode is valid or not
func isValidSSLMode(sslMode ...string) (flag bool) {
	found := false
	if len(sslMode) > 0 {
		for _, curSSLMode := range sslMode {
			if _, exists := validSSLMode[curSSLMode]; exists == true {
				found = true
				break
			}
		}
	} else {
		return true
	}

	if found == true {
		flag = true
	}
	return
}
