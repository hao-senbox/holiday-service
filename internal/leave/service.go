package leave

import (
	"context"
	"errors"
	"fmt"
	"time"
	"worktime-service/internal/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LeaveService interface {
	CreateRequestLeave(ctx context.Context, req *CreateLeaveRequest) error
	GetAllLeaveCalendar(ctx context.Context, date string) ([]*DailyLeaveSolt, error)
	GetDetailLeaveCalendar(ctx context.Context, id string) (*DailyLeaveSolt, error)
	GetSettings(ctx context.Context) *Setting
	UpdateSetting(ctx context.Context, req *SettingRequest, id string) (*Setting, error)
	GetMyRequest(ctx context.Context, userID string) (map[string][]*LeaveRequests, error)
	EditMaxSlot(ctx context.Context, req *EditSlotRequest, id string) error
	GetPendingRequest(ctx context.Context) ([]*LeaveRequests, error)
	DeleteRequestLeave(ctx context.Context, req *DeleteLeaveRequest) error
	UpdateRequestLeave(ctx context.Context, req *UpdateRequest, id string) error
	GetStatistical(ctx context.Context, dateFrom string, dateTo string) (*LeaveStatistical, error)
	// AddCronLeavesBalance(ctx context.Context) error
	// GetLeaveBalanceUser(ctx context.Context, userID string) (interface{}, error)
}

type leaveService struct {
	leaveRepository LeaveRepository
	userService     user.UserService
}

func NewLeaveService(leaveRepository LeaveRepository, userService user.UserService) LeaveService {
	return &leaveService{
		leaveRepository: leaveRepository,
		userService:     userService,
	}
}

func (s *leaveService) CreateRequestLeave(ctx context.Context, req *CreateLeaveRequest) error {

	if req.UserID == "" {
		return errors.New("user ID is empty")
	}

	if req.LeaveDate == "" {
		return errors.New("leave date is empty")
	}

	user, err := s.userService.GetUserInfor(ctx, req.UserID)
	if err != nil {
		return err
	}

	dateParse, err := time.Parse("2006-01-02", req.LeaveDate)
	if err != nil {
		return err
	}

	leaveItem := LeaveRequests{
		LeaveDate:   dateParse,
		UserID:      req.UserID,
		UserName:    user.UserName,
		Reason:      req.Reason,
		RequestedAt: time.Now(),
	}

	err = s.leaveRepository.CreateLeave(ctx, &leaveItem)
	if err != nil {
		return err
	}

	return nil
}

