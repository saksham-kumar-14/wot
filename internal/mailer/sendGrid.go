package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendGrid(apiKey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)
	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (m *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	var subject bytes.Buffer
	if err := tmpl.ExecuteTemplate(&subject, "subject", data); err != nil {
		return fmt.Errorf("failed to execute subject template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.ExecuteTemplate(&body, "body", data); err != nil {
		return fmt.Errorf("failed to execute body template: %w", err)
	}

	msg := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	msg.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	for i := 0; i < maxRetries; i++ {
		res, err := m.client.Send(msg)
		if err != nil {
			log.Printf("Failed to send email to %v, attempt %d of %d", email, i+1, maxRetries)
			log.Printf("Error: %v", err.Error())

			time.Sleep(time.Second * time.Duration(i+1))
		}
		log.Printf("Email sent with status code %v", res.StatusCode)
		return nil
	}

	return fmt.Errorf("failed to send email after %d attempts", maxRetries)
}
