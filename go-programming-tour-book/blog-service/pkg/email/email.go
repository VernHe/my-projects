package email

import (
	"crypto/tls"
	"gopkg.in/mail.v2"
)

type Email struct {
	*SMTPInfo
}

type SMTPInfo struct {
	Host     string
	Port     int
	IsSSL    bool
	UserName string
	Password string
	From     string
}

func NewEmail(info *SMTPInfo) *Email {
	return &Email{SMTPInfo: info}
}

// SendMail 发送邮件
func (e *Email) SendMail(to []string, subject, body string) error {
	// 要发送的内容
	message := mail.NewMessage()
	message.SetHeader("From", e.From)     // 发件人
	message.SetHeader("To", to...)        // 收件人
	message.SetHeader("Subject", subject) // 邮件主题
	message.SetBody("text/html", body)    // 邮件内容
	// 创建拨号器
	dialer := mail.NewDialer(e.Host, e.Port, e.UserName, e.Password)
	// 配置
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: e.IsSSL}
	// 发送邮件
	return dialer.DialAndSend(message)
}
