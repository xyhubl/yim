package responsity

const (
	// AckServerRcv 服务端ack确认收到消息
	AckServerRcv = 1
	// AckServerSndToClient 服务端已发送至接收者
	AckServerSndToClient = 2
	// AckServerSndToClientConfirm 接收者ack确认完成
	AckServerSndToClientConfirm = 3
)

type Base struct {
	CreatedAt int `bson:"created_at" json:"created_At"`
	UpdateAt  int `bson:"update_at" json:"update_at"`
}
