# stackrus
[Stackdriver logging](https://cloud.google.com/logging/) plugin for [logrus](https://github.com/sirupsen/logrus)

See the example for usage

Note - to fully flush the hook's buffer, you need to be sure `Close` is called
on the hook before exiting.

Logrus data field values are all converted to string label values in
Stackdriver log entries per the Stackdriver API.

Levels are mapped between Logrus levels and Stackdriver severity:

	{logrus.DebugLevel, logging.Debug},
	{logrus.InfoLevel, logging.Info},
	{logrus.WarnLevel, logging.Warning},
	{logrus.ErrorLevel, logging.Error},
	{logrus.FatalLevel, logging.Critical},
	{logrus.PanicLevel, logging.Emergency},

This hook specifically uses "[cloud.google.com/go/logging](https://godoc.org/cloud.google.com/go/logging)" while many other related projects are only using the deprecated "[google.golang.org/api/logging/v2](https://godoc.org/google.golang.org/api/logging/v2)".

Not an official Google product.