package protocol

const (
	OpHeartbeat      = int32(2)
	OpHeartbeatReply = int32(3)

	OpAuth      = int32(7)
	OpAuthReply = int32(8)

	OpProtoReady  = int32(10)
	OpProtoFinish = int32(11)
)
