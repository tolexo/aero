package activity

import (
	"time"

	"github.com/tolexo/aero/activity/model"
	"github.com/tolexo/aero/db/tmongo"
	mgo "gopkg.in/mgo.v2"
)

const (
	DB_CONTAINER = "database.omni"
)

//Log User activity
func LogActivity(url string, body interface{},
	resp interface{}, respCode int, respTime float64) {
	apiDetail := model.APIDetail{
		Url:      url,
		Body:     body,
		Resp:     resp,
		RespCode: respCode,
		RespTime: respTime,
		Time:     time.Now(),
	}
	if sess, mdb, err := tmongo.GetMongoConn(DB_CONTAINER); err == nil {
		defer sess.Close()
		sess.SetSafe(&mgo.Safe{W: 0})
		sess.DB(mdb).C("activity").Insert(apiDetail)
	}
}
