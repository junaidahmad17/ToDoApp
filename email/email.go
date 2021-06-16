package email

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(user string, subject string, msg string) {

  // Sender data.
  from := os.Getenv("E_USERNAME")
  password := os.Getenv("E_PASSWORD")

  // Receiver email address.
  to := []string{
    user,
  }

  // smtp server configuration.
  smtpHost := "smtp.gmail.com"
  smtpPort := "587"

  // Message.
  message := []byte("From: the name <"+os.Getenv("E_USERNAME")+">\r\n" +
  "To: "+user+"\r\n" +
  "Subject: "+subject+"\r\n" +
  "\r\n" +
  msg+"\r\n")
  
  // Authentication.
  auth := smtp.PlainAuth("", from, password, smtpHost)
  
  // Sending email.
  err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println("Email Sent Successfully!")
}