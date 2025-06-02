package leave

import (
	"context"
	"fmt"
	"service-holiday/internal/user"
	"time"

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
		fmt.Printf("User ID is empty\n")
	}

	if req.LeaveDate == "" {
		fmt.Printf("Leave date is empty\n")
	}

	user, err := s.userService.GetUserInfor(req.UserID)
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
		MaxEmployeesPerDay:  req.MaxEmployeesPerDay,
		AdvanceBookingDays:  req.AdvanceBookingDays,
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