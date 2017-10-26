package protocol

import (
	"testing"

	"github.com/morya/utils/log"
	"github.com/morya/sms/protocol/unified"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func TestEncodeBindReq(t *testing.T) {
	var account, pswd string
	var coder = unified.GetCoder("cmpp")
	if coder == nil {
		t.Errorf("no coder named %v", "cmpp")
		return
	}

	var m = &unified.MsgBindReq{Account: account, Password: pswd}
	m.Seq = 15
	data, err := coder.Encode(m)
	if err != nil {
		t.Error(err)
		return
	}

	v, err := coder.Decode(data)
	if err != nil {
		t.Error(err)
	}

	newmsg, ok := v.(*unified.MsgBindReq)
	if !ok {
		t.Errorf("return msg with wrong type, %T", v)
	}
	if newmsg.Seq != 15 {
		t.Errorf("seq convert fail")
		return
	}

	if newmsg.Length == 0 {
		t.Errorf("length convert fail")
		return
	}
	return
}

func TestEncodeSubmitReq(t *testing.T) {
	var coder = unified.GetCoder("cmpp")
	if coder == nil {
		t.Errorf("no coder named %v", "cmpp")
		return
	}

	var m = &unified.MsgSubmitReq{}
	m.Seq = 15
	data, err := coder.Encode(m)
	if err != nil {
		t.Errorf("%v", err)
	}

	v, err := coder.Decode(data)
	if err != nil {
		t.Error(err)
	}

	newmsg, ok := v.(*unified.MsgSubmitReq)
	if !ok {
		t.Errorf("return msg with wrong type, %T", v)
	}
	if newmsg.Seq != 15 {
		t.Errorf("seq convert fail")
		return
	}

	if newmsg.Length == 0 {
		t.Errorf("length convert fail")
		return
	}
}
