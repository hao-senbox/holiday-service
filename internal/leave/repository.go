package leave

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LeaveRepository interface {
	CreateLeave(ctx context.Context, leaveItem *LeaveRequests) error
	GetSettings(ctx context.Context) *Setting
	UpdateSetting(ctx context.Context, setting *Setting, id primitive.ObjectID) (*Setting, error)
	GetDailyLeaveSlots(ctx context.Context, date *time.Time) ([]*DailyLeaveSolt, error)
	GetDetailLeaveSlots(ctx context.Context, id primitive.ObjectID) (*DailyLeaveSolt, error)
	GetMyRequest(ctx context.Context, userID string) ([]*LeaveRequests, error)
	EditMaxSlot(ctx context.Context, maxSlot int, availableSlot int, id primitive.ObjectID) error
	GetPendingRequest(ctx context.Context) ([]*LeaveRequests, error)
	DeleteRequestLeave(ctx context.Context, date *time.Time, userID string) error
	UpdateRequestLeave(ctx context.Context, types string, id primitive.ObjectID) error
	GetAllLeaves(ctx context.Context, dateFrom *time.Time, dateTo *time.Time) ([]*LeaveRequests, error)
	GetAllLeaveBalance(ctx context.Context) ([]*UserLeaveBalance, error)
	CreateLeaveBalance(ctx context.Context, leaveBalance []*UserLeaveBalance) error
}

type leaveRepository struct {
	collectionLeave           *mongo.Collection
	collectionLeaveSetting    *mongo.Collection
	collectionDailyLeaveSlots *mongo.Collection
	collectionLeaveBalance    *mongo.Collection
}

func NewLeaveRepository(collectionLeave *mongo.Collection,
	collectionLeaveSetting *mongo.Collection,
	collectionDailyLeaveSlots *mongo.Collection,
	collectionLeaveBalance *mongo.Collection) LeaveRepository {
	return &leaveRepository{
		collectionLeave:           collectionLeave,
		collectionLeaveSetting:    collectionLeaveSetting,
		collectionDailyLeaveSlots: collectionDailyLeaveSlots,
		collectionLeaveBalance:    collectionLeaveBalance,
	}
}

func (r *leaveRepository) GetSettings(ctx context.Context) *Setting {

	var setting Setting

	err := r.collectionLeaveSetting.FindOne(ctx, bson.M{}).Decode(&setting)
	if err == mongo.ErrNoDocuments {

		setting = Setting{
			ID:                 primitive.NewObjectID(),
			MaxEmployeesPerDay: 5,
			AdvanceBookingDays: 7,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}

		r.collectionLeaveSetting.InsertOne(ctx, setting)
	}

	return &setting
}

func (r *leaveRepository) UpdateSetting(ctx context.Context, setting *Setting, id primitive.ObjectID) (*Setting, error) {

	filter := bson.M{"_id": id}

	update := bson.M{
		"$set": bson.M{
			"max_employees_per_day": setting.MaxEmployeesPerDay,
			"advance_booking_days":  setting.AdvanceBookingDays,
			"updated_at":            time.Now(),
		},
	}

	_, err := r.collectionLeaveSetting.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	var settingUpdated Setting

	err = r.collectionLeaveSetting.FindOne(ctx, filter).Decode(&settingUpdated)
	if err != nil {
		return nil, err
	}

	return &settingUpdated, nil

}

