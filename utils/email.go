package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/mail"
	"net/smtp"
	"strings"

	"gopkg.in/gomail.v2"
)

type Mail struct {
	Sender      string
	To          []string
	Cc          []string
	Subject     string
	Body        string
	Attachments map[string][]byte
}

type SmtpConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type EmailAttachment struct {
	FileName string
	FileMine string
	Data     []byte
}

type EmailContent struct {
	From        mail.Address
	To          []mail.Address
	Cc          []mail.Address
	Bcc         []mail.Address
	Subject     string
	Body        string
	ContentType string
	Attachments []EmailAttachment
}

func (attachment *EmailAttachment) WriteAttachment(w io.Writer) error {
	_, err := w.Write([]byte(attachment.Data))
	return err
}

func emailsToString(emails []mail.Address) string {
	if len(emails) > 0 {
		tmp := make([]string, len(emails))
		for idx, item := range emails {
			tmp[idx] = item.String()
		}
		return strings.Join(tmp, "; ")
	}
	return ""
}

func emailOnly(emails []mail.Address) []string {
	returnVal := make([]string, len(emails))
	for idx, item := range emails {
		returnVal[idx] = item.Address
	}
	return returnVal
}

func buildEmailBody(data EmailContent) []byte {
	buf := bytes.NewBuffer(nil)
	withAttachments := len(data.Attachments) > 0

	// build header
	buf.WriteString(fmt.Sprintf("Subject: %s\n", data.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\n", emailsToString(data.To)))
	if len(data.Cc) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\n", emailsToString(data.Cc)))
	}

	if len(data.Bcc) > 0 {
		buf.WriteString(fmt.Sprintf("Bcc: %s\n", emailsToString(data.Bcc)))
	}

	buf.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	} else {
		if strings.TrimSpace(data.ContentType) == "" {
			data.ContentType = "text/html"
		}
		buf.WriteString(fmt.Sprintf("Content-Type: %v; charset=utf-8\n", data.ContentType))
	}

	buf.WriteString(data.Body)
	if withAttachments {
		for _, attachment := range data.Attachments {
			fileMine := attachment.FileMine
			if strings.TrimSpace(fileMine) == "" {
				fileMine = http.DetectContentType(attachment.Data)
			}
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", fileMine))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", attachment.FileName))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(attachment.Data)))
			base64.StdEncoding.Encode(b, attachment.Data)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}

		buf.WriteString("--")
	}
	return buf.Bytes()
}

func SendMailWithData(senderConfig SmtpConfig, data EmailContent) error {
	return SendRawMail(senderConfig, data.From.Address, emailOnly(data.To), buildEmailBody(data))
}

func buildListAddresses(msg *gomail.Message, addresses []mail.Address) []string {
	returnVal := []string{}
	for _, v := range addresses {
		if len(strings.TrimSpace(v.Address)) > 0 {
			if len(strings.TrimSpace(v.Name)) > 0 {
				returnVal = append(returnVal, msg.FormatAddress(v.Address, v.Name))
			} else {
				returnVal = append(returnVal, v.Address)
			}
		}
	}
	return returnVal
}

func SendMailWithDataV2(senderConfig SmtpConfig, data EmailContent) error {
	m := gomail.NewMessage()
	if len(strings.TrimSpace(data.From.Name)) == 0 {
		m.SetHeader("From", data.From.Address)
	} else {
		m.SetAddressHeader("From", data.From.Address, data.From.Name)
	}

	// to
	to := buildListAddresses(m, data.To)
	m.SetHeader("To", to...)

	// cc
	cc := buildListAddresses(m, data.Cc)
	if len(cc) > 0 {
		m.SetHeader("Cc", cc...)
	}

	// bcc
	bcc := buildListAddresses(m, data.Bcc)
	if len(bcc) > 0 {
		m.SetHeader("Bcc", bcc...)
	}

	m.SetHeader("Subject", data.Subject)

	m.SetBody(data.ContentType, data.Body)

	for idx, _ := range data.Attachments {
		m.Attach(
			data.Attachments[idx].FileName,
			gomail.SetCopyFunc(data.Attachments[idx].WriteAttachment),
		)
	}
	// Send the email
	d := gomail.NewDialer(senderConfig.Host, senderConfig.Port, senderConfig.Username, senderConfig.Password) // Replace with your SMTP details

	return d.DialAndSend(m)
}

func SendRawMail(senderConfig SmtpConfig, fromEmail string, to []string, msg []byte) error {
	smtpAuth := smtp.PlainAuth("", senderConfig.Username, senderConfig.Password, senderConfig.Host)
	// mail server address
	smtpAddress := fmt.Sprintf("%v:%v", senderConfig.Host, senderConfig.Port)
	// send mail
	if strings.TrimSpace(fromEmail) == "" {
		fromEmail = senderConfig.Username
	}
	return smtp.SendMail(smtpAddress, smtpAuth, fromEmail, to, msg)
}

func FormatEmailAddress(emailAddress, name string) string {
	return gomail.NewMessage().FormatAddress(emailAddress, name)
}
