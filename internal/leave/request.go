package leave

import "time"

type CreateLeaveRequest struct {
	LeaveDate   string    `bson:"leave_date" json:"leave_date"`
	UserID      string    `bson:"user_id" json:"user_id"`
	Reason      *string   `bson:"reason" json:"reason"`
	RequestedAt time.Time `bson:"requested_at" json:"requested_at"`
}

type DeleteLeaveRequest struct {
	LeaveDate string `bson:"leave_date" json:"leave_date"`
	UserID    string `bson:"user_id" json:"user_id"`
}

type SettingRequest struct {
	MaxEmployeesPerDay int `bson:"max_employees_per_day" json:"max_employees_per_day"`
	AdvanceBookingDays int `bson:"advance_booking_days" json:"advance_booking_days"`
}

type EditSlotRequest struct {
	MaxSlot int `bson:"max_slot" json:"max_slot"`
}

type UpdateRequest struct {
	Types string `bson:"types" json:"types"`
}

