package activity

import (
	"fmt"
	"github.com/tolexo/aero/activity/model"
	"github.com/tolexo/aero/db/tmongo"
	"gopkg.in/mgo.v2"
	"time"
)

const (
	DB_CONTAINER = "database.activity"
)

//Log User activity
func LogActivity(url string, body interface{},
	resp interface{}, respCode int, respTime float64) {
	sTime := time.Now()
	apiDetail := model.APIDetail{
		Url:      url,
		Body:     body,
		Resp:     resp,
		RespCode: respCode,
		RespTime: respTime,
		Time:     sTime,
	}
	if sess, mdb, err := tmongo.GetMongoConn(DB_CONTAINER); err == nil {
		defer sess.Close()
		mdb = fmt.Sprintf("%v_%v_%v_%v", mdb, sTime.Day(), sTime.Month(), sTime.Year())
		sess.SetSafe(&mgo.Safe{W: 0})
		sess.DB(mdb).C("activity").Insert(apiDetail)
	}
}
