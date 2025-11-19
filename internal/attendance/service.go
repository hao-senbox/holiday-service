package attendance

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"
	"worktime-service/helper"
	attendance "worktime-service/internal/attendance/usecase"
	"worktime-service/internal/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AttendanceService interface {
	CheckIn(c context.Context, req *CheckInRequest) error
	CheckOut(c context.Context, req *CheckOutRequest) error
	AttendanceStudent(c context.Context, req *AttendanceStudentRequest) error
	GetMyAttendance(c context.Context, userID string, month string, year string) ([]*DailyAttendance, error)
	GetAttendanceStudent(c context.Context, userID string, month string, year string) ([]*AttendanceStudent, error)
	GetAllAttendances(c context.Context, userID string, date string, page int, limit int) (*DailyAttendanceResponsePagination, error)
	GetStudentTemperatureChart(c context.Context, userID string) (*StudentTemperatureChartResponse, error)
}

type attendanceService struct {
	repo                              AttendanceRepository
	userService                       user.UserService
	getStudentTemperatureChartUsecase attendance.GetStudentTemperatureChartUsecase
}

func NewAttendanceService(repo AttendanceRepository, userService user.UserService, getStudentTemperatureChartUsecase attendance.GetStudentTemperatureChartUsecase) AttendanceService {
	return &attendanceService{
		repo:                              repo,
		userService:                       userService,
		getStudentTemperatureChartUsecase: getStudentTemperatureChartUsecase,
	}
}