func (r *leaveRepository) getDailyLeaveSlots(ctx context.Context, date string) *DailyLeaveSolt {

	dateParse, _ := time.Parse("2006-01-02", date)

	filter := bson.M{"date": dateParse}

	var dailyLeaveSlot DailyLeaveSolt

	err := r.collectionDailyLeaveSlots.FindOne(ctx, filter).Decode(&dailyLeaveSlot)

	if mongo.ErrNoDocuments == err {
		setting := r.GetSettings(ctx)
		dailyLeaveSlot = DailyLeaveSolt{
			ID:              primitive.NewObjectID(),
			MaxSlot:         setting.MaxEmployeesPerDay,
			AvailableSlot:   setting.MaxEmployeesPerDay,
			Date:            dateParse,
			ConfirmedLeaves: []ConfirmedLeave{},
			PendingRequests: []PendingRequest{},
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		r.collectionDailyLeaveSlots.InsertOne(ctx, dailyLeaveSlot)
	}

	return &dailyLeaveSlot
}

func (r *leaveRepository) CreateLeave(ctx context.Context, leaveItem *LeaveRequests) error {

	dailyLeaveSlot := r.getDailyLeaveSlots(ctx, leaveItem.LeaveDate.Format("2006-01-02"))

	requestType := "immediate"
	status := "confirmed"

	if dailyLeaveSlot.AvailableSlot <= 0 {
		requestType = "wishlist"
		status = "pending"
	}

	leaveItem = &LeaveRequests{
		ID:          primitive.NewObjectID(),
		LeaveDate:   leaveItem.LeaveDate,
		UserID:      leaveItem.UserID,
		UserName:    leaveItem.UserName,
		Reason:      leaveItem.Reason,
		RequestedAt: leaveItem.RequestedAt,
		Status:      status,
		RequestType: requestType,
	}

	check, msg := r.checkRequestExists(ctx, leaveItem)
	if check {
		return fmt.Errorf("leave request exists: %s", msg)
	}

	_, err := r.collectionLeave.InsertOne(ctx, leaveItem)
	if err != nil {
		return err
	}

	err = r.updateDailyLeaveSlots(ctx, leaveItem.LeaveDate.Format("2006-01-02"), leaveItem)
	if err != nil {
		return err
	}

	return nil

}

func (r *leaveRepository) updateDailyLeaveSlots(ctx context.Context, date string, leaveItem *LeaveRequests) error {

	dateParse, _ := time.Parse("2006-01-02", date)

	if leaveItem.RequestType == "immediate" {

		confirmedLeave := ConfirmedLeave{
			UserID:    leaveItem.UserID,
			UserName:  leaveItem.UserName,
			ApproveAt: time.Now(),
		}

		filter := bson.M{"date": dateParse}
		update := bson.M{
			"$inc":  bson.M{"available_slot": -1},
			"$push": bson.M{"confirmed_leaves": confirmedLeave},
		}

		_, err := r.collectionDailyLeaveSlots.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
	} else {

		pendingRequest := PendingRequest{
			LeaveID:   leaveItem.ID,
			UserID:    leaveItem.UserID,
			UserName:  leaveItem.UserName,
			Status:    leaveItem.Status,
			RequestAt: leaveItem.RequestedAt,
		}

		filter := bson.M{"date": dateParse}
		update := bson.M{
			"$push": bson.M{"pending_requests": pendingRequest},
		}

		_, err := r.collectionDailyLeaveSlots.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}

	}

	return nil

}

func (r *leaveRepository) checkRequestExists(ctx context.Context, leaveItem *LeaveRequests) (bool, string) {

	filterConfirmed := bson.M{
		"date":                     leaveItem.LeaveDate,
		"confirmed_leaves.user_id": leaveItem.UserID,
	}

	filterPending := bson.M{
		"date":                     leaveItem.LeaveDate,
		"pending_requests.user_id": leaveItem.UserID,
	}

	countConfirmed, err := r.collectionDailyLeaveSlots.CountDocuments(ctx, filterConfirmed)
	if err != nil {
		return false, ""
	}
	if countConfirmed > 0 {
		return true, "User has successfully registered for leave"
	}

	countPending, err := r.collectionDailyLeaveSlots.CountDocuments(ctx, filterPending)
	if err != nil {
		return false, ""
	}
	if countPending > 0 {
		return true, "User has pending leave request"
	}

	return false, ""
}

func (r *leaveRepository) GetDailyLeaveSlots(ctx context.Context, date *time.Time) ([]*DailyLeaveSolt, error) {

	var dailyLeaveSlots []*DailyLeaveSolt

	filter := bson.M{}
	if date != nil {
		start := date.AddDate(0, 0, -30)
		end := date.AddDate(0, 0, 30)

		filter["date"] = bson.M{
			"$gte": start,
			"$lte": end,
		}

	}

	cursor, err := r.collectionDailyLeaveSlots.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &dailyLeaveSlots); err != nil {
		return nil, err
	}

	return dailyLeaveSlots, nil

}

