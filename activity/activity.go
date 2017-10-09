package activity

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/tolexo/aero/activity/model"
	"github.com/tolexo/aero/db/tmongo"
	mgo "gopkg.in/mgo.v2"
)

const (
	MGO_CONN_REFRESH_COUNTER = 3
	DB_CONTAINER             = "database.omni"
)

var (
	once sync.Once
	sess *mgo.Session
	mdb  string
)

//Log User activity
func LogActivity(url string, body interface{},
	resp reflect.Value, respCode int, respTime float64) {
	var (
		err error
	)
	apiDetail := model.APIDetail{
		Url:      url,
		Body:     body,
		Resp:     resp.Interface(),
		RespCode: respCode,
		RespTime: respTime,
		Time:     time.Now(),
	}
	once.Do(func() {
		if sess, mdb, err = tmongo.GetMongoConn(DB_CONTAINER); err != nil {
			fmt.Println("LogActivity: mongo connection error", err.Error())
		}
	})
	go func() {
		if err = sess.DB(mdb).C("activity").Insert(apiDetail); err != nil {
			for counter := 0; counter < MGO_CONN_REFRESH_COUNTER; counter++ {
				fmt.Println("LogActivity: Refreshing mongo connection", counter)
				sess.Refresh()
				err = sess.DB(mdb).C("activity").Insert(apiDetail)
				if err == nil {
					break
				}
			}
		}
	}()
}
