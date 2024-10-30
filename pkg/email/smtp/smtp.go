package smtp

import (
	"errors"

	"auth-api/pkg/email"

	"github.com/go-gomail/gomail"
)

type SenderConfig struct {
	Email    string `yaml:"email" env:"EMAIL_SENDER_ADDRESS" env-required:"true"`
	Password string `yaml:"password" env:"EMAIL_SENDER_ADDRESS" env-required:"true"`
	Host     string `yaml:"host" env:"EMAIL_SENDER_ADDRESS" env-required:"true"`
	Port     int    `yaml:"port" env:"EMAIL_SENDER_ADDRESS" env-required:"true"`
}

type SMTPSender struct {
	from string
	pass string
	host string
	port int
}

func NewSMTPSender(senderConfig *SenderConfig) (*SMTPSender, error) {
	if !email.IsEmailValid(senderConfig.Email) {
		return nil, errors.New("invalid from email")
	}

	return &SMTPSender{from: senderConfig.Email, pass: senderConfig.Password, host: senderConfig.Host, port: senderConfig.Port}, nil
}

func (s *SMTPSender) Send(input email.SendEmailInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", s.from)
	msg.SetHeader("To", input.To)
	msg.SetHeader("Subject", input.Subject)
	msg.SetBody("text/html", input.Body)

	dialer := gomail.NewDialer(s.host, s.port, s.from, s.pass)

	return dialer.DialAndSend(msg)
}
