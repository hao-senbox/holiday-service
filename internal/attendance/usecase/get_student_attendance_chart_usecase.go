package usecase

import (
	"context"
	"fmt"
	"sort"
	"time"

	"worktime-service/internal/gateway"
	"worktime-service/internal/shared"
)

var weekdayOrder = []time.Weekday{
	time.Monday,
	time.Tuesday,
	time.Wednesday,
	time.Thursday,
	time.Friday,
	time.Saturday,
	time.Sunday,
}

type GetStudentTemperatureChartUsecase interface {
	Execute(c context.Context, req shared.GetStudentTemperatureChartRequest) ([]*shared.StudentTemperatureChartResponse, error)
}

type getStudentTemperatureChartUsecase struct {
	attendanceRepository shared.AttendanceRepository
	termGateway          gateway.TermGateway
}

func NewGetStudentTemperatureChartUsecase(attendanceRepository shared.AttendanceRepository, termGateway gateway.TermGateway) GetStudentTemperatureChartUsecase {
	return &getStudentTemperatureChartUsecase{attendanceRepository: attendanceRepository, termGateway: termGateway}
}

func (u *getStudentTemperatureChartUsecase) Execute(c context.Context, req shared.GetStudentTemperatureChartRequest) ([]*shared.StudentTemperatureChartResponse, error) {

	term, err := u.termGateway.GetCurrentTermByOrgID(c, req.OrgID)
	if err != nil {
		return nil, err
	}

	StartDate, _ := time.Parse("2006-01-02", term.StartDate)
	EndDate, _ := time.Parse("2006-01-02", term.EndDate)

	records, err := u.attendanceRepository.GetStudentAttendanceInDateRange(c, req.StudentID, StartDate, EndDate)
	if err != nil {
		return nil, err
	}

	// Group records by year-week
	recordsByWeek := make(map[string][]*shared.AttendanceStudent)
	for _, r := range records {
		if r.Temperature != 0 {
			year, week := r.Date.ISOWeek()
			weekStr := fmt.Sprintf("%d-%02d", year, week)
			recordsByWeek[weekStr] = append(recordsByWeek[weekStr], r)
		}
	}

	var responses []*shared.StudentTemperatureChartResponse

	// Get sorted week keys
	var weekKeys []string
	for weekKey := range recordsByWeek {
		weekKeys = append(weekKeys, weekKey)
	}
	sort.Strings(weekKeys)

	// Create chart for each week in order
	for _, weekKey := range weekKeys {
		tempByWeekday := make(map[time.Weekday]float64)

		for _, r := range recordsByWeek[weekKey] {
			tempByWeekday[r.DayOfWeek] = r.Temperature
		}

		labels := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
		temperatures := make([]float64, len(weekdayOrder))

		for i, day := range weekdayOrder {
			if val, exists := tempByWeekday[day]; exists {
				temperatures[i] = val
			} else {
				temperatures[i] = 0
			}
		}

		// Find week date range for title
		var weekStart, weekEnd time.Time
		for _, r := range recordsByWeek[weekKey] {
			if weekStart.IsZero() || r.Date.Before(weekStart) {
				weekStart = r.Date
			}
			if weekEnd.IsZero() || r.Date.After(weekEnd) {
				weekEnd = r.Date
			}
		}

		title := fmt.Sprintf("Week %s",
			weekKey[len(weekKey)-2:],
		)

		response := &shared.StudentTemperatureChartResponse{
			Title:        title,
			Unit:         "°C",
			Labels:       labels,
			Temperatures: temperatures,
		}

		responses = append(responses, response)
	}

	return responses, nil
}
