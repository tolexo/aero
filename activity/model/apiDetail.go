package model

import "time"

//Complete API detail
type APIDetail struct {
	Url       string      `bson:"url"`
	ServiceID string      `bson:"service"`
	Body      interface{} `bson:"body"`
	Resp      interface{} `bson:"response"`
	RespCode  int         `bson:"response_code"`
	RespTime  float64     `bson:"response_time"`
	Time      time.Time   `bson:"time"`
}
