package lg

import (
	"fmt"
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "", log.LstdFlags)
}

func printf(prefix string, format string, v ...interface{}) {
	logger.SetPrefix(fmt.Sprintf("%-7s", prefix))
	logger.Printf(format, v...)
}

func Logf(format string, v ...interface{}) {
	printf("LOG", format, v...)
}

func Errorf(format string, v ...interface{}) {
	printf("ERROR", format, v...)
}
