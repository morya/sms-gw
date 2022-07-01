package session

import (
	"bufio"
	"fmt"
	"net"

	options "github.com/morya/sms-gw/options"
	_ "github.com/morya/sms-gw/protocol/cmpp30"
	"github.com/morya/sms-gw/protocol/unified"
)

type ConnMode int

const (
	CONN_MODE_CLIENT ConnMode = 0 + iota
	CONN_MODE_SERVER
)

type Conn struct {
	opt *options.Options

	lastErr error

	coder   unified.Coder
	sock    net.Conn
	tcpMode ConnMode

	writer *bufio.Writer
	reader *bufio.Reader

	chanMsgSend chan []byte
	chanMsgRecv chan interface{}
}

func NewConn(opt *options.Options, codername string) (*Conn, error) {

	sock, err := net.Dial("tcp", opt.RemoteAddr)
	if err != nil {
		return nil, err
	}

	var coder = unified.GetCoder(codername)
	if coder == nil {
		return nil, fmt.Errorf("coder not found name=[%s]", codername)
	}

	c := &Conn{
		opt:   opt,
		coder: coder,
		sock:  sock,

		writer: bufio.NewWriter(sock),
		reader: bufio.NewReader(sock),

		chanMsgSend: make(chan []byte, 16),
		chanMsgRecv: make(chan interface{}, 10),
	}
	return c, nil
}

func (c *Conn) Send(data []byte) {
	c.chanMsgSend <- data
}

func (c *Conn) ClientRun() error {
	c.tcpMode = CONN_MODE_CLIENT
	return c.Loop()
}
