package monit

import (
	"fmt"
	"github.com/tolexo/aero/conf"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	fmt.Fprintf(w, "%s", "4041 page not found")

	dataDogAgent := GetDataDogAgent()
	var statusCode int64 = 4041
	dataDogAgent.Count(FormatHttpStatusCode(statusCode), 1, nil, 1)
}

func FormatHttpStatusCode(httpStatusCode int64) string {
	return "http_" + strconv.FormatInt(httpStatusCode, 10)
}

func GetTimeIntervalTag(respTime int64) (ret []string) {
	interval := conf.String("monitor.interval", "")
	if interval != "" {
		intervalArr := strings.Split(interval, ",")
		var lowerLimit, higherLimit int64
		var matched bool
		for i := range intervalArr {
			higherLimit, _ = strconv.ParseInt(intervalArr[i], 10, 64)
			if respTime <= higherLimit {
				matched = true
				break
			} else {
				lowerLimit = higherLimit
			}
		}
		var intervalGroup string
		if matched == true {
			intervalGroup = "from_" + strconv.FormatInt(lowerLimit, 10) + "_to_" + strconv.FormatInt(higherLimit, 10)
		} else {
			intervalGroup = "above_" + strconv.FormatInt(higherLimit, 10)
		}
		ret = append(ret, intervalGroup)
	}
	return
}

func ModRecorder() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := httptest.NewRecorder()
			next.ServeHTTP(rec, r)

			for k, v := range rec.Header() {
				w.Header()[k] = v
			}

			len := rec.Body.Len()
			w.Header().Set("Content-Length", strconv.Itoa(len))
			w.WriteHeader(rec.Code)
			w.Write(rec.Body.Bytes())
			dur := time.Since(start).Seconds() * 1000
			fmt.Println("Time taken by request, Header code, time interval: ", dur, rec.Code, GetTimeIntervalTag(int64(dur)))

			dataDogAgent := GetDataDogAgent()
			dataDogAgent.Histogram("APICall", dur, GetTimeIntervalTag(int64(dur)), 1)
			dataDogAgent.Count(FormatHttpStatusCode(int64(rec.Code)), 1, nil, 1)

		})
	}
}
