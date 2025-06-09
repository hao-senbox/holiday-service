package attendance

import "worktime-service/internal/user"

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
