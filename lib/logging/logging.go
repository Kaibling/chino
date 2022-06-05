package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func init() {
	//Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.Out = os.Stdout
	Logger.SetOutput(os.Stdout)
	Logger.SetLevel(logrus.InfoLevel)
}
