package sublogrus

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

func Sublogrusfunc() *logrus.Logger {
	//f, err := os.Create("D:\\MattCode\\go_practice\\gologrus\\output.log")
	f, err := os.OpenFile("D:\\go_practice\\gologrus\\output.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	log := &logrus.Logger{
		// Log into f file handler and on os.Stdout
		Out: io.MultiWriter(f, os.Stdout),
		//Out:   os.Stdout,
		Level: logrus.TraceLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "%lvl%: [%time%] - %msg%\n",
		},
	}

	return log
}
