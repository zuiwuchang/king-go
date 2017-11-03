package pool

import (
	"encoding/binary"
	"errors"
	"github.com/zuiwuchang/king-go/net/echo"
	"net"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	laddr := "127.0.0.1:60000"
	__runTestServer(t, laddr)
	__runTestPool0(t, laddr)
}
func __runTestPool0(t *testing.T, laddr string) {

	p, e := NewPool(
		&__testPoolTemplate{NewPoolTemplate(laddr, time.Hour, 2, 10), laddr},
		0,
	)
	if e != nil {
		t.Fatal(e)
	}
	defer p.Close()

	for i := 0; i < 10; i++ {
		var c *Conn
		c, e = p.Get()
		if e != nil {
			t.Fatal(e)
		}
		defer p.Put(c)

		e = testWrite(c, __TestCmdEcho, nil)
		if e != nil {
			c.Close()
			t.Fatal(e)
		}

		c0 := c.Get().(echo.IClient)
		if b, e := c0.GetMessage(0); e != nil {
			c.Close()
			t.Fatal(e)
		} else {
			m := __testMsg{b: b}
			if m.Cmd() != __TestCmdEcho {
				c.Close()
				t.Fatal("bad echo")
			}
		}
	}

	c0, e := p.Get()
	if e != nil {
		c0.Close()
		t.Fatal(e)
	}
	p.Put(c0)
	c1, e := p.Get()
	if e != nil {
		c1.Close()
		t.Fatal(e)
	}
	p.Put(c1)
	if c0 != c1 {
		t.Fatal("bad get put")
	}
}
func testWrite(c net.Conn, cmd uint16, data []byte) error {
	n := 4 + len(data)
	b := make([]byte, n)
	binary.LittleEndian.PutUint16(b, uint16(n))
	binary.LittleEndian.PutUint16(b[2:], cmd)
	copy(b[4:], data)
	_, e := c.Write(b)
	return e
}
func __runTestServer(t *testing.T, addr string) {
	s, e := echo.NewServer(addr, 0, &__testServerTemplate{})
	if e != nil {
		t.Fatal(e)
	}
	s.Run()
}

const (
	__TestCmdLogin = 1
	__TestCmdEcho  = 2
)

type __testServerSession struct {
	Ok bool
}
type __testServerTemplate struct {
}

func (s *__testServerTemplate) GetHeaderSize() int {
	return 4
}

func (s *__testServerTemplate) GetMessageSize(session echo.Session, header []byte) (n int, e error) {
	size := binary.LittleEndian.Uint16(header)
	if size < 4 {
		return 0, errors.New("bad size")
	}
	return int(size), nil
}

func (s *__testServerTemplate) NewSession(c net.Conn) (session echo.Session, e error) {
	return &__testServerSession{}, nil
}
func (s *__testServerTemplate) DeleteSession(c net.Conn, session echo.Session) {

}
func (s *__testServerTemplate) Message(c net.Conn, session echo.Session, b []byte) error {
	s0 := session.(*__testServerSession)
	cmd := binary.LittleEndian.Uint16(b[2:])
	if cmd == __TestCmdLogin {
		s0.Ok = true
	} else {
		if !s0.Ok {
			return errors.New("not login")
		}
	}

	_, e := c.Write(b)
	return e
}

type __testClientTemplate struct {
}

func (c *__testClientTemplate) GetHeaderSize() int {
	return 4
}
func (c *__testClientTemplate) GetMessageSize(header []byte) (n int, e error) {
	size := binary.LittleEndian.Uint16(header)
	if size < 4 {
		return 0, errors.New("bad size")
	}
	return int(size), nil
}

type __testPoolTemplate struct {
	*PoolTemplate
	laddr string
}
type __testMsg struct {
	b []byte
}

func (m *__testMsg) Cmd() uint16 {
	return binary.LittleEndian.Uint16(m.b[2:])
}

func (p *__testPoolTemplate) Conect() (net.Conn, error) {
	c, e := echo.NewClient(p.laddr, &__testClientTemplate{})
	if e != nil {
		return nil, e
	}
	e = testWrite(c, __TestCmdLogin, nil)
	if e != nil {
		return nil, e
	}
	if b, e := c.GetMessage(0); e != nil {
		return nil, e
	} else {
		m := &__testMsg{b: b}
		if m.Cmd() != __TestCmdLogin {
			return nil, errors.New("bad login")
		}
	}
	return c, nil
}
