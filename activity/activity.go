package activity

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tolexo/aero/activity/model"
	"github.com/tolexo/aero/conf"
	"github.com/tolexo/aero/db/tmongo"
	"github.com/tolexo/aero/panik"
	"gopkg.in/mgo.v2"
)

const (
	DB_CONTAINER = "database.activity"
)

//Log User activity
func LogActivity(url string, serviceID string, body interface{},
	resp interface{}, respCode int, respTime float64) {
	sTime := time.Now()
	apiDetail := model.APIDetail{
		Url:       url,
		ServiceID: serviceID,
		Body:      body,
		Resp:      resp,
		RespCode:  respCode,
		RespTime:  respTime,
		Time:      sTime,
	}
	if sess, mdb, err := tmongo.GetMongoConn(DB_CONTAINER); err == nil {
		defer sess.Close()
		mdb = fmt.Sprintf("%v_%v_%v_%v", mdb, sTime.Day(), sTime.Month(), sTime.Year())
		sess.SetSafe(&mgo.Safe{W: 0})
		sess.DB(mdb).C("activity").Insert(apiDetail)
	}
}

//LogCSV : will log in csv
func LogCSV(serviceID, url string, respTime float64, respCode int64) {
	go func() {
		logfile := conf.String("monitor.filelog", "")
		if logfile != "" {
			t := time.Now()
			logfileSuffix := fmt.Sprintf("%d-%d-%d", t.Month(), t.Day(), t.Hour())
			f, err := os.OpenFile(logfile+logfileSuffix+".csv", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			panik.On(err)
			l := log.New(f, "", log.LstdFlags)
			l.Printf("%s %s %f %d", serviceID, url, respTime, respCode)
		}
	}()
}
