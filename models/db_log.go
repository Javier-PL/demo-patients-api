package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DBLog struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" `
	ItemID     string			  `json:"itemid" form:"itemid" bson:"itemid,omitempty"`
	Action     string             `json:"action" form:"action" bson:"action,omitempty"`
	Result     string             `json:"result" form:"result" bson:"result,omitempty"`
	Object     string             `json:"object" form:"object" bson:"object,omitempty"`
	Itemjson   string             `json:"itemjson" form:"itemjson" bson:"itemjson,omitempty"`
	Updatejson string             `json:"updatejson" form:"updatejson" bson:"updatejson,omitempty"`
	Date       time.Time          `json:"date" form:"date" bson:"date,omitempty"`
	
}
