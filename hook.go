// Copyright 2019 Google LLC

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     https://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stackrus

import (
	"context"
	"fmt"

	"cloud.google.com/go/logging"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

// Hook provides a Logrus logging hook that sends logs to Stackdriver Logging
type Hook struct {
	client *logging.Client
	logger *logging.Logger
	levels []logrus.Level
}

// NewHook create a new Hook with project and logname configuration
func NewHook(ctx context.Context, projectID, logName string, level logrus.Level, opts ...option.ClientOption) *Hook {
	hook := &Hook{}

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
		fmt.Println("Stackrus hook ERR: ", e)
	}

	return hook
}

// Fire logs entry to Stackdriver
func (l *Hook) Fire(e *logrus.Entry) error {
	le := logrusToStackdriverEntry(e)
	l.logger.Log(le)
	return nil
}

// Flush clears any pending logs to the network
func (l *Hook) Flush() {
	l.logger.Flush()
}

// Close flushes any log entries before closing the hook
func (l *Hook) Close() {
	l.logger.Flush()
	l.client.Close()
}

// Levels returns the available logging levels
func (l *Hook) Levels() []logrus.Level {
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
