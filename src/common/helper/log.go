package helper

import (
	"encoding/json"

	"github.com/dwprz/prasorganic-user-service/src/common/log"
	"github.com/sirupsen/logrus"
)

func LogJSON(value any) {
	jsonData, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "helper.LogJSON", "section": "json.MarshalIndent"}).Error(err)
	}

	log.Logger.Info(string(jsonData))
}
