package email

import (
	"bytes"
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"
)

// Config holds email service configuration
type Config struct {
	SMTPHost            string
	SMTPPort            string
	SMTPUsername        string
	SMTPPassword        string
	FromEmail           string
	FromName            string
	BaseURL             string // Base URL for email links (e.g., https://app.lightshare.com)
	MobileDeepLinkScheme string // Custom URL scheme for mobile deep links (e.g., lightshare)
}

// Service handles email sending
type Service struct {
	config Config
	dialer *gomail.Dialer
}

// New creates a new email service
func New(cfg *Config) *Service {
	port, err := strconv.Atoi(cfg.SMTPPort)
	if err != nil {
		port = 587 // default to standard SMTP submission port
	}

	dialer := gomail.NewDialer(cfg.SMTPHost, port, cfg.SMTPUsername, cfg.SMTPPassword)
	// Use SSL for port 465, STARTTLS for others (587, 25)
	dialer.SSL = (port == 465)

	return &Service{
		config: *cfg,
		dialer: dialer,
	}
}

// EmailMessage represents an email to send
type EmailMessage struct {
	To      string
	Subject string
	Body    string
	IsHTML  bool
}

// Send sends an email using gomail (supports OVH and other SMTP providers)
func (s *Service) Send(msg EmailMessage) error {
	m := gomail.NewMessage()

	// Set headers
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail))
	m.SetHeader("To", msg.To)
	m.SetHeader("Subject", msg.Subject)

	// Set body
	if msg.IsHTML {
		m.SetBody("text/html", msg.Body)
	} else {
		m.SetBody("text/plain", msg.Body)
	}

	// Send email
	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendVerificationEmail sends an email verification email
func (s *Service) SendVerificationEmail(to, token string) error {
	verificationURL := fmt.Sprintf("%s://verify-email?token=%s", s.config.MobileDeepLinkScheme, token)

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Verify Your Email</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #2563eb;">Welcome to LightShare!</h1>
        <p>Thank you for signing up. Please verify your email address by clicking the button below:</p>
        <div style="text-align: center; margin: 30px 0;">
            <a href="{{.URL}}" style="background-color: #2563eb; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block;">
                Verify Email
            </a>
        </div>
        <p style="color: #666; font-size: 14px;">
            Or copy and paste this link into your browser:<br>
            <a href="{{.URL}}">{{.URL}}</a>
        </p>
        <p style="color: #666; font-size: 14px;">
            This link will expire in 24 hours. If you didn't create an account with LightShare, you can safely ignore this email.
        </p>
    </div>
</body>
</html>
`

	t, err := template.New("verification").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var body bytes.Buffer
	if err := t.Execute(&body, map[string]string{"URL": verificationURL}); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return s.Send(EmailMessage{
		To:      to,
		Subject: "Verify your LightShare email",
		Body:    body.String(),
		IsHTML:  true,
	})
}

// SendMagicLinkEmail sends a magic link login email
func (s *Service) SendMagicLinkEmail(to, token string) error {
	magicLinkURL := fmt.Sprintf("%s://magic-link?token=%s", s.config.MobileDeepLinkScheme, token)

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Your Magic Link</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #2563eb;">Login to LightShare</h1>
        <p>Click the button below to securely log in to your account:</p>
        <div style="text-align: center; margin: 30px 0;">
            <a href="{{.URL}}" style="background-color: #2563eb; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block;">
                Login to LightShare
            </a>
        </div>
        <p style="color: #666; font-size: 14px;">
            Or copy and paste this link into your browser:<br>
            <a href="{{.URL}}">{{.URL}}</a>
        </p>
        <p style="color: #666; font-size: 14px;">
            This link will expire in 15 minutes. If you didn't request this login link, you can safely ignore this email.
        </p>
    </div>
</body>
</html>
`

	t, err := template.New("magiclink").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var body bytes.Buffer
	if err := t.Execute(&body, map[string]string{"URL": magicLinkURL}); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return s.Send(EmailMessage{
		To:      to,
		Subject: "Your LightShare login link",
		Body:    body.String(),
		IsHTML:  true,
	})
}

// SendPasswordResetEmail sends a password reset email
func (s *Service) SendPasswordResetEmail(to, token string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.config.BaseURL, token)

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Reset Your Password</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #2563eb;">Reset Your Password</h1>
        <p>You requested to reset your password. Click the button below to create a new password:</p>
        <div style="text-align: center; margin: 30px 0;">
            <a href="{{.URL}}" style="background-color: #2563eb; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block;">
                Reset Password
            </a>
        </div>
        <p style="color: #666; font-size: 14px;">
            Or copy and paste this link into your browser:<br>
            <a href="{{.URL}}">{{.URL}}</a>
        </p>
        <p style="color: #666; font-size: 14px;">
            This link will expire in 1 hour. If you didn't request a password reset, you can safely ignore this email.
        </p>
    </div>
</body>
</html>
`

	t, err := template.New("reset").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var body bytes.Buffer
	if err := t.Execute(&body, map[string]string{"URL": resetURL}); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return s.Send(EmailMessage{
		To:      to,
		Subject: "Reset your LightShare password",
		Body:    body.String(),
		IsHTML:  true,
	})
}

// ValidateEmail performs basic email validation
func ValidateEmail(email string) bool {
	email = strings.TrimSpace(strings.ToLower(email))
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	if !strings.Contains(email, "@") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if parts[0] == "" || len(parts[1]) < 3 {
		return false
	}
	if !strings.Contains(parts[1], ".") {
		return false
	}
	return true
}
