package auth

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/tolexo/aero/cache"
	"github.com/tolexo/aero/conf"
	"github.com/tolexo/aero/db/orm"
	"net/http"
	"strconv"
	"time"
)

const TIMESTAMP_VALIDITY int = 60000
const APP_SECRET string = "TOLEXOANDROIDAPP"

type AuthParams struct {
	AppSession string
	Timestamp  string
	Hash       string
}

type Session struct {
	CustomerId   int
	CustomerType string
	ExpiryDate   time.Time
}

func GetSessionFromHeader(r *http.Request) (int, string) {

	a := getAuthParamsFromHeader(r)
	code := CheckTimestamp(a)
	if code != 200 {
		return code, ""
	}
	code = CheckHash(a)
	if code != 200 {
		return code, ""
	}

	code, session := getSessionInfo(a)
	se, err := json.Marshal(session)
	if err == nil {
		return 200, string(se[:])
	}

	return 604, ""

}

func getAuthParamsFromHeader(r *http.Request) (authParam AuthParams) {

	if len(r.Header["Appsession"]) > 0 {
		authParam.AppSession = r.Header["Appsession"][0]
	} else {
		authParam.AppSession = ""
	}
	if len(r.Header["Timestamp"]) > 0 {
		authParam.Timestamp = r.Header["Timestamp"][0]
	} else {
		authParam.Timestamp = ""
	}
	if len(r.Header["Hash"]) > 0 {
		authParam.Hash = r.Header["Hash"][0]
	} else {
		authParam.Hash = ""
	}

	return authParam
}

func CheckTimestamp(a AuthParams) (code int) {
	today := time.Now()
	i, err := strconv.ParseInt(a.Timestamp, 10, 64)
	if err != nil {
		return 602 //error code for invalid timestamp
	}
	t := time.Unix(i, 0)

	duration := -1
	if err == nil {
		duration = int(today.Sub(t).Minutes())
	}

	if duration < 0 || duration > TIMESTAMP_VALIDITY {
		return 602 //error code for invalid timestamp
	}
	return 200
}

func CheckHash(a AuthParams) (code int) {
	data := a.Timestamp + "|" + a.AppSession + "|" + APP_SECRET
	//hash := common.GetMD5Hash(data)	//TODO move md5 to common place

	h := md5.New()
	h.Write([]byte(data))
	hash := hex.EncodeToString(h.Sum(nil))

	if hash != a.Hash {
		return 603 //error code for invalid hash
	}
	return 200

}

func getSessionInfo(a AuthParams) (code int, s Session) {
	s = Session{}
	s, ok := getSessionDetailsFromCache(a.AppSession)

	if ok {
		return 200, s
	}

	s, ok = getSessionDetailsFromDb(a.AppSession)
	if ok {
		saveSessionInCache(a.AppSession, s)
		return 200, s
	}

	return 604, s
}

func getSessionDetailsFromCache(appSession string) (s Session, ok bool) {
	cType := conf.String("primary.type", "")
	if cType != "redis" {
		panic("redis configuration not found: " + cType)
	}

	c := cache.RedisFromConfig("primary")

	key := "APPSESSION-" + appSession
	val, err := c.Get(key)
	if err == nil {
		err := json.Unmarshal(val, &s)
		if err == nil {
			return s, true
		}
	}
	return s, false
}

func saveSessionInCache(appSession string, s Session) {
	cType := conf.String("primary.type", "")
	if cType != "redis" {
		panic("redis configuration not found: " + cType)
	}
	c := cache.RedisFromConfig("primary")
	key := "APPSESSION-" + appSession
	val, _ := json.Marshal(s)
	ttl := s.ExpiryDate.Sub(time.Now())
	c.Set(key, val, ttl)
}

func getSessionDetailsFromDb(appSession string) (s Session, ok bool) {

	db := orm.Get(false)

	sql := "Select customer_id, customer_type, expiry_date from customer_auth where app_session = ? and status = 'active'"
	db.Raw(sql, appSession).Row().Scan(&s.CustomerId, &s.CustomerType, &s.ExpiryDate)

	return s, true
}
