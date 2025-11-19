package shared

import (
	"context"
	"time"
	"worktime-service/internal/attendance"
)

type StudentTemperatureChartResponse interface {
	GetStudentAttendanceInDateRange(c context.Context, studentID string, startDate time.Time, endDate time.Time) ([]*attendance.AttendanceStudent, error)
}
