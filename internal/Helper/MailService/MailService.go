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

	// Use Gmail SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("EMAILID"), os.Getenv("PASSWORD"))

	// Use TLS explicitly (optional, 587 already uses STARTTLS)
	d.TLSConfig = nil

	if err := d.DialAndSend(m); err != nil {
		log.Errorf("Could not send email: %v", err)
		return false
	}

	log.Println("Email sent successfully!")
	return true
}
