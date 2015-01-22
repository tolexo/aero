package orm

import (
	"aero/db"
	"github.com/jinzhu/gorm"
)

var engines map[string]gorm.DB

func init() {
	engines = make(map[string]gorm.DB)
}

func Get(writable bool) gorm.DB {
	connStr := db.GetConnString(writable)
	return getOrm(connStr)
}

func GetFromConf(container string) gorm.DB {
	connStr := db.ParseConnStringFromConfig(container)
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
		if CustomInit != nil {
			CustomInit(&ormObj)
		}
		engines[connStr] = ormObj
		return engines[connStr]
	} else {
		panic(err)
	}
}

var CustomInit func(*gorm.DB) = func(o *gorm.DB) {
	a := *o
	a.SingularTable(true)
}
