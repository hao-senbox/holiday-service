package shared

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AttendanceStudent struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	UserID       string             `json:"user_id" bson:"user_id"`
	Temperature  string             `json:"temperature" bson:"temperature"`
	DayOfWeek    time.Weekday       `json:"day_of_week" bson:"day_of_week"`
	Date         time.Time          `json:"date" bson:"date"`
	CheckInTime  *time.Time         `json:"check_in_time" bson:"check_in_time"`
	CheckOutTime *time.Time         `json:"check_out_time" bson:"check_out_time"`
	Types        string             `json:"types" bson:"types"`
	Note         string             `json:"note" bson:"note"`
	CreatedBy    string             `json:"created_by" bson:"created_by"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}
