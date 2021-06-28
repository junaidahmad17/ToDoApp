package email

import (
	"log"
	"net/smtp"
	"os"
)

func SendEmail(user string, subject string, msg string) {

  from := os.Getenv("E_USERNAME")
  password := os.Getenv("E_PASSWORD")

  to := []string{
    user,
  }

  smtpHost := "smtp.gmail.com"
  smtpPort := "587"

  message := []byte("From: the name <"+os.Getenv("E_USERNAME")+">\r\n" +
  "To: "+user+"\r\n" +
  "Subject: "+subject+"\r\n" +
  "\r\n" +
  msg+"\r\n")
  
  auth := smtp.PlainAuth("", from, password, smtpHost)
  
  err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
  if err != nil {
    log.Println(err)
    return
  }
  log.Println("email sent!")
}