package logruspretty

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"io"
	"strings"
)

type PrettyHandler struct {
	l     *logrus.Logger
	attrs logrus.Fields
}

// NewPrettyHandler создает новый экземпляр PrettyHandler с указанным выводом
func NewPrettyHandler(out io.Writer) *PrettyHandler {
	logger := logrus.New()
	logger.SetOutput(out)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})

	return &PrettyHandler{
		l: logger,
	}
}

// Format форматирует запись журнала в красивом стиле
func (h *PrettyHandler) Format(entry *logrus.Entry) ([]byte, error) {
	level := strings.ToUpper(entry.Level.String()) + ":"
	msg := color.CyanString(entry.Message)

	// Преобразование уровня лога в цвет
	switch entry.Level {
	case logrus.DebugLevel:
		level = color.MagentaString(level)
	case logrus.InfoLevel:
		level = color.BlueString(level)

	case logrus.WarnLevel:
		level = color.YellowString(level)
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		level = color.RedString(level)
	}

	// Отображение времени
	timeStr := entry.Time.Format("15:04:05.000")

	// Преобразование атрибутов logrus в logrus.Fields
	fields := entry.Data

	// Преобразование атрибутов logrus в JSON строку
	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return nil, err
	}

	formattedMsg := fmt.Sprintf("[%s] %s %s %s\n", timeStr, level, msg, color.WhiteString(string(b)))

	return []byte(formattedMsg), nil
}

// WithAttrs добавляет атрибуты к PrettyHandler
func (h *PrettyHandler) WithAttrs(attrs logrus.Fields) *PrettyHandler {
	h.attrs = attrs
	return h
}
