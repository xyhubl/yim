package responsity

import "go.mongodb.org/mongo-driver/bson/primitive"

type MemberMessage struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Seq        int                `bson:"seq" json:"seq"`
	Ack        int                `bson:"ack" json:"ack"`
	SendMid    string             `bson:"send_mid" json:"send_mid"`
	ReceiveMid string             `bson:"receive_mid" json:"receive_mid"`
	Message    string             `bson:"message" json:"message"`
	Base
}
