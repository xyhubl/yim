package protocol

const (
	OpHeartbeat      = int32(2)
	OpHeartbeatReply = int32(3)

	OpAuth      = int32(7)
	OpAuthReply = int32(8)

	OpProtoReady  = int32(10)
	OpProtoFinish = int32(11)

	OpChangeRoom      = int32(12)
	OpChangeRoomReply = int32(13)

	OpSub      = int32(14)
	OpSubReply = int32(15)

	OpUnsub      = int32(16)
	OpUnsubReply = int32(17)
)
