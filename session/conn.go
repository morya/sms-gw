package session

import (
	"bufio"
	"net"
	"sync"

	"github.com/morya/utils/log"

	_ "github.com/morya/sms/protocol/cmpp30"
	"github.com/morya/sms/protocol/unified"
)

type ConnMode int

const (
	CONN_MODE_CLIENT ConnMode = 0 + iota
	CONN_MODE_SERVER
)

type Conn struct {
	wg      *sync.WaitGroup
	coder   unified.Coder
	sock    net.Conn
	tcpMode ConnMode

	Account    string
	Password   string
	RemoteAddr string

	writer *bufio.Writer
	reader *bufio.Reader

	chanMsgSend chan []byte
	chanMsgRecv chan interface{}
}

func NewConn(account, pswd, remote string) (*Conn, error) {
	log.Infof("sysid = %v, pswd = %v, remote= %v", account, pswd, remote)

	sock, err := net.Dial("tcp", remote)
	if err != nil {
		return nil, err
	}

	c := &Conn{
		coder:      unified.GetCoder("cmpp"),
		sock:       sock,
		Account:    account,
		Password:   pswd,
		RemoteAddr: remote,

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

func (c *Conn) Client() {
	c.tcpMode = CONN_MODE_CLIENT
	c.Loop()
}
