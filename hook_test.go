package stackrus

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/logging"

	"github.com/sirupsen/logrus"
)

var levelTests = []struct {
	logrusLevel      logrus.Level
	stackdriverLevel logging.Severity
}{
	{logrus.DebugLevel, logging.Debug},
	{logrus.InfoLevel, logging.Info},
	{logrus.WarnLevel, logging.Warning},
	{logrus.ErrorLevel, logging.Error},
	{logrus.FatalLevel, logging.Critical},
	{logrus.PanicLevel, logging.Emergency},
}

func TestEntryConversion(t *testing.T) {
	testtime := time.Now()
	logger := logrus.New()
	logger.Out = &bytes.Buffer{}
	entry := logrus.NewEntry(logger)
	entry = entry.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	})
	entry.Time = testtime
	entry.Message = "test payload"
	for _, ltest := range levelTests {
		entry.Level = ltest.logrusLevel
		result := logrusToStackdriverEntry(entry)
		if result.Severity != ltest.stackdriverLevel {
			t.Errorf("incorrectlevel, expected %v, got %v", ltest.logrusLevel, ltest.stackdriverLevel)
		}

		if result.Timestamp != testtime {
			t.Errorf("timestamp mismatch")
		}

		if result.Payload != "test payload" {
			t.Errorf("payload mismatch")
		}

		if result.Labels["animal"] != "walrus" {
			t.Errorf("string label mismatch")
		}

		if result.Labels["size"] != "10" {
			t.Errorf("int conversion label mismatch")
		}

	}

}

// this is a live integration test and will ping the stackdriver
// client in the hook, will only run if it finds a STACKRUS_TEST_PROJECT
// environment variable set. If you are @google - you can use
// STACKRUS_TEST_PROJECT='stackrus-test-project' go test
func TestPingStackdriver(t *testing.T) {
	testProject := os.Getenv("STACKRUS_TEST_PROJECT")
	ctx := context.Background()
	if testProject == "" {
		t.Skip("skipping live integration test")
	}
	testhook := NewStackdriverHook(ctx, testProject, "dummy", logrus.DebugLevel)
	err := testhook.client.Ping(ctx)
	if err != nil {
		t.Errorf("Failed to ping stackdriver client: %s", err.Error())
	}
}
