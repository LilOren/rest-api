package shared

import (
	"fmt"

	"github.com/jordan-wright/email"
)

func MakeEmail(senderName, senderAddress, subject, content, toAddress string) *email.Email {
	mail := email.NewEmail()
	mail.From = fmt.Sprintf("%s <%s>", senderName, senderAddress)
	mail.Subject = subject
	mail.HTML = []byte(content)
	mail.To = []string{toAddress}
	return mail
}
