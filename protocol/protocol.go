package protocol

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"

	_ "github.com/morya/sms-gw/protocol/cmpp30" // plugin like inner package
	_ "github.com/morya/sms-gw/protocol/sgip12" // plugin like inner package
	_ "github.com/morya/sms-gw/protocol/smpp34" // plugin like inner package

	"github.com/morya/sms-gw/protocol/unified"
	"github.com/sirupsen/logrus"
)

func Test(account, pswd, remote string) (err error) {
	var coder = unified.GetCoder("cmpp")
	if coder == nil {
		err = fmt.Errorf("no coder named %v", "cmpp")
		return
	}

	var m = &unified.MsgBindReq{Account: account, Password: pswd}
	m.Seq = 15
	data, err := coder.Encode(m)
	if err != nil {
		return
	}

	logrus.Debug("encode result", hex.Dump(data))
	return
}

func Run(account, pswd, remote string) (err error) {
	var coder = unified.GetCoder("cmpp")
	if coder == nil {
		err = fmt.Errorf("no coder named %v", "cmpp")
		return
	}

	var m = &unified.MsgBindReq{Account: account, Password: pswd}
	m.Seq = 15
	data, err := coder.Encode(m)
	if err != nil {
		return
	}

	logrus.Debug("encode result", hex.Dump(data))
	c, err := net.Dial("tcp", remote)
	if err != nil {
		return
	}

	defer c.Close()
	r := bufio.NewReader(c)
	w := bufio.NewWriter(c)
	w.Write(data)
	w.Flush()
	data, err = coder.NextMsg(r)
	if err != nil {
		return
	}

	ack, err := coder.Decode(data)
	logrus.Debugf("decode msg %v, data\n%s", ack, hex.Dump(data))
	return
}
