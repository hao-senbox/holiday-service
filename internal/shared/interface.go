package shared

import (
	"context"
	"time"
)

type AttendanceRepository interface {
	GetStudentAttendanceInDateRange(c context.Context, studentID string, startDate time.Time, endDate time.Time) ([]*AttendanceStudent, error)
}
