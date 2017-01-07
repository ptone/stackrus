package stackrus

import (
	"context"
	"fmt"

	"cloud.google.com/go/logging"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

type StackdriverHook struct {
	client *logging.Client
	logger *logging.Logger
	levels []logrus.Level
}

func NewStackdriverHook(ctx context.Context, projectID, logName string, level logrus.Level, opts ...option.ClientOption) *StackdriverHook {
	hook := &StackdriverHook{}

	c, err := logging.NewClient(ctx, projectID)
	if err != nil {
		panic("unable to create logging client")
	}
	hook.client = c
	// TODO consider common map with logrus tag
	hook.logger = c.Logger(logName)

	logLevels := []logrus.Level{}
	for _, l := range []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	} {
		if l <= level {
			logLevels = append(logLevels, l)
		}
	}
	hook.levels = logLevels
	hook.client.OnError = func(e error) {
		fmt.Println("Stackrus hook ERR: ", e.Error())
	}

	return hook
}

func (l *StackdriverHook) Fire(e *logrus.Entry) error {
	le := logrusToStackdriverEntry(e)
	l.logger.Log(le)
	return nil
}

func (l *StackdriverHook) Flush() {
	l.logger.Flush()
}

func (l *StackdriverHook) Close() {
	l.logger.Flush()
	l.client.Close()
}

func (l *StackdriverHook) Levels() []logrus.Level {
	return l.levels
}

func logrusToStackdriverEntry(ge *logrus.Entry) logging.Entry {
	le := logging.Entry{}
	le.Timestamp = ge.Time
	le.Payload = ge.Message
	labels := make(map[string]string, 0)
	for k, v := range ge.Data {
		labels[k] = fmt.Sprintf("%v", v)
	}
	le.Labels = labels
	// TODO severity
	switch ge.Level {
	case logrus.DebugLevel:
		le.Severity = logging.Debug
	case logrus.InfoLevel:
		le.Severity = logging.Info
	case logrus.WarnLevel:
		le.Severity = logging.Warning
	case logrus.ErrorLevel:
		le.Severity = logging.Error
	case logrus.FatalLevel:
		le.Severity = logging.Critical
	case logrus.PanicLevel:
		le.Severity = logging.Emergency
	}

	return le
}
