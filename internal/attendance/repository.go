package attendance

import (
	"context"
	"fmt"
	"time"
	"worktime-service/helper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AttendanceRepository interface {
	CreateAttendanceLog(c context.Context, attendanceLog *AttendanceLog) error
	CreateDailyAttendance(c context.Context, dailyAttendance *DailyAttendance) error
	existingDailyAttendance(c context.Context, userID string, today time.Time) (*DailyAttendance, error)
	UpdatedDailyAttendance(c context.Context, userID string, today time.Time, dailyAttendance *DailyAttendance) error
	GetMyAttendance(c context.Context, userID string, firstDay time.Time, lastDay time.Time) ([]*DailyAttendance, error)
}

type attendanceRepository struct {
	collectionAttendance      *mongo.Collection
	collectionDailyAttendance *mongo.Collection
}

func NewAttendanceRepository(collectionAttendance *mongo.Collection, collectionDailyAttendance *mongo.Collection) AttendanceRepository {
	return &attendanceRepository{
		collectionAttendance:      collectionAttendance,
		collectionDailyAttendance: collectionDailyAttendance,
	}
}

func (r *attendanceRepository) existingDailyAttendance(c context.Context, userID string, today time.Time) (*DailyAttendance, error) {

	var existingDailyAttendance DailyAttendance

	startOfDay := today
	endOfDay := helper.GetEndOfDay(startOfDay)

	filter := bson.M{
		"user_id": userID,
		"date": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}

	err := r.collectionDailyAttendance.FindOne(c, filter).Decode(&existingDailyAttendance)

	if mongo.ErrNoDocuments == err {
		return nil, nil
	}

	return &existingDailyAttendance, nil

}

func (r *attendanceRepository) CreateAttendanceLog(c context.Context, attendanceLog *AttendanceLog) error {
	_, err := r.collectionAttendance.InsertOne(c, attendanceLog)
	return err
}

func (r *attendanceRepository) CreateDailyAttendance(c context.Context, dailyAttendance *DailyAttendance) error {
	_, err := r.collectionDailyAttendance.InsertOne(c, dailyAttendance)
	return err
}

func (r *attendanceRepository) UpdatedDailyAttendance(c context.Context, userID string, today time.Time, dailyAttendance *DailyAttendance) error {

	startOfDay := today
	endOfDay := helper.GetEndOfDay(startOfDay)

	filter := bson.M{
		"user_id": userID,
		"date": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"check_out_time":      dailyAttendance.CheckoutTime,
			"lunch_duration":      dailyAttendance.LunchDuration,
			"emotion_check_out":   dailyAttendance.EMotionCheckOut,
			"percent_work_day":    dailyAttendance.PercentWorkDay,
			"total_working_hours": dailyAttendance.TotalWorkingHours,
			"updated_at":          dailyAttendance.UpdatedAt,
		},
	}

	_, err := r.collectionDailyAttendance.UpdateOne(c, filter, update)

	return err
}

func (r *attendanceRepository) GetMyAttendance(c context.Context, userID string, firstDay time.Time, lastDay time.Time) ([]*DailyAttendance, error) {

	var dailyAttendances []*DailyAttendance

	filter := bson.M{
		"user_id": userID,
		"date": bson.M{
			"$gte": firstDay,
			"$lt":  lastDay,
		},
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"date": 1})

	cursor, err := r.collectionDailyAttendance.Find(c, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c)

	err = cursor.All(c, &dailyAttendances)
	if err != nil {
		return nil, err
	}

	if len(dailyAttendances) == 0 {
		return nil, fmt.Errorf("no attendance records found for user %s", userID)
	}

	return dailyAttendances, nil

}
