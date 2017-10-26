package protocol

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"

	_ "github.com/morya/sms/protocol/cmpp30" // plugin like inner package
	_ "github.com/morya/sms/protocol/sgip12" // plugin like inner package
	_ "github.com/morya/sms/protocol/smpp34" // plugin like inner package

	"github.com/morya/utils/log"
	"github.com/morya/sms/protocol/unified"
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

	log.Info("encode result\n", hex.Dump(data))
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

	log.Info("encode result\n", hex.Dump(data))
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
	log.Infof("decode msg %v, data\n%s", ack, hex.Dump(data))
	return
}
