package main

import (
	"log"
	"net/smtp"

	"github.com/jordan-wright/email"
)

func main() {
	// email account and password(once)
	smtpServer := "smtp.gmail.com"
	senderEmail := "wenliangmatt@gmail.com"
	senderPassword := "cylx tdls gtou pbsr"

	// Receiver
	toEmail := "wenliangmatt@gmail.com"

	// Build a new email
	e := email.NewEmail()
	e.From = senderEmail
	e.To = []string{toEmail}
	e.Subject = "Hello, Golang Gmail Example"
	e.Text = []byte("This is the email body text.")

	// Send Email
	err := e.Send(smtpServer+":587", smtp.PlainAuth("", senderEmail, senderPassword, smtpServer))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Email sent successfully!")
}