func (r *leaveRepository) GetDetailLeaveSlots(ctx context.Context, id primitive.ObjectID) (*DailyLeaveSolt, error) {

	var dailyLeaveSlot *DailyLeaveSolt

	filter := bson.M{"_id": id}
	err := r.collectionDailyLeaveSlots.FindOne(ctx, filter).Decode(&dailyLeaveSlot)
	if err != nil {
		return nil, err
	}

	return dailyLeaveSlot, nil

}

func (r *leaveRepository) GetMyRequest(ctx context.Context, userID string) ([]*LeaveRequests, error) {

	var leaveRequests []*LeaveRequests

	filter := bson.M{"user_id": userID}
	cursor, err := r.collectionLeave.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &leaveRequests); err != nil {
		return nil, err
	}

	return leaveRequests, nil

}

func (r *leaveRepository) EditMaxSlot(ctx context.Context, maxSlot int, availableSlot int, id primitive.ObjectID) error {

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"max_slot":       maxSlot,
			"available_slot": availableSlot,
		},
	}

	_, err := r.collectionDailyLeaveSlots.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil

}

func (r *leaveRepository) GetPendingRequest(ctx context.Context) ([]*LeaveRequests, error) {

	var leaveRequests []*LeaveRequests

	filter := bson.M{"status": "pending"}

	cursor, err := r.collectionLeave.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &leaveRequests)
	if err != nil {
		return nil, err
	}

	return leaveRequests, nil

}

func (r *leaveRepository) DeleteRequestLeave(ctx context.Context, date *time.Time, userID string) error {

	var leaveRequests LeaveRequests
	var dailyLeaveSlot DailyLeaveSolt
	leaveFilter := bson.M{"leave_date": date, "user_id": userID}

	err := r.collectionLeave.FindOneAndDelete(ctx, leaveFilter).Decode(&leaveRequests)
	if err != nil {
		return err
	}

	slotFilter := bson.M{"date": date}
	var update bson.M

	if leaveRequests.RequestType == "immediate" {

		update = bson.M{
			"$inc": bson.M{"available_slot": 1},
			"$pull": bson.M{
				"confirmed_leaves": bson.M{"user_id": userID},
			},
		}

	} else {

		update = bson.M{
			"$pull": bson.M{
				"pending_requests": bson.M{"user_id": userID},
			},
		}

	}

	_, err = r.collectionDailyLeaveSlots.UpdateOne(ctx, slotFilter, update)
	if err != nil {
		return err
	}

	err = r.collectionDailyLeaveSlots.FindOne(ctx, slotFilter).Decode(&dailyLeaveSlot)
	if err != nil {
		return err
	}

	if len(dailyLeaveSlot.ConfirmedLeaves) == 0 && len(dailyLeaveSlot.PendingRequests) == 0 {
		_, err := r.collectionDailyLeaveSlots.DeleteOne(ctx, slotFilter)
		if err != nil {
			return err
		}
	}

	return nil

}

func (r *leaveRepository) UpdateRequestLeave(ctx context.Context, types string, id primitive.ObjectID) error {

	filter := bson.M{"_id": id, "request_type": "wishlist"}
	update := bson.M{"$set": bson.M{"status": types}}

	_, err := r.collectionLeave.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil

}

func (r *leaveRepository) GetAllLeaves(ctx context.Context, dateFrom *time.Time, dateTo *time.Time) ([]*LeaveRequests, error) {

	var leaveRequests []*LeaveRequests

	filter := bson.M{"leave_date": bson.M{"$gte": *dateFrom, "$lte": *dateTo}}

	cursor, err := r.collectionLeave.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &leaveRequests)
	if err != nil {
		return nil, err
	}

	return leaveRequests, nil

}

func (r *leaveRepository) GetAllLeaveBalance(ctx context.Context) ([]*UserLeaveBalance, error) {

	var userLeaveBalances []*UserLeaveBalance

	cursor, err := r.collectionLeaveBalance.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &userLeaveBalances)
	if err != nil {
		return nil, err
	}

	if userLeaveBalances == nil {
		return []*UserLeaveBalance{}, nil
	}

	return userLeaveBalances, nil

}

func (r *leaveRepository) CreateLeaveBalance(ctx context.Context, leaveBalance []*UserLeaveBalance) error {

	docs := make([]interface{}, len(leaveBalance))
	for i, leave := range leaveBalance {
		docs[i] = leave
	}

	_, err := r.collectionLeaveBalance.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	return nil
	
}