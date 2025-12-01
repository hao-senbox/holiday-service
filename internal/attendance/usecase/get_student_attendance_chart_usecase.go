package usecase

import (
	"context"
	"fmt"
	"time"

	"worktime-service/internal/gateway"
	"worktime-service/internal/shared"
	"worktime-service/internal/user"
)
type GetStudentTemperatureChartUsecase interface {
	Execute(c context.Context, req shared.GetStudentTemperatureChartRequest) ([]*shared.StudentTemperatureChartResponse, error)
}

type getStudentTemperatureChartUsecase struct {
	attendanceRepository shared.AttendanceRepository
	userService          user.UserService
	termGateway          gateway.TermGateway
}

func NewGetStudentTemperatureChartUsecase(attendanceRepository shared.AttendanceRepository, userService user.UserService, termGateway gateway.TermGateway) GetStudentTemperatureChartUsecase {
	return &getStudentTemperatureChartUsecase{attendanceRepository: attendanceRepository, userService: userService, termGateway: termGateway}
}

func (u *getStudentTemperatureChartUsecase) Execute(c context.Context, req shared.GetStudentTemperatureChartRequest) ([]*shared.StudentTemperatureChartResponse, error) {
	user, err := u.userService.GetCurrentUser(c)
	if err != nil {
		return nil, err
	}
	fmt.Printf("organization_id: %s", user.OrganizationAdmin.ID)
	term, err := u.termGateway.GetCurrentTermByOrgID(c, user.OrganizationAdmin.ID)
	if err != nil {
		return nil, err
	}

	StartDate, _ := time.Parse("2006-01-02", term.StartDate)
	EndDate, _ := time.Parse("2006-01-02", term.EndDate)

	records, err := u.attendanceRepository.GetStudentAttendanceInDateRange(c, req.StudentID, StartDate, EndDate)
	if err != nil {
		return nil, err
	}

	// Group records by date
	recordsByDate := make(map[string]*shared.AttendanceStudent)
	for _, r := range records {
		if r.Temperature != 0 {
			dateStr := r.Date.Format("2006-01-02")
			recordsByDate[dateStr] = r
		}
	}

	var responses []*shared.StudentTemperatureChartResponse

	// Create response for each day in the term range
	for d := StartDate; d.Before(EndDate) || d.Equal(EndDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")

		var temperature float64
		if record, exists := recordsByDate[dateStr]; exists {
			temperature = record.Temperature
		}

		response := &shared.StudentTemperatureChartResponse{
			Date:         dateStr,
			Unit:         "°C",
			Labels:       d.Weekday().String()[:3],
			Temperatures: temperature,
		}

		responses = append(responses, response)
	}

	return responses, nil
}
