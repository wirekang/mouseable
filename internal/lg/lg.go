package lg

import (
	"fmt"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

func printf(prefix string, format string, v ...interface{}) {
	logger.SetPrefix(fmt.Sprintf("%-7s", prefix))
	v = append([]interface{}{prefix}, v...)
	logger.Printf("%-7s"+format, v...)
}

func Printf(format string, v ...interface{}) {
	printf("LOG", format, v...)
}

func Errorf(format string, v ...interface{}) {
	printf("ERROR", format, v...)
}
