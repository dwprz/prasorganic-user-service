package helper

import (
	"github.com/sirupsen/logrus"
)

func HandlePanic(name string, l *logrus.Logger) {
	message := recover()

	if message != nil {
		l.Errorf("%v | %v", name, message)
	}

}
