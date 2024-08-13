package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logger = logrus.New()

func init() {
	appStatus := os.Getenv("PRASORGANIC_APP_STATUS")

	Logger.SetFormatter(&StackFormatter{
		logrus.TextFormatter{
			DisableColors:    false,
			DisableTimestamp: false,
			FullTimestamp:    true,
			DisableQuote:     true,
		},
	})

	Logger.SetLevel(logrus.InfoLevel)

	if appStatus == "DEVELOPMENT" {
		return
	}

	if _, err := os.Stat("./app.log"); os.IsNotExist(err) {
		if err := os.Mkdir("./tmp", os.ModePerm); err != nil {
			Logger.WithFields(logrus.Fields{"location": "log.init", "section": "os.Mkdir"}).Fatal(err)
		}
	}

	file, err := os.OpenFile("./app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Logger.WithFields(logrus.Fields{"location": "log.init", "section": "helper.CheckExistDir"}).Fatal(err)
	}

	Logger.Out = file
}
