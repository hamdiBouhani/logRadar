package logradar

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
	"sync"
	"time"
)

// EmailSink sends alerts via SMTP.
// The zero value is not usable; always build it with NewEmailSink.
type EmailSink struct {
	host     string
	port     string
	username string
	password string
	from     string
	to       []string
	subject  string // optional prefix
	mu       sync.Mutex
}

// EmailSinkOptions configures an EmailSink.
type EmailSinkOptions struct {
	Host     string   // e.g. "smtp.gmail.com"
	Port     string   // e.g. "587"
	Username string   // SMTP user
	Password string   // SMTP password or app-password
	From     string   // e.g. "logradar@example.com"
	To       []string // recipients
	Subject  string   // optional prefix; default "[logRadar]"
}

func NewEmailSink(opts EmailSinkOptions) *EmailSink {
	subject := opts.Subject
	if subject == "" {
		subject = "[logRadar]"
	}
	return &EmailSink{
		host:     opts.Host,
		port:     opts.Port,
		username: opts.Username,
		password: opts.Password,
		from:     opts.From,
		to:       opts.To,
		subject:  subject,
	}
}

func (e *EmailSink) Send(_ context.Context, alert Alert) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	subj := fmt.Sprintf("%s %s - %s", e.subject, alert.Entry.Level, alert.Reason)
	body := fmt.Sprintf(
		"Time:    %s\r\nLevel:   %s\r\nReason:  %s\r\nMessage: %s\r\n",
		alert.Entry.Timestamp.Format(time.RFC3339),
		alert.Entry.Level,
		alert.Reason,
		alert.Entry.Message,
	)

	msg := "From: " + e.from + "\r\n" +
		"To: " + strings.Join(e.to, ", ") + "\r\n" +
		"Subject: " + subj + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
		"\r\n" + body

	auth := smtp.PlainAuth("", e.username, e.password, e.host)
	return smtp.SendMail(
		e.host+":"+e.port,
		auth,
		e.from,
		e.to,
		[]byte(msg),
	)
}

func (e *EmailSink) Close() error { return nil } // nothing to close
