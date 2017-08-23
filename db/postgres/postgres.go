package postgres

import (
	"errors"
	"github.com/go-pg/pg"
	"github.com/tolexo/aero/conf"
	"math/rand"
)

var dbPostgresWrite *pg.DB
var dbPostgresRead []*pg.DB

//init master connection
func initMaster() (err error) {
	lookup := "database.master"
	if conf.Exists(lookup) {
		if dbPostgresWrite != nil {
			return
		} else {
			var postgresWriteOption pg.Options
			if postgresWriteOption, err = getPostgresOptions(lookup); err == nil {
				dbPostgresWrite = pg.Connect(&postgresWriteOption)
			}
		}
	} else {
		err = errors.New("Master config does not exists")
	}
	return
}

//init slave connections
func initSlaves() (err error) {
	lookup := "database.slaves"
	if conf.Exists(lookup) {
		slaves := conf.StringSlice(lookup, []string{})
		if dbPostgresRead == nil {
			dbPostgresRead = make([]*pg.DB, len(slaves))
		}
		for i, container := range slaves {
			if dbPostgresRead[i] == nil {
				var postgresReadOption pg.Options
				if postgresReadOption, err = getPostgresOptions(container); err != nil {
					break
				}
				dbPostgresRead[i] = pg.Connect(&postgresReadOption)
			}
		}
	} else {
		err = errors.New("Slaves config does not exists")
	}
	return
}

//create new master connection
func CreateMaster() (err error) {
	lookup := "database.master"
	if conf.Exists(lookup) {
		var postgresWriteOption pg.Options
		if postgresWriteOption, err = getPostgresOptions(lookup); err == nil {
			dbPostgresWrite = pg.Connect(&postgresWriteOption)
		}
	} else {
		err = errors.New("Master config does not exists")
	}
	return
}

//create new slave connections
func CreateSlave() (err error) {
	lookup := "database.slaves"
	if conf.Exists(lookup) {
		slaves := conf.StringSlice(lookup, []string{})
		if dbPostgresRead == nil {
			dbPostgresRead = make([]*pg.DB, len(slaves))
		}
		for i, container := range slaves {
			var postgresReadOption pg.Options
			if postgresReadOption, err = getPostgresOptions(container); err != nil {
				break
			}
			dbPostgresRead[i] = pg.Connect(&postgresReadOption)
		}
	} else {
		err = errors.New("Slaves config does not exists")
	}
	return
}

//set postgres connection options from conf
func getPostgresOptions(container string) (pgOption pg.Options, err error) {
	if !conf.Exists(container) {
		err = errors.New("Container for postgres configuration not found")
		return
	}
	host := conf.String(container+".host", "")
	port := conf.String(container+".port", "")
	addr := ""
	if host != "" && port != "" {
		addr = host + ":" + port
	}
	pgOption.Addr = addr
	pgOption.User = conf.String(container+".username", "")
	pgOption.Password = conf.String(container+".password", "")
	pgOption.Database = conf.String(container+".db", "")
	pgOption.MaxRetries = conf.Int(container+".maxRetries", 3)
	pgOption.RetryStatementTimeout = conf.Bool(container+".retryStmTimeout", false)
	return
}

//Get postgres connection
func Conn(writable bool) (dbConn *pg.DB, err error) {
	if writable {
		err = initMaster()
		return dbPostgresWrite, err
	} else {
		err = initSlaves()
		if dbPostgresRead == nil || len(dbPostgresRead) == 0 {
			err = initMaster()
			return dbPostgresWrite, err
		}
		return dbPostgresRead[rand.Intn(len(dbPostgresRead))], err
	}
}
