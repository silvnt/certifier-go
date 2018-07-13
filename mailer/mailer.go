package mailer

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Sender contains the information needed to send emails
type Sender struct {
	Name               string
	Address            string
	DefaultSubject     string
	DefaultTextContent string
}

// SendEmail sends emails given the remitee's information & filename to be sent
func (s *Sender) SendEmail(remiteeName, remiteeAddress, filename string) error {
	m := mail.NewV3Mail()

	from := mail.NewEmail(s.Name, s.Address)
	content := mail.NewContent("text/html", "<p>"+s.DefaultTextContent+"</p>")
	to := mail.NewEmail(remiteeName, remiteeAddress)
	m.SetFrom(from)
	m.AddContent(content)

	personalization := mail.NewPersonalization()
	personalization.AddTos(to)
	personalization.Subject = s.DefaultSubject
	m.AddPersonalizations(personalization)

	aPDF := mail.NewAttachment()
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(dat))
	aPDF.SetContent(encoded)
	aPDF.SetType("application/pdf")
	aPDF.SetFilename(filename)
	aPDF.SetDisposition("attachment")
	aPDF.SetContentID("Test Attachment")

	m.AddAttachment(aPDF)

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

	return err
}
