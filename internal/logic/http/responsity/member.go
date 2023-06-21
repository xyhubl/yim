package responsity

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// 正常
	MemberNormalStatus = 1
	// 禁用
	MemberUnNormalStatus = 2
)

// 成员
type Member struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Nickname     string             `bson:"nickname" json:"nickname"`
	Password     string             `bson:"password" json:"password"`
	Email        string             `bson:"email" json:"email"`
	Phone        string             `bson:"phone" json:"phone"`
	MemberStatus int                `bson:"member_status" json:"member_status"`
	Base
}

type MemberModel struct {
	tx       *mongo.Client
	database string
}

func (m *MemberModel) Collection() *mongo.Collection {
	return m.tx.Database(m.database).Collection("member")

}

func (m *MemberModel) InsertOne(ctx context.Context, member *Member) (*mongo.InsertOneResult, error) {
	return m.Collection().InsertOne(ctx, member)
}
