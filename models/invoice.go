package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" `
	Patient        string             `json:"patient" form:"patient" bson:"patient,omitempty"`
	PatientID      string             `json:"patientid" form:"patientid" bson:"patientid,omitempty"`
	InvoiceNumber  int64              `json:"invoice_number" form:"invoice_number" bson:"invoice_number,omitempty"`
	Description    string             `json:"desc" form:"desc" bson:"desc,"`
	Units          int                `json:"units" form:"units" bson:"units,"`
	Price          float64            `json:"price" form:"price" bson:"price"`
	Retention      int                `json:"retention" form:"retention" bson:"retention,"`
	Payed          float64            `json:"payed" form:"payed" bson:"payed,"`
	Date           time.Time          `json:"date" form:"date" bson:"date,omitempty"`
	PatientDNI     string             `json:"patientdni" form:"patientdni" bson:"patientdni"`             //not omit, cause when moving invoice to a patient without invoice data, those fields has to be emptied
	PatientAddress string             `json:"patientaddress" form:"patientaddress" bson:"patientaddress"` //not omit, same reason
	IsOrg          bool               `json:"isorg" form:"isorg" bson:"isorg"`                            //nede for updateInvoice, when changing invoice between org/patient
}
