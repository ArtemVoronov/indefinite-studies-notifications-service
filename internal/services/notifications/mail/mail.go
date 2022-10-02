package mail

import (
	"fmt"
	"net"
	"net/smtp"
	"time"
)

// TODO: add tls using
// TODO: add authenication
// TODO: add keep alive connection
type EmailNotificationsService struct {
	smtpServerAddr string
	connectTimeout time.Duration
}

func CreateEmailNotificationsService(smtpServerAddr string, connectTimeout time.Duration) *EmailNotificationsService {
	return &EmailNotificationsService{
		smtpServerAddr: smtpServerAddr,
		connectTimeout: connectTimeout,
	}
}

func (s *EmailNotificationsService) Shutdown() error {
	return nil
}

func (s *EmailNotificationsService) SendEmail(sender string, recipient string, subject string, body string) error {
	msg := "To: " + recipient + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n"

	conn, err := net.DialTimeout("tcp", s.smtpServerAddr, s.connectTimeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	c, err := smtp.NewClient(conn, s.smtpServerAddr)
	if err != nil {
		return err
	}
	defer c.Quit()

	if err := c.Mail(sender); err != nil {
		return err
	}
	if err := c.Rcpt(recipient); err != nil {
		return err
	}

	wc, err := c.Data()
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(wc, msg)
	if err != nil {
		return err
	}
	err = wc.Close()
	if err != nil {
		return err
	}

	return nil
}
