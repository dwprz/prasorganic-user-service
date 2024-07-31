package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
	"github.com/sirupsen/logrus"
)

type customJSONFormatter struct{}

func (f *customJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields, len(entry.Data)+3)
	for k, v := range entry.Data {
		data[k] = v
	}
	data["level"] = entry.Level.String()
	data["msg"] = entry.Message
	data["time"] = entry.Time.Format(time.RFC3339)

	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %w", err)
	}
	return b.Bytes(), nil
}
