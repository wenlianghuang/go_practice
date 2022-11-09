package sublogrus

import (
	"os"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

func Sublogrusfunc() *logrus.Logger {
	log := &logrus.Logger{
		// Log into f file handler and on os.Stdout
		//Out:   io.MultiWriter(f, os.Stdout),
		Out:   os.Stdout,
		Level: logrus.TraceLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006/01/02 15:04:05",
			LogFormat:       "%lvl%: [%time%] - %msg%\n",
		},
	}

	return log
}
