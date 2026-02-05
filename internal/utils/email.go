package utils

import (
	"HMS-GO/internal/models"
	"bytes"
	"crypto/tls"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

// ? Email template parser

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendEmail(user *models.User, data *EmailData, emailTemp string) {

	// Sender data.
	from := os.Getenv("EMAIL_FROM")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	smtpUser := os.Getenv("SMTP_USER")
	to := user.Email
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := 587

	var body bytes.Buffer

	templatePath := "views"
	template, err := ParseTemplateDir(templatePath)
	if err != nil {
		log.Fatal("Could not parse template", err)
	}

	template.ExecuteTemplate(&body, emailTemp, &data)

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Could not send email: ", err)
	}

}
