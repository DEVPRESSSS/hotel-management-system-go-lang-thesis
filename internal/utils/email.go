package utils

import (
	"HMS-GO/internal/models"
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
	Year      int
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
	// Parse all templates in the directory
	tmpl, err := ParseTemplateDir("views")
	if err != nil {
		log.Printf("Template parse error: %v\n", err)
		return
	}

	// Execute the specific template into a buffer
	var body bytes.Buffer
	if err := tmpl.ExecuteTemplate(&body, emailTemp, data); err != nil {
		log.Printf("Template execute error: %v\n", err)
		return
	}

	from := mail.NewEmail("devpress", "xmontemorjerald@gmail.com")
	subject := data.Subject
	to := mail.NewEmail(user.Username, user.Email)
	htmlContent := body.String()
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Printf("SendGrid Error: %v\n", err)
	} else {
		log.Printf("Status Code: %d\n", response.StatusCode)
		log.Printf("Headers: %v\n", response.Headers)
	}
}
