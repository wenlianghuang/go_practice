package sublogrus

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

func Sublogrusfunc(relativePath string) *logrus.Logger {
	//f, err := os.Create("D:\\MattCode\\go_practice\\gologrus\\output.log")
	//f, err := os.OpenFile("D:\\go_practice\\gologrus\\output.log",os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// 取得當前日期時間
	currentTime := time.Now()
	// 使用指定的日期時間格式化字串來建立檔案名稱
	fileName := fmt.Sprintf("%d-%02d-%02d_%02d-%02d-%02d.log", currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), currentTime.Second())
	//Check the folder of log is existed

	if _, err := os.Stat("log"); os.IsNotExist(err) {
		err = os.MkdirAll("log", 0755)
		if err != nil {
			panic(err)
		}
	}

	f, err := os.OpenFile(relativePath+"\\log\\"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
