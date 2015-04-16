package orm

import (
	"github.com/jinzhu/gorm"
	"github.com/thejackrabbit/aero/db"
)

var engines map[string]gorm.DB

func init() {
	engines = make(map[string]gorm.DB)
}

func Get(writable bool) gorm.DB {
	connStr := db.GetDefaultMySqlConn(writable)
	return getOrm(connStr)
}

func GetFromConf(container string) gorm.DB {
	connStr := db.GetMySqlConnFromConfig(container)
	return getOrm(connStr)
}

func getOrm(connStr string) gorm.DB {
	var ormObj gorm.DB
	var ok bool
	var err error

	if ormObj, ok = engines[connStr]; ok {
		return ormObj
	}
	if ormObj, err = gorm.Open("mysql", connStr); err == nil {
		if OrmInit != nil {
			OrmInit(&ormObj)
		}
		engines[connStr] = ormObj
		return engines[connStr]
	} else {
		panic(err)
	}
}

var OrmInit func(*gorm.DB)