func (s *attendanceService) CheckIn(c context.Context, req *CheckInRequest) error {

	if req.UserID == "" {
		return fmt.Errorf("user id is required")
	}

	_, err := s.userService.GetUserInfor(c, req.UserID)
	if err != nil {
		return err
	}

	now := time.Now()
	today := helper.GetStartOfDay(now)

	result, _ := s.repo.existingDailyAttendance(c, req.UserID, today)
	if result != nil {
		if result.CheckInTime != nil {
			return fmt.Errorf("user has already checked in today")
		}
	} else {
		attendanceLog := AttendanceLog{
			ID:        primitive.NewObjectID(),
			UserID:    req.UserID,
			LogDate:   today,
			LogTime:   now,
			LogType:   "check_in",
			Emotion:   req.Emotion,
			Notes:     &req.Notes,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err = s.repo.CreateAttendanceLog(c, &attendanceLog)
		if err != nil {
			return err
		}

		attendaceDaily := DailyAttendance{
			ID:                primitive.NewObjectID(),
			UserID:            req.UserID,
			DayOfWeek:         today.Weekday(),
			Date:              today,
			Status:            "working",
			CheckInTime:       &now,
			EmotionCheckIn:    req.Emotion,
			LunchDuration:     0,
			PercentWorkDay:    0,
			TotalWorkingHours: 0,
			CreatedAt:         now,
			UpdatedAt:         now,
		}

		err = s.repo.CreateDailyAttendance(c, &attendaceDaily)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *attendanceService) CheckOut(c context.Context, req *CheckOutRequest) error {

	if req.UserID == "" {
		return fmt.Errorf("user id is required")
	}

	_, err := s.userService.GetUserInfor(c, req.UserID)
	if err != nil {
		return err
	}

	now := time.Now()
	today := helper.GetStartOfDay(now)

	result, err := s.repo.existingDailyAttendance(c, req.UserID, today)
	if err != nil {
		return err
	}
	if result == nil || result.CheckInTime == nil {
		return fmt.Errorf("user has not checked in today")
	}
	if result.CheckoutTime != nil {
		return fmt.Errorf("user has already checked out today")
	}

	attendanceLog := AttendanceLog{
		ID:        primitive.NewObjectID(),
		UserID:    req.UserID,
		LogDate:   today,
		LogTime:   now,
		LogType:   "check_out",
		Emotion:   req.Emotion,
		Notes:     &req.Notes,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.repo.CreateAttendanceLog(c, &attendanceLog)
	if err != nil {
		return err
	}

	totalWorkingHours := helper.CalculateWorkingHours(*result.CheckInTime, now, req.DurianLunch)
	percentWorkday := (totalWorkingHours / 8) * 100

	result.CheckoutTime = &now
	result.LunchDuration = req.DurianLunch
	result.EMotionCheckOut = req.Emotion
	result.PercentWorkDay = percentWorkday
	result.TotalWorkingHours = totalWorkingHours
	result.UpdatedAt = now

	err = s.repo.UpdatedDailyAttendance(c, req.UserID, today, result)
	if err != nil {
		return err
	}

	return nil

}

func (s *attendanceService) GetMyAttendance(c context.Context, userID string, month string, year string) ([]*DailyAttendance, error) {

	if userID == "" {
		return nil, fmt.Errorf("user id is required")
	}

	if month == "" {
		return nil, fmt.Errorf("month is required")
	}

	if year == "" {
		return nil, fmt.Errorf("year is required")
	}

	monthInt, err := strconv.Atoi(month)
	if err != nil {
		return nil, err
	}

	if monthInt < 1 || monthInt > 12 {
		return nil, fmt.Errorf("invalid month")
	}

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return nil, err
	}

	firstDay := time.Date(yearInt, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)

	dailyAttendances, err := s.repo.GetMyAttendance(c, userID, firstDay, lastDay)

	if err != nil {
		return nil, err
	}

	data, err := s.getMonthlyAttandance(dailyAttendances, monthInt, yearInt)
	if err != nil {
		return nil, err
	}

	return data, nil

}

func (s *attendanceService) getMonthlyAttandance(myAttedances []*DailyAttendance, monthInt int, yearInt int) ([]*DailyAttendance, error) {

	firstDay := time.Date(yearInt, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)

	attandanceMap := make(map[string]*DailyAttendance)
	for _, item := range myAttedances {
		attandanceMap[item.Date.Format("2006-01-02")] = item
	}

	var dataMonthly []*DailyAttendance
	var status string

	for day := firstDay; !day.After(lastDay); day = day.AddDate(0, 0, 1) {

		dataKey := day.Format("2006-01-02")

		if item, ok := attandanceMap[dataKey]; ok {
			dataMonthly = append(dataMonthly, item)
		} else {
			if day.Weekday() == time.Sunday {
				status = "holiday"
			} else {
				status = "absent"
			}

			dataMonthly = append(dataMonthly, &DailyAttendance{
				ID:                primitive.NewObjectID(),
				UserID:            "",
				DayOfWeek:         day.Weekday(),
				Date:              day,
				Status:            status,
				CheckInTime:       nil,
				EmotionCheckIn:    "",
				CheckoutTime:      nil,
				LunchDuration:     0,
				EMotionCheckOut:   "",
				PercentWorkDay:    0,
				TotalWorkingHours: 0,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			})
		}

	}

	return dataMonthly, nil
}

func (s *attendanceService) AttendanceStudent(c context.Context, req *AttendanceStudentRequest) error {

	if req.UserID == "" {
		return fmt.Errorf("user id is required")
	}

	if req.Types == "" {
		return fmt.Errorf("types is required")
	}

	if req.Temperature == "" {
		return fmt.Errorf("temperature is required")
	}

	now := time.Now()
	today := helper.GetStartOfDay(now)

	switch req.Types {
	case "arrive":
		result, _ := s.repo.existingDailyAttendanceStudent(c, req.UserID, today)
		if result != nil {
			// if result.CheckInTime != nil {
			// 	return fmt.Errorf("student has already checked in today")
			// }

			result.CheckInTime = &now
			result.UpdatedAt = now

			if err := s.repo.UpdateDailyAttendanceStudent(c, result); err != nil {
				return err
			}
		} else {
			if err := s.repo.CreateAttendanceLog(c, &AttendanceLog{
				ID:          primitive.NewObjectID(),
				UserID:      req.UserID,
				Temperature: req.Temperature,
				LogDate:     today,
				LogTime:     now,
				LogType:     "arrive",
				Notes:       &req.Notes,
				CreatedBy:   &req.CreatedBy,
				CreatedAt:   now,
				UpdatedAt:   now,
			}); err != nil {
				return err
			}

			if err := s.repo.CreateDailyAttendanceStudent(c, &AttendanceStudent{
				ID:          primitive.NewObjectID(),
				UserID:      req.UserID,
				DayOfWeek:   today.Weekday(),
				Temperature: req.Temperature,
				Date:        today,
				CheckInTime: &now,
				Types:       req.Types,
				Note:        req.Notes,
				CreatedBy:   req.CreatedBy,
				CreatedAt:   now,
				UpdatedAt:   now,
			}); err != nil {
				return err
			}
		}

	case "leave":
		result, _ := s.repo.existingDailyAttendanceStudent(c, req.UserID, today)
		// if result == nil {
		// 	return fmt.Errorf("student has not checked in today")
		// }

		// if result.CheckOutTime != nil {
		// 	return fmt.Errorf("student has already checked out today")
		// }

		// if result.CheckInTime == nil {
		// 	return fmt.Errorf("student has not checked in today")
		// }

		if err := s.repo.CreateAttendanceLog(c, &AttendanceLog{
			ID:        primitive.NewObjectID(),
			UserID:    req.UserID,
			LogDate:   today,
			LogTime:   now,
			LogType:   "leave",
			Notes:     &req.Notes,
			CreatedBy: &req.CreatedBy,
			CreatedAt: now,
			UpdatedAt: now,
		}); err != nil {
			return err
		}

		result.CheckOutTime = &now
		result.UpdatedAt = now

		if err := s.repo.UpdateDailyAttendanceStudent(c, result); err != nil {
			return err
		}

	case "absent":
		if req.Date == "" {
			return fmt.Errorf("date is required for absent")
		}

		dateParse, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return fmt.Errorf("invalid date format, must be YYYY-MM-DD")
		}

		if err := s.repo.CreateAttendanceLog(c, &AttendanceLog{
			ID:        primitive.NewObjectID(),
			UserID:    req.UserID,
			LogDate:   dateParse,
			LogTime:   now,
			LogType:   "absent",
			Notes:     &req.Notes,
			CreatedBy: &req.CreatedBy,
			CreatedAt: now,
			UpdatedAt: now,
		}); err != nil {
			return err
		}

		if err := s.repo.CreateDailyAttendanceStudent(c, &AttendanceStudent{
			ID:           primitive.NewObjectID(),
			UserID:       req.UserID,
			DayOfWeek:    dateParse.Weekday(),
			Date:         dateParse,
			CheckInTime:  nil,
			CheckOutTime: nil,
			Types:        req.Types,
			Note:         req.Notes,
			CreatedBy:    req.CreatedBy,
			CreatedAt:    now,
			UpdatedAt:    now,
		}); err != nil {
			return err
		}

	default:
		return fmt.Errorf("invalid type, must be arrive, leave, or absent")
	}

	return nil
}

func (s *attendanceService) GetAttendanceStudent(c context.Context, userID string, month string, year string) ([]*AttendanceStudent, error) {

	monthInt, err := strconv.Atoi(month)
	if err != nil {
		return nil, err
	}

	if monthInt < 1 || monthInt > 12 {
		return nil, fmt.Errorf("invalid month")
	}

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return nil, err
	}

	firstDay := time.Date(yearInt, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)

	data, err := s.repo.GetAttendanceStudent(c, userID, firstDay, lastDay)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *attendanceService) GetAllAttendances(c context.Context, userID string, date string, page int, limit int) (*DailyAttendanceResponsePagination, error) {

	var dateParse *time.Time
	if date != "" {
		t, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil, err
		}
		dateParse = &t
	}

	attendances, totalCount, err := s.repo.GetAllAttendances(c, userID, dateParse, page, limit)
	if err != nil {
		return nil, err
	}

	totalPages := (totalCount + int64(limit) - 1) / int64(limit)
	var data []*DailyAttendanceUser

	userCache := make(map[string]*user.UserInfor)
	for _, attendance := range attendances {

		var userInfo *user.UserInfor
		if attendance.UserID != "" {
			if cached, ok := userCache[attendance.UserID]; ok {
				userInfo = cached
			} else {
				user, err := s.userService.GetUserInfor(c, attendance.UserID)
				if err != nil {
					log.Println("Failed to get user information:", err)
				}
				userCache[attendance.UserID] = user
				userInfo = user
			}
		}

		data = append(data, &DailyAttendanceUser{
			ID:                attendance.ID.Hex(),
			UserInfor:         userInfo,
			DayOfWeek:         attendance.DayOfWeek,
			Date:              attendance.Date,
			Status:            attendance.Status,
			CheckInTime:       formatTimePtr(attendance.CheckInTime),
			EmotionCheckIn:    attendance.EmotionCheckIn,
			CheckoutTime:      formatTimePtr(attendance.CheckoutTime),
			LunchDuration:     attendance.LunchDuration,
			EMotionCheckOut:   attendance.EMotionCheckOut,
			PercentWorkDay:    attendance.PercentWorkDay,
			TotalWorkingHours: attendance.TotalWorkingHours,
			CreatedAt:         formatTimePtr(&attendance.CreatedAt),
			UpdatedAt:         formatTimePtr(&attendance.UpdatedAt),
		})
	}

	paginate := Pagination{
		TotalCount: totalCount,
		TotalPages: totalPages,
		Page:       int64(page),
		Limit:      int64(limit),
	}

	return &DailyAttendanceResponsePagination{
		Pagination:      paginate,
		DailyAttendance: data,
	}, nil

}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Add(7 * time.Hour).Format("2006-01-02 15:04:05")
}

func (s *attendanceService) GetStudentTemperatureChart(c context.Context, userID string) (*StudentTemperatureChartResponse, error) {
	return s.getStudentTemperatureChartUsecase.Execute(c, userID)
}
