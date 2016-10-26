package monit

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/tolexo/aero/conf"
)

var (
	logFp      *os.File
	fileErr    error
	logger     *log.Logger
	prevDay    int
	currentDay int
)

func PanicLogger(panicMsg interface{}) {

	sTime := time.Now().UTC()
	currentDay = sTime.Day()
	if currentDay != prevDay {
		if prevDay != 0 {
			logFp.Close()
		}
		path := conf.String("logs.panic_log", "panic_log")
		path = fmt.Sprintf("%s_%d-%d-%d.log", path, sTime.Day(), sTime.Month(), sTime.Year())
		if logFp, fileErr = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); fileErr != nil {
			fmt.Println("Could not create the panic log file")
			panic(panicMsg)
		}
		prevDay = currentDay
	}
	logger = log.New(logFp, "panic", log.Lshortfile)
	logger.Print(string(debug.Stack()))
	logger.Panic(panicMsg)
}
