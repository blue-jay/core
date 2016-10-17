// Package email provides email sending via SMTP.
package email

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
)

// Info holds the details for the SMTP server.
type Info struct {
	Username string
	Password string
	Hostname string
	Port     int
	From     string
}

// Send an email.
func (c Info) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", c.Username, c.Password, c.Hostname)

	// Create the header
	header := make(map[string]string)
	header["From"] = c.From
	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = `text/plain; charset="utf-8"`
	header["Content-Transfer-Encoding"] = "base64"

	// Set the message
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	// Send the email
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", c.Hostname, c.Port),
		auth,
		c.From,
		[]string{to},
		[]byte(message),
	)

	return err
}
