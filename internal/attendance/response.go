package attendance

import (
	"time"
	"worktime-service/internal/user"
)

type DailyAttendanceResponse struct {
	Date              string  `json:"date"`
	DayOfWeek         string  `json:"day_of_week"`
	Status            string  `json:"status"`
	CheckInTime       string  `json:"check_in_time"`
	EmotionCheckIn    string  `json:"emotion_check_in"`
	CheckoutTime      string  `json:"checkout_time"`
	LunchDuration     int     `json:"lunch_duration"`
	EMotionCheckOut   string  `json:"emotion_check_out"`
	PercentWorkDay    float64 `json:"percent_work_day"`
	TotalWorkingHours float64 `json:"total_working_hours"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

type MonthlyAttendanceResponse struct {
	Employee          user.UserInfor            `json:"employee"`
	YearMonth         string                    `json:"year_month"`
	Summary           MonthlySummary            `json:"summary"`
	MonthlyAttendance []DailyAttendanceResponse `json:"monthly_attendance"`
}

type MonthlySummary struct {
	TotalWorkDays  int     `json:"total_work_days"`
	PresentDays    int     `json:"present_days"`
	LeaveDays      float64 `json:"leave_days"`
	AbsentDays     int     `json:"absent_days"`
	TotalWorkHours float64 `json:"total_work_hours"`
}

type DailyAttendanceResponsePagination struct {
	DailyAttendance []*DailyAttendanceUser `json:"daily_attendance"`
	Pagination      Pagination             `json:"pagination"`
}

type DailyAttendanceUser struct {
	ID                string          `json:"id" bson:"_id"`
	UserInfor         *user.UserInfor `json:"user_infor" bson:"user_infor"`
	DayOfWeek         time.Weekday    `json:"day_of_week" bson:"day_of_week"`
	Date              time.Time       `json:"date" bson:"date"`
	Status            string          `json:"status" bson:"status"`
	CheckInTime       string          `json:"check_in_time" bson:"check_in_time"`
	EmotionCheckIn    string          `json:"emotion_check_in" bson:"emotion_check_in"`
	CheckoutTime      string          `json:"check_out_time" bson:"check_out_time"`
	LunchDuration     int             `json:"lunch_duration" bson:"lunch_duration"`
	EMotionCheckOut   string          `json:"emotion_check_out" bson:"emotion_check_out"`
	PercentWorkDay    float64         `json:"percent_work_day" bson:"percent_work_day"`
	TotalWorkingHours float64         `json:"total_working_hours" bson:"total_working_hours"`
	CreatedAt         string          `json:"created_at" bson:"created_at"`
	UpdatedAt         string          `json:"updated_at" bson:"updated_at"`
}

type Pagination struct {
	TotalCount int64 `json:"total_count"`
	TotalPages int64 `json:"total_pages"`
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
}

type StudentTemperatureChartResponse struct {
	Title        string    `json:"title"`
	Unit         string    `json:"unit"`
	Labels       []string  `json:"labels"`
	Temperatures []float64 `json:"temperatures"`
}
