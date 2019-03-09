package main

import (
	"context"

	"github.com/ptone/stackrus"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()
var hook *stackrus.Hook

func main() {
	ctx := context.Background()
	hook = stackrus.NewHook(ctx, "stackrus-test-project", "logrus-test-log", logrus.InfoLevel)
	log.Hooks.Add(hook)
	log.Formatter = new(logrus.TextFormatter)
	log.Level = logrus.DebugLevel
	defer hook.Close()
	defer func() {
		err := recover()
		if err != nil {
			log.WithFields(logrus.Fields{
				"omg":    true,
				"err":    err,
				"number": 100,
			}).Fatal("The ice breaks!")
		}
	}()

	log.WithFields(logrus.Fields{
		"animal": "walrus",
		"number": 8,
	}).Debug("Started observing beach")

	log.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	log.WithFields(logrus.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")

	// you won't see this on in the logs - as the level is set to Info or higher
	log.WithFields(logrus.Fields{
		"temperature": -4,
	}).Debug("Temperature changes")
}
