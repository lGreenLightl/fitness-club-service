package logs

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	prefixedFormatter "github.com/x-cray/logrus-prefixed-formatter"
)

func InitLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyTime:  "time",
		},
	})

	if isLocalEnvironment, _ := strconv.ParseBool(os.Getenv("LOCAL_ENV")); isLocalEnvironment {
		logrus.SetFormatter(&prefixedFormatter.TextFormatter{
			ForceFormatting: true,
			FullTimestamp:   true,
		})
	}

	logrus.SetLevel(logrus.DebugLevel)
}
