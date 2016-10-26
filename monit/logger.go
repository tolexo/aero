package monit

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tolexo/aero/conf"
)

func PanicLogger(panicMsg interface{}) {
	var (
		logFp   *os.File
		fileErr error
		logger  *log.Logger
	)
	sTime := time.Now().UTC()
	path := conf.String("logs.panic_log", "panic_log")
	path = fmt.Sprintf("%s_%d-%d-%d.log", path, sTime.Day(), sTime.Month(), sTime.Year())
	if logFp, fileErr = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); fileErr == nil {
		logger = log.New(logFp, "panic", log.Lshortfile)
		logger.Panic(panicMsg)
	} else {
		fmt.Println("Could not create the panic log file")
		panic(panicMsg)
	}
}
