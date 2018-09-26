package v1

import (
	"log"
	"net/smtp"

	logger "github.com/adcamie/adserver/logger"
)

var from = "camiemail@adcamie.com"
var pass = "!adcamie12"

func SendEmail(subject string, body string, to []string) {

	msg := "From: " + from + "\n" +
		"Subject:" + subject + "\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, to, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		go logger.ErrorLogger(err.Error(), "EmailSending", "Failed to send email:-"+subject)
		return
	}

	log.Print("Email Sent Successfully")
}
