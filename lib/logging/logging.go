package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func init() {
	// Log as JSON instead of the default ASCII formatter.
	//Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.Out = os.Stdout

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	Logger.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	Logger.SetLevel(logrus.InfoLevel)
}

//   contextLogger := log.WithFields(log.Fields{
//     "common": "this is a common field",
//     "other": "I also should be logged always",
//   })