func (s *leaveService) GetAllLeaveCalendar(ctx context.Context, date string) ([]*DailyLeaveSolt, error) {

	var dateFilter *time.Time

	if date != "" {
		parsed, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil, err
		}
		dateFilter = &parsed
	}

	data, err := s.leaveRepository.GetDailyLeaveSlots(ctx, dateFilter)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *leaveService) GetDetailLeaveCalendar(ctx context.Context, id string) (*DailyLeaveSolt, error) {

	if id == "" {
		return nil, fmt.Errorf("id is empty")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	data, err := s.leaveRepository.GetDetailLeaveSlots(ctx, objectID)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *leaveService) GetSettings(ctx context.Context) *Setting {
	return s.leaveRepository.GetSettings(ctx)
}

func (s *leaveService) UpdateSetting(ctx context.Context, req *SettingRequest, id string) (*Setting, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	setting := Setting{
		MaxEmployeesPerDay: req.MaxEmployeesPerDay,
		AdvanceBookingDays: req.AdvanceBookingDays,
		UpdatedAt:          time.Now(),
	}

	settingResult, err := s.leaveRepository.UpdateSetting(ctx, &setting, objectID)
	if err != nil {
		return nil, err
	}

	return settingResult, nil

}

func (s *leaveService) GetMyRequest(ctx context.Context, userID string) (map[string][]*LeaveRequests, error) {

	data, err := s.leaveRepository.GetMyRequest(ctx, userID)
	if err != nil {
		return nil, err
	}

	grouped := s.groupedRequestType(data)

	return grouped, nil

}

func (s *leaveService) EditMaxSlot(ctx context.Context, req *EditSlotRequest, id string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	detailLeaveSlots, err := s.leaveRepository.GetDetailLeaveSlots(ctx, objectID)
	if err != nil {
		return err
	}

	var AvailableSlot int

	if req.MaxSlot > len(detailLeaveSlots.ConfirmedLeaves) {
		AvailableSlot = req.MaxSlot - len(detailLeaveSlots.ConfirmedLeaves)
	} else {
		AvailableSlot = 0
	}

	if req.MaxSlot <= 0 {
		return fmt.Errorf("max slot is required")
	}

	return s.leaveRepository.EditMaxSlot(ctx, req.MaxSlot, AvailableSlot, objectID)

}

func (s *leaveService) groupedRequestType(leaveRequest []*LeaveRequests) map[string][]*LeaveRequests {

	grouped := make(map[string][]*LeaveRequests)

	for _, item := range leaveRequest {
		grouped[item.RequestType] = append(grouped[item.RequestType], item)
	}

	return grouped

}

func (s *leaveService) GetPendingRequest(ctx context.Context) ([]*LeaveRequests, error) {

	return s.leaveRepository.GetPendingRequest(ctx)

}

func (s *leaveService) DeleteRequestLeave(ctx context.Context, req *DeleteLeaveRequest) error {

	if req.LeaveDate == "" {
		return fmt.Errorf("leave date is required")
	}

	if req.UserID == "" {
		return fmt.Errorf("user id is required")
	}

	dateParse, err := time.Parse("2006-01-02", req.LeaveDate)
	if err != nil {
		return err
	}

	return s.leaveRepository.DeleteRequestLeave(ctx, &dateParse, req.UserID)

}

func (s *leaveService) UpdateRequestLeave(ctx context.Context, req *UpdateRequest, id string) error {

	if req.Types == "" {
		return fmt.Errorf("types is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return s.leaveRepository.UpdateRequestLeave(ctx, req.Types, objectID)

}

func (s *leaveService) GetStatistical(ctx context.Context, dateFrom string, dateTo string) (*LeaveStatistical, error) {

	if dateFrom == "" {
		return nil, fmt.Errorf("date from is empty")
	}

	if dateTo == "" {
		return nil, fmt.Errorf("date to is empty")
	}

	dateFromParse, err := time.Parse("2006-01-02", dateFrom)
	if err != nil {
		return nil, err
	}

	dateToParse, err := time.Parse("2006-01-02", dateTo)
	if err != nil {
		return nil, err
	}

	leaves, err := s.leaveRepository.GetAllLeaves(ctx, &dateFromParse, &dateToParse)
	if err != nil {
		return nil, err
	}

	return caculateStatistical(leaves), nil

}

func caculateStatistical(leaves []*LeaveRequests) *LeaveStatistical {

	stats := &LeaveStatistical{}

	stats.TotalRequested = len(leaves)

	statusCount := make(map[string]int)
	requestTypeCount := make(map[string]int)
	userCount := make(map[string]*UserStats)
	monthCount := make(map[string]*MonthlyStats)
	weekdayCount := make(map[string]int)

	for _, item := range leaves {

		statusCount[item.Status]++
		requestTypeCount[item.RequestType]++

		if userStat, ok := userCount[item.UserID]; ok {
			userStat.TotalRequests++
			if item.Status == "approved"  || item.Status == "confirmed" {
				userStat.ApprovedRequests++
			}
		} else {
			var approvedRequests int
			if item.Status == "approved" || item.Status == "confirmed" {
				approvedRequests = 1
			} else {
				approvedRequests = 0
			}
			userStat := &UserStats{
				UserID:           item.UserID,
				UserName:         item.UserName,
				TotalRequests:    1,
				ApprovedRequests: approvedRequests,
			}
			userCount[item.UserID] = userStat
		}

		monthKey := item.LeaveDate.Format("2006-01")

		if monthStat, ok := monthCount[monthKey]; ok {
			monthStat.Count++
			if item.Status == "approved" {
				monthStat.Approved++
			} else if item.Status == "rejected" {
				monthStat.Rejected++
			}
		} else {
			monthStat := &MonthlyStats{
				Month:    monthKey,
				Count:    1,
				Approved: 0,
				Rejected: 0,
			}
			if item.Status == "approved" {
				monthStat.Approved++
			} else if item.Status == "rejected" {
				monthStat.Rejected++
			}
			monthCount[monthKey] = monthStat
		}

		weekdayCount[item.LeaveDate.Weekday().String()]++
	}

	stats.TotalConfirmed = statusCount["approved"]
	stats.TotalPending = statusCount["pending"]
	stats.TotalRejected = statusCount["rejected"]
	stats.ImmediateRequests = requestTypeCount["immediate"]

	if stats.TotalRequested > 0 {
		stats.ApproveRate = float64(stats.TotalConfirmed) / float64(stats.TotalRequested) * 100
	}

	for _, monthStat := range monthCount {
		stats.RequestsByMonth = append(stats.RequestsByMonth, *monthStat)
	}

	for weekday, count := range weekdayCount {
		stats.RequestByWeekDay = append(stats.RequestByWeekDay, WeeklyStats{
			Weekday: weekday,
			Count:   count,
		})
	}

	for _, userStat := range userCount {
		if userStat.TotalRequests > 0 {
			userStat.ApprovalRate = float64(userStat.ApprovedRequests) / float64(userStat.TotalRequests) * 100
		}
		stats.TopRequestUsers = append(stats.TopRequestUsers, *userStat)
	}

	if len(userCount) > 0 {
		stats.AverageRequestsPerUser = float64(stats.TotalRequested) / float64(len(userCount))
	}

	return stats
}

// func (s *leaveService) GetLeaveBalanceUser(ctx context.Context, userID string) (interface{}, error) {
	
// 	data, err := s.userService.GetAllUser()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return data, nil
// }

// func (s *leaveService) AddCronLeavesBalance(ctx context.Context) error {

// 	dataUser, err := s.userService.GetAllUser()
// 	if err != nil {
// 		return err
// 	}

// 	dataLeavesBanlance, err := s.leaveRepository.GetAllLeaveBalance(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	balanceMap := make(map[string]*UserLeaveBalance)
// 	currentYear := time.Now().Year()

// 	for i := range dataLeavesBanlance {
// 		key := fmt.Sprintf("%s_%d", dataLeavesBanlance[i].UserID, currentYear)
// 		balanceMap[key] = dataLeavesBanlance[i]
// 	}

// 	var newBalance []*UserLeaveBalance
// 	var updateBalance []UserLeaveBalance
// 	var transaction []LeaveTransaction

// 	for _, item := range dataUser {

// 		key := fmt.Sprintf("%s_%d", item.UserID, currentYear)

// 		if existingBalance, ok := balanceMap[key]; ok {

// 			existingBalance.TotalLeaveBanlance += 1
// 			existingBalance.RemaindingLeave += 1
// 			existingBalance.LastUpdated = time.Now()

// 			updateBalance = append(updateBalance, *existingBalance)

// 		} else {

// 			newBalance = append(newBalance, &UserLeaveBalance{
// 				ID:                 primitive.NewObjectID(),
// 				UserID:             item.UserID,
// 				Year:               currentYear,
// 				TotalLeaveBanlance: 1,
// 				UsedLeave:          0,
// 				RemaindingLeave:    1,
// 				LastUpdated:        time.Now(),
// 			})

// 			reason := fmt.Sprintf("Monthly earned leave - %s", time.Now().Format("2006-01"))

// 			transaction = append(transaction, LeaveTransaction{
// 				ID:              primitive.NewObjectID(),
// 				UserID:          item.UserID,
// 				TransactionType: "EARNED",
// 				Date:            time.Now(),
// 				Reason:          &reason,
// 				CreatedAt:       time.Now(),
// 				UpdatedAt:       time.Now(),
// 			})

// 		}
// 	}

// 	if len(newBalance) > 0 {
// 		err = s.leaveRepository.CreateLeaveBalance(ctx, newBalance)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	if len(updateBalance) > 0 {
// 		err = s.leaveRepository.UpdateLeaveBalance(ctx, updateBalance)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	if len(transaction) > 0 {
// 		err = s.leaveRepository.CreateLeaveTransaction(ctx, transaction)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil

// }