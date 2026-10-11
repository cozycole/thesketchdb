package mail

import (
	"bytes"
	"context"
	"embed"
	"html/template"
	"time"

	"github.com/mailgun/mailgun-go/v5"
)

//go:embed "templates"
var templateFS embed.FS

type MailerInterface interface {
	Send(string, string, any) error
}

type Mailer struct {
	client *mailgun.Client
	domain string
	sender string
}

func New(domain, api_key, sender string) Mailer {
	mg := mailgun.NewMailgun(api_key)
	return Mailer{
		client: mg,
		domain: domain,
		sender: sender,
	}
}

// Define a Send() method on the Mailer type. This takes the recipient email address
// as the first parameter, the name of the file containing the templates, and any
// dynamic data for the templates as an any parameter.
func (m Mailer) Send(recipient, templateFile string, data any) error {
	// Use the ParseFS() method to parse the required template file from the embedded
	// file system.
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	// Execute the named template "subject", passing in the dynamic data and storing the
	// result in a bytes.Buffer variable.
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	// Follow the same pattern to execute the "plainBody" template and store the result
	// in the plainBody variable.
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := mailgun.NewMessage(
		m.domain,
		m.sender,
		subject.String(),
		plainBody.String(),
		recipient,
	)

	msg.SetHTML(htmlBody.String())

	for i := 1; i <= 3; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		_, err = m.client.Send(ctx, msg)

		if nil == err {
			return nil
		}

		time.Sleep(10 * time.Second)
	}

	return err
}
