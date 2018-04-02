package email

import (
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/smtp"
	"strings"
)

const (
	// TypeHTML .
	TypeHTML = "html"
)

// ErrWriteBusy tcp 發送緩衝區 已滿
var ErrWriteBusy = errors.New("write busy")

func writeMsg(w io.Writer,
	user, to,
	subject, body,
	mailtype string) (e error) {

	// msg
	msg := "To: " + to +
		"\r\nFrom: " + user + "<" + user +
		">\r\nSubject: " + subject + "\r\n" +
		"Content-Type: text/" + mailtype + "; charset=UTF-8\r\n\r\n" + body

	// write
	var n int
	n, e = w.Write([]byte(msg))
	if e != nil {
		return
	} else if n != len(msg) {
		e = ErrWriteBusy
	}
	return
}

func sendEmail(c *smtp.Client,
	user, to,
	subject, body,
	mailtype string) (e error) {

	// 設置 發信地址
	if e = c.Mail(user); e != nil {
		return
	}

	//設置 收信 地址
	sendToS := strings.Split(to, ";")
	for _, sendTo := range sendToS {
		if e = c.Rcpt(sendTo); e != nil {
			return
		}
	}

	// 創建 writer
	var w io.WriteCloser
	w, e = c.Data()
	if e != nil {
		return
	}

	// 寫入 發送數據
	e = writeMsg(w,
		user, to,
		subject, body,
		mailtype,
	)
	// 發送 結束
	e = w.Close()
	return
}

// SendEmail 使用非 加密 發送 email
func SendEmail(host,
	user, pwd, to,
	subject, body,
	mailtype string,
) (e error) {
	// 創建 客戶端
	var c *smtp.Client
	c, e = NewSMTPClient(host, user, pwd)
	if e != nil {
		return
	}

	// 發送 email
	e = sendEmail(c,
		user, to,
		subject, body,
		mailtype,
	)
	if e == nil {
		c.Quit()
	} else {
		c.Close()
	}
	return
}

// SendSSLEmail 使用 ssl 發送 email
func SendSSLEmail(host,
	user, pwd, to,
	subject, body,
	mailtype string,
) (e error) {
	// 創建 客戶端
	var c *smtp.Client
	c, e = NewSMTPSSLClient(host, user, pwd)
	if e != nil {
		return
	}

	// 發送 email
	e = sendEmail(c,
		user, to,
		subject, body,
		mailtype,
	)
	if e == nil {
		c.Quit()
	} else {
		c.Close()
	}
	return
}

// NewSMTPClient 創建一個 未加密 的 smtp 客戶端
func NewSMTPClient(host, user, pwd string) (c *smtp.Client, e error) {
	// 創建 smtp 客戶端
	c, e = smtp.Dial(host)
	if e != nil { // 失敗會自動 tc.Close
		return
	}

	// 驗證
	e = c.Auth(&InsecureAuth{
		User: user,
		Pwd:  pwd,
	})
	if e != nil {
		c = nil
	}
	return
}

// NewSMTPSSLClient 創建一個 使用 tls 的 smtp 客戶端
func NewSMTPSSLClient(host, user, pwd string) (c *smtp.Client, e error) {
	// 建立 tls 連接
	var serverName string
	serverName, _, e = net.SplitHostPort(host)
	if e != nil {
		return
	}

	var tc *tls.Conn
	tc, e = tls.Dial(
		"tcp",
		host,
		&tls.Config{
			ServerName: serverName,
			//InsecureSkipVerify: true, // 如果為 true 不驗證 證書
		},
	)
	if e != nil {
		return
	}

	// 創建 smtp 客戶端
	c, e = smtp.NewClient(tc, serverName)
	if e != nil { // 失敗會自動 tc.Close
		return
	}

	// 驗證
	e = c.Auth(smtp.PlainAuth("", user, pwd, serverName))
	if e != nil { // 失敗會自動 Close
		c = nil
	}
	return
}
