package easyMail

import (
	"os"
	"strings"
	"testing"
)

type testCreds struct {
	email    string
	password string
	recivers []string
}

func getTestCreds() testCreds {
	mail := os.Getenv("TEST_EMAIL")
	password := os.Getenv("TEST_EMAIL_PASSWORD")
	recivers := os.Getenv("TEST_EMAIL_RECIVERS")

	if mail == "" || password == "" || recivers == "" {
		panic("Please set TEST_EMAIL, TEST_PASSWORD and TEST_EMAIL_RECIVERS environment variables to run this test")
	}

	reciversList := strings.Split(recivers, ",")

	return testCreds{
		email:    mail,
		password: password,
		recivers: reciversList,
	}
}

func TestSendMail(t *testing.T) {
	creds := getTestCreds()

	mserver := NewMailServer(creds.email, creds.password, GMAIL)

	mail := NewMail(
		creds.email,
		creds.recivers,
		"Normal Test",
		"Nabajit gay, nabajit rohan er ta mukh e nay",
		false,
	)

	err := mserver.SendMail(mail)

	if err != nil {
		t.Error(err)
	}
}

func TestSendMailWithHtml(t *testing.T) {
	creds := getTestCreds()

	mserver := NewMailServer(creds.email, creds.password, GMAIL)

	mail := NewMail(
		"Kaushik",
		creds.recivers,
		"Html Test",
		"<h1>Hello</h1><p>Nabajit is gay </br> nabajit rohan er ta mukh e nay</p>",
		true,
	)

	err := mserver.SendMail(mail)

	if err != nil {
		t.Error(err)
	}
}

func TestSendMailWithAttachment(t *testing.T) {
	creds := getTestCreds()

	mserver := NewMailServer(creds.email, creds.password, GMAIL)

	mail := NewMail(
		"Kaushik Chowdhury",
		creds.recivers,
		"Attachment Test",
		"Noice to meet ya",
		false,
	)

	mail.AddAttachment("testAttachment.md")

	err := mserver.SendMail(mail)

	if err != nil {
		t.Error(err)
	}
}

func TestSendMailWithHtmlFile(t *testing.T) {
	creds := getTestCreds()
	mserver := NewMailServer(creds.email, creds.password, GMAIL)

	mail := NewMail(
		"Kaushik Chowdhury",
		creds.recivers,
		"Html File Test",
		"Noice to meet ya",
		false,
	)

	mail.AddHtmlFile("testHtmlFile.html")

		if err := mserver.SendMail(mail); err != nil {
				t.Error(err)
		}

}

