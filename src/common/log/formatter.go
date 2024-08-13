package log

import (
	"github.com/sirupsen/logrus"
)

type StackFormatter struct {
	logrus.TextFormatter
}

func (f *StackFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	stack, ok := entry.Data["stack"]
	if ok {
			delete(entry.Data, "stack")
	}
	res, err := f.TextFormatter.Format(entry)
	if stack, ok := stack.(string); ok && stack != "" {
			res = append(res, []byte("stack trace:\n"+stack)...)
	}
	return res, err
}