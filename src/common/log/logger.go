package log

import (
	"os"

	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
)

var Logger = logrus.New()

func init() {
	appStatus := os.Getenv("PRASORGANIC_APP_STATUS")

	Logger.SetFormatter(&ecslogrus.Formatter{
		PrettyPrint: true,
	})

	Logger.SetLevel(logrus.InfoLevel)

	if appStatus == "DEVELOPMENT" {
		return
	}

	Logger.SetFormatter(&ecslogrus.Formatter{
		PrettyPrint: false,
	})

	file, err := os.OpenFile("./app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Logger.WithFields(logrus.Fields{"location": "log.init", "section": "os.OpenFile"}).Fatal(err)
	}

	Logger.Out = file
}
