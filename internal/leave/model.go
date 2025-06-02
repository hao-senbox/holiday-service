package leave

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Setting struct {
	ID                 primitive.ObjectID `bson:"_id" json:"id"`
	MaxEmployeesPerDay int                `bson:"max_employees_per_day" json:"max_employees_per_day"`
	AdvanceBookingDays int                `bson:"advance_booking_days" json:"advance_booking_days"`
	CreatedAt          time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt          time.Time          `bson:"updated_at" json:"updated_at"`
}

type DailyLeaveSolt struct {
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	Date            time.Time          `bson:"date" json:"date"`
	MaxSlot         int                `bson:"max_slot" json:"max_slot"`
	AvailableSlot   int                `bson:"available_slot" json:"available_slot"`
	ConfirmedLeaves []ConfirmedLeave   `bson:"confirmed_leaves" json:"confirmed_leaves"`
	PendingRequests []PendingRequest   `bson:"pending_requests" json:"pending_requests"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

type ConfirmedLeave struct {
	UserID    string    `bson:"user_id" json:"user_id"`
	UserName  string    `bson:"user_name" json:"user_name"`
	ApproveAt time.Time `bson:"approve_at" json:"approve_at"`
}

type PendingRequest struct {
	LeaveID   primitive.ObjectID `bson:"leave_id" json:"leave_id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	UserName  string             `bson:"user_name" json:"user_name"`
	Status    string             `bson:"status" json:"status"`
	RequestAt time.Time          `bson:"request_at" json:"request_at"`
}

type LeaveRequests struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	LeaveDate   time.Time          `bson:"leave_date" json:"leave_date"`
	UserID      string             `bson:"user_id" json:"user_id"`
	RequestType string             `bson:"request_type" json:"request_type"`
	UserName    string             `bson:"user_name" json:"user_name"`
	Reason      *string            `bson:"reason" json:"reason"`
	RequestedAt time.Time          `bson:"requested_at" json:"requested_at"`
	Status      string             `bson:"status" json:"status"`
}
