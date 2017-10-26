package cmpp30

type CmppMsgHead struct {
	Length uint32
	CmdID  uint32
	Seq    uint32
}

type CmppMsgBindReq struct {
	CmppMsgHead

	SourceAddr    [6]byte
	Authenticator [16]byte
	Version       uint8
	Timestamp     [4]byte
}

type CmppMsgBindAck struct {
	CmppMsgHead

	/*
		0: 正确
		1: 消息结构错
		2: 非法源地址
		3: 认证错
		4: 版本太高
		5: 其它错误
	*/
	Status            uint32
	AuthenticatorISMG [16]byte
	Version           uint8 // 对于cmpp3.0 高位4bit是0x3 低位4bit是0x0
}

// submit length vary according to usercount and msg-content length
// we could only build FIXed length front part
type CmppMsgSubmitFrontReq struct {
	CmppMsgHead

	MsgID    [8]byte
	PkTotal  uint8
	PkNumber uint8

	RegisteredDeliver uint8
	MsgLevel          uint8
	ServiceID         [10]byte
	FeeUserType       uint8
	FeeTerminalID     [32]byte
	FeeTerminalType   uint8
	TP_Pid            uint8
	TP_udhi           uint8
	/*
	   0:ascii
	   3:sms write SIM card
	   4:binary
	   8:UCS2
	   15:GBK
	*/
	MsgFmt uint8
	MsgSrc [6]byte
	// 01:free to feeUser
	// 02:charge feeUser by count
	// 03:charge feeUser by month
	FeeType       [2]byte
	FeeCode       [6]byte
	ValidTime     [17]byte
	AtTime        [17]byte
	SrcID         [21]byte
	DestUserCount uint8 // max 100
}

type CmppMsgSubmitAck struct {
	MsgID [8]byte
}

type CmppMsgMultiSubmitReq struct {
}

type CmppMsgHeartBeatAck struct {
	CmppMsgHead
	Status byte
}
