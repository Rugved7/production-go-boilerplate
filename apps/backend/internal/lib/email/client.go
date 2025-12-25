package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/Rugved7/go-boilerplate/internal/config"
	"github.com/resendlabs/resend-go"
	"github.com/rs/zerolog"
)

type Client struct {
	client *resend.Client
	logger *zerolog.Logger
}

func NewClient(cfg *config.Config, logger *zerolog.Logger) *Client {
	return &Client{
		client: resend.NewClient(cfg.Integration.ResendAPIKey),
		logger: logger,
	}
}

func (c *Client) SendEmail(to, subject string, templateName Template, data map[string]string) error {
	tmplPath := fmt.Sprintf("templates/emails/%s.html", templateName)

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse email template %s: %w", templateName, err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template %s: %w", templateName, err)
	}

	params := &resend.SendEmailRequest{
		From:    "Alfred <onboarding@resend.dev>",
		To:      []string{to},
		Subject: subject,
		Html:    body.String(),
	}

	if _, err := c.client.Emails.Send(params); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
