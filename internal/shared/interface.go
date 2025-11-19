package shared

import (
	"context"
	"time"
)

// AttendanceRepository defines methods needed to fetch student attendance
// data for shared use-cases (e.g., temperature chart).
type AttendanceRepository interface {
	GetStudentAttendanceInDateRange(c context.Context, studentID string, startDate time.Time, endDate time.Time) ([]*AttendanceStudent, error)
}
