package job

import (
	"github.com/xyhubl/yim/api/comet"
	pb "github.com/xyhubl/yim/api/logic"
	"github.com/xyhubl/yim/api/protocol"
	"github.com/xyhubl/yim/pkg/bytes"
	"golang.org/x/net/context"
	"log"
)

func (j *Job) push(ctx context.Context, pushMsg *pb.PushMsg) error {
	var err error
	switch pushMsg.Type {
	case pb.Type_PUSH:
		j.pushKeys(pushMsg.Operation, pushMsg.Server, pushMsg.Keys, pushMsg.Msg)
	case pb.Type_ROOM:
		j.getRoom(pushMsg.Room).Push(pushMsg.Operation, pushMsg.Msg)
	}

	return err
}

// zh: 单聊
func (j *Job) pushKeys(op int32, serverID string, subKeys []string, body []byte) {
	// 64 字节预留协议扩展
	buf := bytes.NewWriterSize(len(body) + 64)
	p := &protocol.Proto{
		Ver:  1,
		Op:   op,
		Body: body,
	}
	p.WriteTo(buf)
	p.Body = buf.Buffer()
	p.Op = protocol.OpRaw

	req := &comet.PushMsgReq{
		Keys:    subKeys,
		ProtoOp: op,
		Proto:   p,
	}
	if c, ok := j.cometServers[serverID]; ok {
		c.Push(req)
		log.Println("[INFO] PUSH: pushServer: ", serverID, "comets: ", len(j.cometServers))
	}
}
