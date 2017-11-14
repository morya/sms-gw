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
	TPPid             uint8
	TPudhi            uint8
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

type CmppMsgDeliverReq struct {
	CmppMsgHead

	MsgID uint64

	// 目的号码
	DestID [21]byte

	// 业务标识 数字、字母和符号的组合
	ServiceID [10]byte
	TPPid     uint8
	TPudhi    uint8

	/*
	 * 0:ascii
	 * 3: write Sim-sms
	 * 4: binary sms
	 * 8: UCS2 coded msg
	 * 15: GBK coded msg
	 */
	MsgFmt             uint8
	SrcTerminalID      [32]byte
	SrcTerminalType    uint8
	RegisteredDelivery uint8 // 是否为状态报告 0:no, 1:yes
	MsgLength          uint8
	// 后续字段是变长的
	// MsgContent []byte
	// LinkID  [20]byte
}

type CmppMsgInnerReport struct {
	MsgID uint64
	Stat  [7]byte
}

type CmppMsgHeartBeatAck struct {
	CmppMsgHead
	Status         [7]byte
	SubmitTime     [10]byte
	DoneTime       [10]byte
	DestTerminalID [21]byte
	SMSCSequence   uint32
}

const (
	CMPP_STAT_Delivered      = "DELIVRD"
	CMPP_STAT_Expired        = "EXPIRED"
	CMPP_STAT_Deleted        = "DELETED"
	CMPP_STAT_DUndeliverable = "UNDELIV"
	CMPP_STAT_Accepted       = "ACCEPTD"
	CMPP_STAT_Unknown        = "UNKNOWN"
	CMPP_STAT_Rejected       = "REJECTED"

	// smsc 不返回响应消息时的状态报告 "MA:xxxx"
	CMPP_STAT_SmscNoAck = `MA:\w{4}`
	// smsc 返回错误应答时的状态报告 "MB:xxxx" xxxx是错误码
	CMPP_STAT_SmscAckErr = `MB:\w{4}`

	// SCP 不返回响应消息的状态报告
	CMPP_STAT_ScpNoAck = `CA:\w{4}`
	// SCP 返回错误应答时的状态报告
	CMPP_STAT_ScpAckErr = `CB:\w{4}`
)

type CmppMsgDeliverAck struct {
	MsgID uint64

	Result uint32
}
