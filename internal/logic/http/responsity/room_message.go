package responsity

import "go.mongodb.org/mongo-driver/bson/primitive"

type RoomMessage struct {
	Id      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Seq     int                `bson:"seq" json:"seq"`
	Ack     int                `bson:"ack" json:"ack"`
	Mid     string             `bson:"mid" json:"mid"`
	RoomId  string             `bson:"room_id" json:"room_id"`
	Message string             `bson:"message" json:"message"`
	Base
}
