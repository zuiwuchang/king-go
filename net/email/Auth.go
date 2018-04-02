package email

import (
	"errors"
	"net/smtp"
)

var errUnexpectedServerChallenge = errors.New("unexpected server challenge")

// InsecureAuth 非加密的 email 驗證
type InsecureAuth struct {
	User, Pwd string
}

// Start .
func (a *InsecureAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	resp := []byte("\x00" + a.User + "\x00" + a.Pwd)
	return "PLAIN", resp, nil
}

// Next .
func (*InsecureAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		// We've already sent everything.
		return nil, errUnexpectedServerChallenge
	}
	return nil, nil
}
