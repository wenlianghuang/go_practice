// ref: https://www.golinuxcloud.com/golang-logrus/
package main

import (
	"fmt"
	sublogrus "gologrus/Sublogrus"
	"os"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}
func main() {

	//f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + "log.txt")
		panic(err)
	}
	defer f.Close()
	fmt.Println("Test")
	/*
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
	*/
	log := sublogrus.Sublogrusfunc()
	log.Trace("Trace message")
	log.Info("Info message")
	log.Warn("Warn message")
	log.Error("Error message")
	log.Fatal("Fatal message")

}
