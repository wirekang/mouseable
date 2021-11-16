package lg

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

var LogPath = strings.ReplaceAll(
	path.Join(
		os.Getenv("USERPROFILE"), "mouseable.log",
	), "\\", "/",
)

var logger *log.Logger

func init() {
	f, err := os.OpenFile(LogPath, os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	logger = log.New(f, "", log.Llongfile|log.LstdFlags)
}

func printf(prefix string, format string, v ...interface{}) {
	logger.SetPrefix(fmt.Sprintf("%-10s", prefix))
	logger.Printf(format, v...)
}

func Logf(format string, v ...interface{}) {
	printf("LOG", format, v...)
}

func Errorf(format string, v ...interface{}) {
	printf("ERROR", format, v...)
}
