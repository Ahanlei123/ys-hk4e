package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ChatMsgTypeText = iota
	ChatMsgTypeIcon
)

type ChatMsg struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Sequence uint32             `bson:"-"`
	Time     uint32             `bson:"Time"`
	Uid      uint32             `bson:"Uid"`
	ToUid    uint32             `bson:"ToUid"`
	IsRead   bool               `bson:"IsRead"`
	MsgType  uint8              `bson:"MsgType"`
	Text     string             `bson:"Text"`
	Icon     uint32             `bson:"Icon"`
}
