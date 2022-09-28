package mail

import (
	"fmt"
	"net/smtp"
)

// TODO: add authenication
type EmailNotificationsService struct {
	host string
}

func CreateEmailNotificationsService(host string) *EmailNotificationsService {
	return &EmailNotificationsService{
		host: host,
	}
}

func (s *EmailNotificationsService) Shutdown() error {
	return nil
}

func (s *EmailNotificationsService) Send(sender string, recipient string, body string) error {
	c, err := smtp.Dial(s.host)
	if err != nil {
		return err
	}

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
	_, err = fmt.Fprint(wc, body)
	if err != nil {
		return err
	}
	err = wc.Close()
	if err != nil {
		return err
	}

	err = c.Quit()
	if err != nil {
		return err
	}

	return nil
}
