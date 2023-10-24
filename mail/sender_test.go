package mail

import (
	"log"
	"testing"
)

func TestMail(t *testing.T) {
	gmailSender := NewGmailSender("Fiber Api", "elishasilas87@gmail.com", "dkuwfacjzvjfhsja")
	subject := "A test email"
	content := `
	<h1>Hello, world!</h1>
	`
	to := []string{"silaselisha66@gmail.com"}
	err := gmailSender.SendEmail(nil, to, nil, subject, content)

	if err != nil {
		log.Fatal(err)
	}
}
