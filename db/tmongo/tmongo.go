package tmongo

import (
	"errors"
	"fmt"
	"github.com/tolexo/aero/conf"
	"gopkg.in/mgo.v2"
	"reflect"
	"time"
)

const (
	MGO_TIMEOUT = "timeout"
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
func validateContainer(container string) (db string, err error) {

	if conf.Exists(container) == false {
		err = errors.New("mongo configuration missing")
		return
	}
	db = container
	return
}

// create mongo connection
// TODO introduce parameter for connection additional settings like socket timeout
func GetMongoConn(container string, param ...map[string]interface{}) (sess *mgo.Session, mdb string, err error) {
	var db string
	if db, err = validateContainer(container); err == nil {
		var conn string
		conn, mdb = getMongoConnStr(db)
		if mdb == "" {
			err = errors.New("mongo database name missing")
			return
		}
		pLen := len(param)
		if pLen > 1 {
			err = errors.New("Parameters other than the Map not supported")
		} else if pLen == 1 {
			if tValue, tExist := param[0][MGO_TIMEOUT]; tExist == true {
				if reflect.TypeOf(tValue).Kind() == reflect.Int64 {
					sess, err = mgo.DialWithTimeout(conn, param[0][MGO_TIMEOUT].(time.Duration))
				} else {
					err = errors.New("Invalid Timeout Value")
				}
			} else {
				sess, err = mgo.Dial(conn)
			}
		} else {
			sess, err = mgo.Dial(conn)
		}
		if err == nil {
			sess.SetMode(mgo.Monotonic, true)
		}
	}
	return
}
