// 有可能在公司無法寄出，回家再試試看
package sendemailex

import (
	"fmt"
	"net/smtp"
)

func Basicinit() {
	host := "smtp.gmail.com"
	port := 587
	email := "wenliangmatt@gmail.com"
	password := "wenliang75"
	//toEmail := "wenliangmatt@gmail.com"
	toEmail := "matt200335@yahoo.com.tw"
	header := make(map[string]string)
	header["From"] = "test" + "<" + email + ">"
	header["To"] = toEmail
	header["Subject"] = "電子郵件標題"
	header["Content-Type"] = "text/html; charset=UTF-8"
	body := "我是一封電子郵件!golang發出."
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body
	auth := smtp.PlainAuth(
		"",
		email,
		password,
		host,
	)
	/*err := SendMailUsingTLS(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		email,
		[]string{toEmail},
		[]byte(message),
	)*/
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		email,
		[]string{toEmail},
		[]byte(message),
	)
	if err != nil {
		panic(err)
	}
}

/*
func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Println("Dialing Error:", err)
		return nil, err
	}
	//分解主機端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

//參考net/smtp的func SendMail()
//使用net.Dial連接tls(ssl)端口时,smtp.NewClient()會卡住且不提示err
//len(to)>1时,to[1]開始提示是密送
func SendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {
	//create smtp client
	c, err := Dial(addr)
	if err != nil {
		log.Println("Create smpt client error:", err)
		return err
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
*/
