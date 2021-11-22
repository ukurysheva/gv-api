package service

import (
	"fmt"
	"net/smtp"
)

func SendMail(to_email, subject, mime, message string) {

	// Sender data.
	from := "globalavia21@gmail.com"
	password := "bwr332gs4"

	// Receiver email address.
	// to := []string{"kuryshevamedia@mail.ru"}
	to := []string{to_email}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.
	messageBytes := []byte(subject + mime + message)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, messageBytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}
