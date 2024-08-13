package test

import (
	"os"

	"github.com/dwprz/prasorganic-user-service/src/common/log"
	"github.com/sirupsen/logrus"
)

// ini untuk merubah working directory path saat menjalankan test supaya path nya berawal dari root

func init() {
	err := os.Chdir(os.Getenv("PRASORGANIC_USER_SERVICE_WORKSPACE"))
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "test.init", "section": "os.Chdir"}).Fatal(err)
	}
}
