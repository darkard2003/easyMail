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

	if email == "" || password == "" {
		t.Error("Please set TEST_EMAIL and TEST_PASSWORD environment variables to run this test")
	}

	mserver := NewMailServer(email, password, GMAIL)

	mail := NewMail(
		email,
		reciversList,
		"Normal Test",
		"Noice to meet ya",
		false,
	)

	err := mserver.SendMail(mail)

	if err != nil {
		t.Error(err)
	}
}

func TestSendMailWithHtml(t *testing.T) {
	email := os.Getenv("TEST_EMAIL")
	password := os.Getenv("TEST_EMAIL_PASSWORD")
	recivers := os.Getenv("TEST_EMAIL_RECIVERS")

	reciversList := strings.Split(recivers, ",")

	t.Logf("Email: %s", email)

	if email == "" || password == "" {
		t.Error("Please set TEST_EMAIL and TEST_PASSWORD environment variables to run this test")
	}

	mserver := NewMailServer(email, password, GMAIL)

	mail := NewMail(
		email,
		reciversList,
		"Html Test",
		"<h1>Hello</h1><p>Noice to meet ya</p>",
		true,
	)

	err := mserver.SendMail(mail)

	if err != nil {
		t.Error(err)
	}
}
