package usecase

import (
	"context"
	"strconv"
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
	Execute(c context.Context, req shared.GetStudentTemperatureChartRequest) (*shared.StudentTemperatureChartResponse, error)
}

type getStudentTemperatureChartUsecase struct {
	attendanceRepository shared.AttendanceRepository
	termGateway          gateway.TermGateway
}

func NewGetStudentTemperatureChartUsecase(attendanceRepository shared.AttendanceRepository, termGateway gateway.TermGateway) GetStudentTemperatureChartUsecase {
	return &getStudentTemperatureChartUsecase{attendanceRepository: attendanceRepository, termGateway: termGateway}
}

func (u *getStudentTemperatureChartUsecase) Execute(c context.Context, req shared.GetStudentTemperatureChartRequest) (*shared.StudentTemperatureChartResponse, error) {

	term, err := u.termGateway.GetTermByID(c, req.TermID)
	if err != nil {
		return nil, err
	}

	StartDate, _ := time.Parse("2006-01-02", term.StartDate)
	EndDate, _ := time.Parse("2006-01-02", term.EndDate)

	records, err := u.attendanceRepository.GetStudentAttendanceInDateRange(c, req.StudentID, StartDate, EndDate)
	if err != nil {
		return nil, err
	}

	tempByWeekday := make(map[time.Weekday]float64)

	for _, r := range records {
		if r.Temperature != "" {
			val, _ := strconv.ParseFloat(r.Temperature, 64)
			tempByWeekday[r.DayOfWeek] = val
		}
	}

	labels := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	temperatures := make([]float64, len(weekdayOrder))
	status := make([]string, len(weekdayOrder))

	for i, day := range weekdayOrder {
		if val, exists := tempByWeekday[day]; exists {
			temperatures[i] = val
			status[i] = evaluateStatus(val)
		} else {
			temperatures[i] = 0
			status[i] = "no_data"
		}
	}

	response := &shared.StudentTemperatureChartResponse{
		Title:        "Student Temperature",
		Unit:         "°C",
		Labels:       labels,
		Temperatures: temperatures,
		Status:       status,
	}

	return response, nil
}

func evaluateStatus(value float64) string {
	switch {
	case value >= 40:
		return "danger"
	case value >= 36.3 && value <= 37.1:
		return "normal"
	default:
		return "warning"
	}
}
