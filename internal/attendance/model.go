package attendance

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AttendanceLog struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserID    string             `json:"user_id" bson:"user_id"`
	LogDate   time.Time          `json:"log_date" bson:"log_date"`
	LogTime   time.Time          `json:"log_time" bson:"log_time"`
	LogType   string             `json:"log_type" bson:"log_type"`
	Emotion   string             `json:"emotion" bson:"emotion"`
	Notes     string             `json:"notes" bson:"notes"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type DailyAttendance struct {
	ID                primitive.ObjectID `json:"id" bson:"_id"`
	UserID            string             `json:"user_id" bson:"user_id"`
	DayOfWeek         time.Weekday       `json:"day_of_week" bson:"day_of_week"`
	Date              time.Time          `json:"date" bson:"date"`
	Status            string             `json:"status" bson:"status"`
	CheckInTime       *time.Time         `json:"check_in_time" bson:"check_in_time"`
	EmotionCheckIn    string             `json:"emotion_check_in" bson:"emotion_check_in"`
	CheckoutTime      *time.Time         `json:"check_out_time" bson:"check_out_time"`
	LunchDuration     int                `json:"lunch_duration" bson:"lunch_duration"`
	EMotionCheckOut   string             `json:"emotion_check_out" bson:"emotion_check_out"`
	PercentWorkDay    float64            `json:"percent_work_day" bson:"percent_work_day"`
	TotalWorkingHours float64            `json:"total_working_hours" bson:"total_working_hours"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
}
