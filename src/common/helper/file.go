package helper

import (
	"fmt"
	"github.com/dwprz/prasorganic-user-service/src/common/log"
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
	"time"
)

func CreateUnixFileName(filename string) string {
	re := regexp.MustCompile(`[ %?#&=]`)
	encodedName := re.ReplaceAllString(filename, "-")
	epochMillis := time.Now().UnixMilli()

	filename = fmt.Sprintf("%d-%s", epochMillis, encodedName)
	return filename
}

func CheckExistDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		return err
	}

	return nil
}

func DeleteFile(path string) {
	if _, err := os.Stat(path); err == nil {
		if err := os.Remove(path); err != nil {
			log.Logger.WithFields(logrus.Fields{"location": "helper.DeleteFile", "section": "os.Remove"}).Error(err)
		}
	}
}
