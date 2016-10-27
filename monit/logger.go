package monit

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/tolexo/aero/conf"
)

var (
	logFp      *os.File
	fileErr    error
	logger     *log.Logger
	prevDay    int
	currentDay int
	isLog      bool
	syncOnce   sync.Once
)

func PanicLogger(panicMsg interface{}) {

	syncOnce.Do(func() {
		isLog = conf.Bool("monitor.panic_log", false)
	})

	if isLog == true {
		sTime := time.Now().UTC()
		currentDay = sTime.Day()
		if currentDay != prevDay {
			if prevDay != 0 && logFp != nil {
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
	} else {
		panic(panicMsg)
	}
}
