package unified

type MsgHead struct {
	Length uint32
	CmdID  uint32
	Seq    uint32
	Status uint32

	RecvConnID uint32
	SendConnID uint32
}

type MsgBindReq struct {
	MsgHead

	Account  string
	Password string
}

type MsgBindAck struct {
	MsgHead
	Succ bool
}

type MsgHeartBeat struct {
	MsgHead
}

type MsgSubmitReq struct {
	MsgHead

	SrcAddr string
	DstAddr string
	FeeAddr string

	LinkID string
}

type MsgSubmitAck struct {
	MsgHead
	MsgID [8]byte // used in submit
}

type MsgMultiSubmitReq struct {
	MsgHead
}

type MsgDeliverReq struct {
	MsgHead
	MsgID [8]byte
}

type MsgDeliverAck struct {
	MsgHead
	MsgID  [8]byte
	Result uint32
}
