package postgres

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-pg/pg"
	"github.com/tolexo/aero/conf"
)

var (
	isDebuggerActive bool
	StartLogging     bool
	dbPostgresWrite  *pg.DB
	dbPostgresRead   []*pg.DB
	MasterContainer  = "database.master"
	SlaveContainer   = "database.slaves"
)

const (
	GO_PG_PACKAGE = "/github.com/go-pg"
)

//init master connection
func initMaster() (err error) {
	if conf.Exists(MasterContainer) {
		if dbPostgresWrite != nil {
			isDebuggerActive = true
			return
		} else {
			isDebuggerActive = false
			var postgresWriteOption pg.Options
			if postgresWriteOption, err = getPostgresOptions(MasterContainer); err == nil {
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
	if conf.Exists(SlaveContainer) {
		slaves := conf.StringSlice(SlaveContainer, []string{})
		if dbPostgresRead == nil {
			isDebuggerActive = false
			dbPostgresRead = make([]*pg.DB, len(slaves))
		} else {
			isDebuggerActive = true
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
	if conf.Exists(MasterContainer) {
		var postgresWriteOption pg.Options
		if postgresWriteOption, err = getPostgresOptions(MasterContainer); err == nil {
			dbPostgresWrite = pg.Connect(&postgresWriteOption)
		}
	} else {
		err = errors.New("Master config does not exists")
	}
	return
}

//create new slave connections
func CreateSlave() (err error) {
	if conf.Exists(SlaveContainer) {
		slaves := conf.StringSlice(SlaveContainer, []string{})
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
		dbConn = dbPostgresWrite
	} else {
		err = initSlaves()
		if len(dbPostgresRead) == 0 {
			err = initMaster()
			dbConn = dbPostgresWrite
		} else {
			dbConn = dbPostgresRead[rand.Intn(len(dbPostgresRead))]
		}
	}
	logQuery(dbConn)
	return
}

//Get postgres connection by container
func ConnByContainer(container string) (*pg.DB, error) {
	if strings.HasSuffix(container, "master") == true {
		MasterContainer = container
		return Conn(true)
	} else if strings.HasSuffix(container, "slaves") == true {
		SlaveContainer = container
		return Conn(false)
	}
	return nil, errors.New("No master or slaves container found in: " + container)
}

//logQuery : Print postgresql query on terminal
func logQuery(conn *pg.DB) {
	if isDebuggerActive == false {
		isDebuggerActive = true
		conn.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
			if StartLogging == true {
				if query, err := event.FormattedQuery(); err == nil {
					var queryError string
					if event.Error != nil {
						queryError = "\nQUERY ERROR: " + event.Error.Error()
					}
					fmt.Println("----DEBUGGER----")
					fmt.Printf("\nFile: %v : %v\nFunction: %v\nQuery Execution Taken: %s\n%s%s\n\n",
						event.File, event.Line, event.Func, time.Since(event.StartTime), query, queryError)
				} else {
					fmt.Println("Debugger Error: " + err.Error())
				}
			}
		})
	}
}
