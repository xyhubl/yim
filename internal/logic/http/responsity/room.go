package responsity

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	// 正常
	RoomNormalStatus = 1
	// 禁用
	RoomUnNormalStatus = 2
)

type Room struct {
	Id               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Typ              int                `bson:"typ" json:"typ"`
	OwnerMid         int                `bson:"owner_mid" json:"owner_mid"` // 群主id
	RoomName         string             `bson:"room_name" json:"room_name"`
	RoomCapacity     string             `bson:"room_capacity" json:"room_capacity"`         // 房间容量 默认500
	RoomAnnouncement string             `bson:"room_announcement" json:"room_announcement"` // 公告
	RoomStatus       int                `bson:"room_status" json:"room_status"`
	Base
}
