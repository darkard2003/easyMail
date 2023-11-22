package easyMail

type Mail struct {
	From    string
	To      []string
	Subject string
	Body    string
	IsHtml  bool
}

func NewMail(
	from string,
	to []string,
	subject string,
	body string,
	isHtml bool,
) *Mail {
	return &Mail{}
}

func (m *Mail) GetStr() string {
	if m.IsHtml {
		return "From: " + m.From + "\r\n" +
			"To: " + m.To[0] + "\r\n" +
			"Subject: " + m.Subject + "\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n\r\n" +
			m.Body
	} else {
		return "From: " + m.From + "\r\n" +
			"To: " + m.To[0] + "\r\n" +
			"Subject: " + m.Subject + "\r\n\r\n" +
			m.Body
	}
}
