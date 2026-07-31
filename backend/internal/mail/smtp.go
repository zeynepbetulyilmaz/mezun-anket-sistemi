package mail

import (
	"fmt"
	"net/smtp"
)

type SMTPClient struct {
	host string
	port string
	user string
	pass string
	from string
}

func NewSMTPClient(host, port, user, pass, from string) *SMTPClient {
	return &SMTPClient{host: host, port: port, user: user, pass: pass, from: from}
}

func (s *SMTPClient) Send(toEmail, subject, body string) error {
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		s.from, toEmail, subject, body)

	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	var auth smtp.Auth
	if s.user != "" {
		auth = smtp.PlainAuth("", s.user, s.pass, s.host)
	}
	return smtp.SendMail(addr, auth, s.from, []string{toEmail}, []byte(msg))
}
