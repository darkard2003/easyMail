package easyMail

import (
	"os"
	"strings"
	"testing"
)

func TestSendMail(t *testing.T) {
	email := os.Getenv("TEST_EMAIL")
	password := os.Getenv("TEST_EMAIL_PASSWORD")
	recivers := os.Getenv("TEST_EMAIL_RECIVERS")

	reciversList := strings.Split(recivers, ",")

	t.Logf("Email: %s", email)

	if email == "" || password == "" || recivers == "" {
		t.Error("Please set TEST_EMAIL and TEST_PASSWORD and TEST_EMAIL_RECIVER environment variables to run this test")
	}

	mserver := NewMailServer(email, password, GMAIL)

	mail := NewMail(
		email,
		reciversList,
		"Test Subject",
		"Test Body",
		false,
	)

	err := mserver.SendMail(mail)

	if err != nil {
		t.Error(err)
	}
}
