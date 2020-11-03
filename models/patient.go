package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Patient struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" `
	Name     string             `json:"name" form:"name" bson:"name"`
	Surname  string             `json:"surname" form:"surname" bson:"surname"`
	Fullname string             `json:"fullname" form:"fullname" bson:"fullname"`
	DNI      string             `json:"dni" form:"dni" bson:"dni"`
	Email    string             `json:"email" form:"email" bson:"email"`
	Address  string             `json:"address" form:"address" bson:"address"`
	Phone    string             `json:"phone" form:"phone" bson:"phone"`
	//Invoices []string           `json:"invoices" form:"invoices" bson:"invoices,omitempty"`
	IsOrg bool `json:"isorg" form:"isorg" bson:"isorg,omitempty"`
}
