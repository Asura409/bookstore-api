package services

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"bytes"
	"github.com/resend/resend-go/v2" // Or your preferred email provider
)

type EmailService struct {
	client         *resend.Client
	fromEmail      string
	templates      map[string]*template.Template
	templateDir    string
}

// NewEmailService creates a configured email service
func NewEmailService(apiKey, fromEmail, templateDir string) (*EmailService, error) {
	// Initialize email client (Resend in this example)
	client := resend.NewClient(apiKey)

	// Preload templates
	templates, err := loadTemplates(templateDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load email templates: %w", err)
	}

	return &EmailService{
		client:      client,
		fromEmail:   fromEmail,
		templates:   templates,
		templateDir: templateDir,
	}, nil
}

// Email Types (enums for template selection)
const (
	EmailTypePasswordReset = "password_reset"
	EmailTypeVerification  = "email_verification"
	EmailTypeWelcome       = "welcome"
)

// SendEmail sends a templated email
func (s *EmailService) SendEmail(to, subject, emailType string, data map[string]interface{}) error {
	tmpl, exists := s.templates[emailType]
	if !exists {
		return fmt.Errorf("email template %s not found", emailType)
	}

	// Execute template
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

	// Send via email provider
	_, err := s.client.Emails.Send(&resend.SendEmailRequest{
		From:    s.fromEmail,
		To:      []string{to},
		Subject: subject,
		Html:    body.String(),
	})

	return err
}

// loadTemplates preloads all email templates from disk
func loadTemplates(templateDir string) (map[string]*template.Template, error) {
	templates := make(map[string]*template.Template)

	// Load all .html files in the template directory
	files, err := os.ReadDir(templateDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".html" {
			name := file.Name()[:len(file.Name())-5] // Remove .html extension
			path := filepath.Join(templateDir, file.Name())
			tmpl, err := template.ParseFiles(path)
			if err != nil {
				return nil, fmt.Errorf("failed to parse template %s: %w", name, err)
			}
			templates[name] = tmpl
		}
	}

	return templates, nil
}