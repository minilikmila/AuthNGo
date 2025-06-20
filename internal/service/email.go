package service

import (
	"bytes"
	"context"
	"html/template"
	"net/smtp"

	config "github.com/minilikmila/standard-auth-go/configs"
	"github.com/sirupsen/logrus"
)

type EmailServiceImpl struct {
	config *config.Config
}

func NewEmailService(config *config.Config) EmailService {
	return &EmailServiceImpl{
		config: config,
	}
}

// sendEmailFromTemplate parses the template, injects data, and sends the email
func (s *EmailServiceImpl) sendEmailFromTemplate(ctx context.Context, to, subject, templatePath string, data interface{}) error {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	smtpHost := s.config.SMTPHost
	smtpPort := s.config.SMTPPort
	smtpUser := s.config.SMTPUser
	smtpPass := s.config.SMTPPass
	from := s.config.SMTPFrom

	logrus.Infoln("smtp envs: ", smtpHost, smtpPass, smtpUser, s.config.CompanyName)

	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"
	msg += "From: " + from + "\n"
	msg += "To: " + to + "\n"
	msg += "Subject: " + subject + "\n\n"
	msg += body.String()

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	addr := smtpHost + ":" + smtpPort
	logrus.Infoln("smtp:address: ", addr)
	if err = smtp.SendMail(addr, auth, from, []string{to}, []byte(msg)); err != nil {
		logrus.Errorln("SMTP sending error: ", err)
		return err
	}
	return nil
}

func (s *EmailServiceImpl) SendVerificationEmail(ctx context.Context, email, token, receiverName string) error {
	if receiverName == "" {
		receiverName = "User"
	}
	data := map[string]interface{}{
		"Name":            receiverName, // Replace with actual name if available
		"VerificationURL": s.config.AppURL + "/verify?token=" + token,
		"Company":         s.config.CompanyName,
	}
	return s.sendEmailFromTemplate(ctx, email, "Verify your email address", "templates/email_verification.html", data)
}

func (s *EmailServiceImpl) SendPasswordResetEmail(ctx context.Context, email, token string) error {
	data := map[string]interface{}{
		"Name":     "User", // Replace with actual name if available
		"ResetURL": s.config.AppURL + "/reset-password?token=" + token,
		"Company":  s.config.CompanyName,
	}
	return s.sendEmailFromTemplate(ctx, email, "Reset your password", "templates/password_reset.html", data)
}

func (s *EmailServiceImpl) SendWelcomeEmail(ctx context.Context, email string, name string) error {
	data := map[string]interface{}{
		"Name":    name,
		"Company": s.config.CompanyName,
	}
	return s.sendEmailFromTemplate(ctx, email, "Welcome to "+s.config.CompanyName, "templates/welcome.html", data)
}
