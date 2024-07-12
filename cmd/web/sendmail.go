package main


import (
    "fmt"
    "net/smtp"
)


type EmailService struct {
    smtpServer string
    smtpPort   int
    username   string
    password   string
}

func NewEmailService(server string, port int, username, password string) *EmailService {
    return &EmailService{
        smtpServer: server,
        smtpPort:   port,
        username:   username,
        password:   password,
    }
}

func (s *EmailService) SendEmail(to, subject, body string) error {
    auth := smtp.PlainAuth("", s.username, s.password, s.smtpServer)

    msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)

    err := smtp.SendMail(
        fmt.Sprintf("%s:%d", s.smtpServer, s.smtpPort),
        auth,
        s.username,
        []string{to},
        []byte(msg),
    )

    if err != nil {
        return fmt.Errorf("failed to send email: %w", err)
    }

    return nil
}