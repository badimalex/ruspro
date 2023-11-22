package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()
	Log.Out = os.Stdout
	Log.SetLevel(logrus.DebugLevel)
}
