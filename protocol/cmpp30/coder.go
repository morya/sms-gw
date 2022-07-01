package cmpp30

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/binary"

	"fmt"
	"io"
	"time"

	"github.com/morya/sms-gw/protocol/unified"
)

// DEBUG debug switch
var DEBUG bool = false

type Coder struct {
}

func (c *Coder) Decode(data []byte) (v interface{}, err error) {
	r := bytes.NewBuffer(data)
	var head = &CmppMsgHead{}
	err = binary.Read(r, binary.BigEndian, head)
	if err != nil {
		return
	}

	if head.Length > MAX_MSG_LENGTH {
		err = fmt.Errorf("msg too long, length=%v", head.Length)
		return
	}

	switch head.CmdID {
	case CMPP_CMD_BindReq:
		v, err = c.decodeBindReq(data)

	case CMPP_CMD_BindAck:
		v, err = c.decodeBindAck(data)

	case CMPP_CMD_ActiveTestAck:
		v, err = c.decodeHeartBeatAck(data)

	case CMPP_CMD_SubmitReq:
		v, err = c.decodeSubmitReq(data)

	default:
		err = fmt.Errorf("cmd not supported %v", head.CmdID)
	}

	return
}

func (c *Coder) decodeBindReq(data []byte) (v interface{}, err error) {
	var dstMsg = &unified.MsgBindReq{}
	var cmppMsg = &CmppMsgBindReq{}
	r := bytes.NewBuffer(data)

	binary.Read(r, binary.BigEndian, cmppMsg)
	dstMsg.Length = cmppMsg.Length // just record this, useless
	dstMsg.CmdID = unified.CMD_BIND_REQ
	dstMsg.Seq = cmppMsg.Seq

	return dstMsg, err
}

func (c *Coder) decodeBindAck(data []byte) (v interface{}, err error) {
	var dstMsg = &unified.MsgBindAck{}
	var cmppMsg = &CmppMsgBindAck{}
	r := bytes.NewBuffer(data)

	binary.Read(r, binary.BigEndian, cmppMsg)
	dstMsg.Length = cmppMsg.Length // just record this, useless
	dstMsg.CmdID = unified.CMD_BIND_ACK
	dstMsg.Seq = cmppMsg.Seq
	dstMsg.Succ = true

	return dstMsg, err
}

func (c *Coder) decodeHeartBeatAck(data []byte) (v interface{}, err error) {
	var dstMsg = &unified.MsgHeartBeat{}
	var cmppMsg = &CmppMsgHeartBeatAck{}
	r := bytes.NewBuffer(data)

	binary.Read(r, binary.BigEndian, cmppMsg)
	dstMsg.Length = cmppMsg.Length // just record this, useless
	dstMsg.CmdID = unified.CMD_HEARTBEAT_ACK
	dstMsg.Seq = cmppMsg.Seq

	return dstMsg, err
}
func (c *Coder) decodeSubmitReq(data []byte) (v interface{}, err error) {
	var dstMsg = &unified.MsgSubmitReq{}
	var cmppFront = &CmppMsgSubmitFrontReq{}
	r := bytes.NewBuffer(data)

	binary.Read(r, binary.BigEndian, cmppFront)
	dstMsg.Length = cmppFront.Length // just record this, useless
	dstMsg.CmdID = unified.CMD_SUBMIT_REQ
	dstMsg.Seq = cmppFront.Seq

	dstMsg.SrcAddr = string(cmppFront.SrcID[:])
	dstMsg.FeeAddr = string(cmppFront.FeeTerminalID[:])
	// dstMsg.DstAddr =  string(cmppFront.DestUserCount)

	return dstMsg, err
}

func (c *Coder) Encode(m interface{}) ([]byte, error) {
	switch msg := m.(type) {
	case *unified.MsgBindReq:
		// log.Println("marshal bind req ", msg)
		return c.encodeBindReq(msg)

	case *unified.MsgHeartBeat:
		return c.encodeHeartBeat(msg)

	case *unified.MsgSubmitReq:
		// log.Println("submit", msg)
		return c.encodeSubmit(msg)
	}

	return []byte{}, fmt.Errorf("not supported msg %v", m)
}

func getAuthChecksum(account, pswd string, timestamp []byte, checksum []byte) {

	var t = time.Now()

	var tsInt int
	tsInt += 100000000 * int(t.Month())
	tsInt += 1000000 * t.Day()
	tsInt += 10000 * t.Hour()
	tsInt += 100 * t.Minute()
	tsInt += t.Second()

	var tsUint = uint32(tsInt)
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, tsUint)

	copy(timestamp, b.Bytes())

	h := md5.New()
	w := bufio.NewWriter(h)
	w.WriteString(account)

	// 9 byte of zero
	var zero byte
	for i := 0; i < 9; i++ {
		w.WriteByte(zero)
	}

	w.WriteString(pswd)
	var tsString = fmt.Sprintf("%02d%02d%02d%02d%02d", t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	w.WriteString(tsString)
	w.Flush()

	copy(checksum, h.Sum(nil))
}

func (c *Coder) encodeBindReq(msg *unified.MsgBindReq) ([]byte, error) {
	var cmpp = &CmppMsgBindReq{}
	cmpp.CmdID = CMPP_CMD_BindReq
	cmpp.Seq = msg.Seq
	copy(cmpp.SourceAddr[:], []byte(msg.Account))
	var timestamp = make([]byte, 4)
	var checksum = make([]byte, 16)
	getAuthChecksum(msg.Account, msg.Password, timestamp, checksum)

	copy(cmpp.Authenticator[:], checksum)
	copy(cmpp.Timestamp[:], timestamp)

	cmpp.Version = 0x30
	cmpp.Length = uint32(binary.Size(cmpp))

	buff := new(bytes.Buffer)
	binary.Write(buff, binary.BigEndian, cmpp)

	return buff.Bytes(), nil
}

func (c *Coder) encodeHeartBeat(msg *unified.MsgHeartBeat) ([]byte, error) {
	var hb = &CmppMsgHead{}
	hb.CmdID = CMPP_CMD_ActiveTestReq
	hb.Seq = 1
	hb.Length = uint32(binary.Size(hb))

	buff := new(bytes.Buffer)
	binary.Write(buff, binary.BigEndian, hb)

	return buff.Bytes(), nil
}

func (c *Coder) encodeSubmit(msg *unified.MsgSubmitReq) ([]byte, error) {
	var cmpp = &CmppMsgSubmitFrontReq{}
	cmpp.CmdID = CMPP_CMD_SubmitReq
	cmpp.Seq = msg.Seq

	// first, we write things here
	// then copy all things togather, so all length can be easily calculated
	var middleBuffer = new(bytes.Buffer)

	cmpp.PkTotal = 1
	cmpp.PkNumber = 1

	cmpp.RegisteredDeliver = 0 // not a report
	cmpp.MsgLevel = 0
	copy(cmpp.ServiceID[:], []byte("ABC"))

	binary.Write(middleBuffer, binary.BigEndian, cmpp)

	h := new(bytes.Buffer)
	binary.Write(h, binary.BigEndian, uint32(len(middleBuffer.Bytes())))

	result := middleBuffer.Bytes()
	copy(result, h.Bytes())
	return result, nil
}

func (c *Coder) NextMsg(r io.Reader) ([]byte, error) {
	var length uint32
	err := binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}

	if length > MAX_MSG_LENGTH {
		return nil, fmt.Errorf("length too big %d", length)
	}

	var b = new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, length)

	var buff = make([]byte, int(length))
	r.Read(buff)

	b.Write(buff)

	return b.Bytes(), nil
}
