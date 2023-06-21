package responsity

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	// 正常
	RoomMemberNormalStatus = 1
	// 禁用
	RoomMemberUnNormalStatus = 2
)

type RoomMember struct {
	Id               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Mid              string             `bson:"mid" json:"mid"`
	MidNickname      string             `bson:"mid_nickname" json:"mid_nickname"`
	RoomMemberStatus int                `bson:"room_member_status" json:"room_member_status"`
	Base
}
