package helper

import (
	"os"

	"github.com/sirupsen/logrus"
)

func (h *HelperImpl) DeleteFile(path string) {
	if _, err := os.Stat(path); err == nil {
		if err := os.Remove(path); err != nil {
			h.logger.WithFields(logrus.Fields{"location": "helper.HelperImpl/DeleteFile", "section": "os.Remove"}).Error(err)
		}
	}
}
