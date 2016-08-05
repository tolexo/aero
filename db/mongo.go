package db

import (
	"errors"
	"fmt"
	"github.com/tolexo/aero/conf"
	"gopkg.in/mgo.v2"
)

// create the mongo connection string
func getMongoConnStr(container string) (string, string) {
	username := conf.String(container+".username", "")
	password := conf.String(container+".password", "")
	host := conf.String(container+".host", "")
	port := conf.String(container+".port", "")
	db := conf.String(container+".db", "")
	replicas := conf.String(container+".replicas", "")
	options := conf.String(container+".options", "")

	if port != "" {
		port = ":" + port
	}
	if replicas != "" {
		replicas = "," + replicas
	}
	if options != "" {
		options = "?" + options
	}
	auth := ""
	if username != "" || password != "" {
		auth = username + ":" + password + "@"
	}

	return fmt.Sprintf("mongodb://%s%s%s%s/%s%s", auth, host, port, replicas, db, options), db
}

// validate the container string
func validateContainer(container ...string) (db string, err error) {
	cLen := len(container)
	if cLen == 0 {
		container = append(container, "database.mongo")
	} else if cLen > 1 {
		err = errors.New("More than one container not supported")
		return
	}
	if conf.Exists(container[0]) == false {
		err = errors.New("mongo configuration missing")
		return
	}
	db = container[0]
	return
}

// create mongo connection
// TODO introduce parameter for additional connection settings like socket timeout
func GetMongoConn(container ...string) (sess *mgo.Session, mdb string, err error) {
	var db string
	if db, err = validateContainer(container...); err == nil {
		conn, mdb := getMongoConnStr(db)
		if mdb == "" {
			err = errors.New("mongo database name missing")
		}
		sess, err = mgo.Dial(conn)
	}
	return
}
