package easyMail

type MailSender struct{
  Email string
  Password string
  Provider string
}

const (
  GMAIL = "gmail"
  YAHOO = "yahoo"
  OUTLOOK = "outlook"
)

