package mailservice

import (
	logger "AuthenticationService/internal/Helper/Logger"
	"os"

	"gopkg.in/gomail.v2"
)

func MailService(toMailer string, htmlContent string, subject string) bool {
	log := logger.InitLogger()
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAILID"))
	m.SetHeader("To", toMailer)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlContent)

	log.Error(os.Getenv("PASSWORD"))

	d := gomail.NewDialer("smtpout.secureserver.net", 465, os.Getenv("EMAILID"), os.Getenv("PASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		log.Errorf("Could not send email: %v", err)
		return false
	}

	log.Println("Email sent successfully!")
	return true
}
