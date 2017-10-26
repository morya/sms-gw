package cmpp30

const (
	CMPP_CMD_BindReq       uint32 = 0x00000001
	CMPP_CMD_UnbindReq     uint32 = 0x00000002
	CMPP_CMD_SubmitReq     uint32 = 0x00000004
	CMPP_CMD_DeliverReq    uint32 = 0x00000003
	CMPP_CMD_ActiveTestReq uint32 = 0x00000008
	CMPP_CMD_BindAck       uint32 = 0x80000001
	CMPP_CMD_UnbindAck     uint32 = 0x80000002
	CMPP_CMD_SubmitAck     uint32 = 0x80000004
	CMPP_CMD_DeliverAck    uint32 = 0x80000003
	CMPP_CMD_ActiveTestAck uint32 = 0x80000008
)

const (
	MAX_MSG_LENGTH = 10240
)
