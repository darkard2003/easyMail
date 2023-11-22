package easyMail

import "net/smtp"

type MailServer struct {
	Sender MailSender
	Auth   smtp.Auth
}

func (s *MailSender) getHost() string {
	switch s.Provider {
	case GMAIL:
		return "smtp.gmail.com"
	case YAHOO:
		return "smtp.mail.yahoo.com"
	case OUTLOOK:
		return "smtp.live.com"
	default:
		return ""
	}
}

func (s *MailSender) getAuth() smtp.Auth {
	return smtp.PlainAuth("", s.Email, s.Password, s.getHost())
}

func NewMailServer(sender MailSender) *MailServer {
	return &MailServer{
		Sender: sender,
		Auth:   sender.getAuth(),
	}
}

func (s *MailServer) SendMail(mail *Mail) error {
	return smtp.SendMail(
		s.Sender.getHost()+":587",
		s.Auth,
		s.Sender.Email,
		mail.To,
		[]byte(mail.GetStr()),
	)
}