package services

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"mime"
	"net"
	"net/smtp"
	"strings"
	"time"
)

type EmailService struct {
	svc *CompanySettingsService
}

type EmailAttachment struct {
	Filename    string
	ContentType string
	Data        []byte
}

func NewEmailService(svc *CompanySettingsService) *EmailService {
	return &EmailService{svc: svc}
}

func (s *EmailService) SendTestEmail(ctx context.Context, to, name string) error {
	subject := "Test Email — FreeFSM"
	body := fmt.Sprintf(`Hi %s,

This is a test email from FreeFSM. Your SMTP settings are working correctly.

If you received this email, your email configuration is ready to use.

- FreeFSM`, name)

	return s.SendEmail(ctx, to, subject, body)
}

func (s *EmailService) SendWelcomeEmail(ctx context.Context, to, name, tempPassword, loginURL string) error {
	subject := "Welcome to FreeFSM"
	body := fmt.Sprintf(`Hi %s,

Your FreeFSM account has been created.

Login: %s
Email: %s
Temporary Password: %s

You'll be required to change this on first login.

- FreeFSM`, name, loginURL, to, tempPassword)

	return s.SendEmail(ctx, to, subject, body)
}

func (s *EmailService) SendPasswordReset(ctx context.Context, to, name, link string) error {
	subject := "Password Reset - FreeFSM"
	body := fmt.Sprintf(`Hi %s,

A password reset was requested for your account.

Click the link below to reset your password:
%s

This link expires in 1 hour.

If you did not request this, please ignore this email.

- FreeFSM`, name, link)

	return s.SendEmail(ctx, to, subject, body)
}

func (s *EmailService) SendEmailWithAttachment(ctx context.Context, to, subject, body, filename, mimeType string, data []byte) error {
	return s.SendEmail(ctx, to, subject, body, EmailAttachment{
		Filename:    filename,
		ContentType: mimeType,
		Data:        data,
	})
}

func (s *EmailService) SendEmail(ctx context.Context, to, subject, body string, attachments ...EmailAttachment) error {
	cs, err := s.svc.Get(ctx)
	if err != nil || cs == nil || cs.SMTPHost == "" {
		return fmt.Errorf("SMTP not configured")
	}

	from := sanitizeHeader(cs.SMTPFrom)
	to = sanitizeHeader(to)
	if strings.TrimSpace(to) == "" {
		return fmt.Errorf("recipient email is required")
	}

	msg, err := buildEmailMessage(from, to, subject, body, attachments)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%d", cs.SMTPHost, cs.SMTPPort)
	if cs.SMTPPort == 465 {
		return s.sendTLS(addr, cs.SMTPUser, cs.SMTPPassword, from, to, msg)
	}
	if cs.SMTPPort == 587 {
		return s.sendWithSTARTTLS(addr, cs.SMTPUser, cs.SMTPPassword, from, to, msg)
	}
	return s.sendPlain(addr, cs.SMTPHost, cs.SMTPUser, cs.SMTPPassword, from, to, msg)
}

func RenderEmailTemplate(template string, values map[string]string) string {
	for key, value := range values {
		template = strings.ReplaceAll(template, "{"+key+"}", value)
	}
	return template
}

func buildEmailMessage(from, to, subject, body string, attachments []EmailAttachment) (string, error) {
	from = sanitizeHeader(from)
	to = sanitizeHeader(to)
	subject = sanitizeHeader(subject)

	if len(attachments) == 0 {
		return fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
			from, to, mime.QEncoding.Encode("UTF-8", subject), body), nil
	}

	boundary := fmt.Sprintf("freefsm-%d", time.Now().UnixNano())
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary=%q\r\n\r\n", from, to, mime.QEncoding.Encode("UTF-8", subject), boundary)
	fmt.Fprintf(&buf, "--%s\r\nContent-Type: text/plain; charset=UTF-8\r\nContent-Transfer-Encoding: 8bit\r\n\r\n%s\r\n", boundary, body)

	for _, attachment := range attachments {
		filename := sanitizeHeader(attachment.Filename)
		if filename == "" {
			return "", fmt.Errorf("attachment filename is required")
		}
		contentType := sanitizeHeader(attachment.ContentType)
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		disposition := mime.FormatMediaType("attachment", map[string]string{"filename": filename})
		fmt.Fprintf(&buf, "\r\n--%s\r\nContent-Type: %s\r\nContent-Transfer-Encoding: base64\r\nContent-Disposition: %s\r\n\r\n", boundary, contentType, disposition)
		writeBase64Lines(&buf, attachment.Data)
		buf.WriteString("\r\n")
	}
	fmt.Fprintf(&buf, "--%s--\r\n", boundary)

	return buf.String(), nil
}

func sanitizeHeader(value string) string {
	value = strings.ReplaceAll(value, "\r", "")
	return strings.ReplaceAll(value, "\n", "")
}

func writeBase64Lines(buf *bytes.Buffer, data []byte) {
	encoded := base64.StdEncoding.EncodeToString(data)
	for len(encoded) > 76 {
		buf.WriteString(encoded[:76])
		buf.WriteString("\r\n")
		encoded = encoded[76:]
	}
	buf.WriteString(encoded)
}

func (s *EmailService) sendWithSTARTTLS(addr, user, password, from, to, msg string) error {
	host := strings.Split(addr, ":")[0]
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(10 * time.Second))

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("SMTP client: %w", err)
	}
	defer client.Quit()

	tlsConfig := &tls.Config{ServerName: host}
	if ok, _ := client.Extension("STARTTLS"); ok {
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("STARTTLS: %w", err)
		}
	}

	return sendSMTPMessage(client, host, user, password, from, to, msg)
}

func (s *EmailService) sendTLS(addr, user, password, from, to, msg string) error {
	host := strings.Split(addr, ":")[0]

	tlsConfig := &tls.Config{ServerName: host}

	dialer := &net.Dialer{Timeout: 10 * time.Second}
	conn, err := tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS connect: %w", err)
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(10 * time.Second))

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("SMTP client: %w", err)
	}
	defer client.Quit()

	return sendSMTPMessage(client, host, user, password, from, to, msg)
}

func (s *EmailService) sendPlain(addr, host, user, password, from, to, msg string) error {
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(10 * time.Second))

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("SMTP client: %w", err)
	}
	defer client.Quit()

	return sendSMTPMessage(client, host, user, password, from, to, msg)
}

func sendSMTPMessage(client *smtp.Client, host, user, password, from, to, msg string) error {
	if user != "" {
		auth := smtp.PlainAuth("", user, password, host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP auth: %w", err)
		}
	}
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("MAIL FROM: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("RCPT TO: %w", err)
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA: %w", err)
	}
	_, err = w.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	w.Close()

	return nil
}
