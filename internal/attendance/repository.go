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
	existingDailyAttendanceStudent(c context.Context, userID string, today time.Time) (*AttendanceStudent, error)
	CreateDailyAttendanceStudent(c context.Context, dailyAttendanceStudent *AttendanceStudent) error
	UpdateDailyAttendanceStudent(c context.Context, dailyAttendanceStudent *AttendanceStudent) error
	GetAttendanceStudent(c context.Context, userID string, firstDay time.Time, lastDay time.Time) ([]*AttendanceStudent, error)
	GetAllAttendances(c context.Context, userID string, date *time.Time, page int, limit int) ([]*DailyAttendance, int64, error)
}

type attendanceRepository struct {
	collectionAttendance             *mongo.Collection
	collectionDailyAttendance        *mongo.Collection
	collectionDailyAttendanceStudent *mongo.Collection
}

func NewAttendanceRepository(collectionAttendance *mongo.Collection, collectionDailyAttendance *mongo.Collection, collectionDailyAttendanceStudent *mongo.Collection) AttendanceRepository {
	return &attendanceRepository{
		collectionAttendance:             collectionAttendance,
		collectionDailyAttendance:        collectionDailyAttendance,
		collectionDailyAttendanceStudent: collectionDailyAttendanceStudent,
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
	fmt.Printf("userID: %s, firstDay: %s, lastDay: %s\n", userID, firstDay, lastDay)
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

func (r *attendanceRepository) existingDailyAttendanceStudent(c context.Context, userID string, today time.Time) (*AttendanceStudent, error) {

	var existingDailyAttendanceStudent AttendanceStudent

	startOfDay := today
	endOfDay := helper.GetEndOfDay(startOfDay)

	filter := bson.M{
		"user_id": userID,
		"date": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}

	err := r.collectionDailyAttendanceStudent.FindOne(c, filter).Decode(&existingDailyAttendanceStudent)

	if mongo.ErrNoDocuments == err {
		return nil, nil
	}

	return &existingDailyAttendanceStudent, nil

}

func (r *attendanceRepository) CreateDailyAttendanceStudent(c context.Context, dailyAttendanceStudent *AttendanceStudent) error {
	_, err := r.collectionDailyAttendanceStudent.InsertOne(c, dailyAttendanceStudent)
	return err
}

func (r *attendanceRepository) UpdateDailyAttendanceStudent(c context.Context, dailyAttendanceStudent *AttendanceStudent) error {
	_, err := r.collectionDailyAttendanceStudent.UpdateOne(c, bson.M{"_id": dailyAttendanceStudent.ID}, bson.M{"$set": dailyAttendanceStudent})
	return err
}

func (r *attendanceRepository) GetAttendanceStudent(c context.Context, userID string, firstDay time.Time, lastDay time.Time) ([]*AttendanceStudent, error) {

	var dailyAttendances []*AttendanceStudent
	fmt.Printf("userID: %s, firstDay: %s, lastDay: %s\n", userID, firstDay, lastDay)
	filter := bson.M{
		"user_id": userID,
		"date": bson.M{
			"$gte": firstDay,
			"$lt":  lastDay,
		},
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"date": 1})

	cursor, err := r.collectionDailyAttendanceStudent.Find(c, filter, findOptions)
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

func (r *attendanceRepository) GetAllAttendances(c context.Context, userID string, date *time.Time, page int, limit int) ([]*DailyAttendance, int64, error) {

	filter := bson.M{}

	if userID != "" {
		filter["user_id"] = userID
	}

	if date != nil {
		startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)

		filter["created_at"] = bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		}
	}

	totalCount, err := r.collectionDailyAttendance.CountDocuments(c, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := int64((page - 1) * limit)

	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collectionDailyAttendance.Find(c, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(c)

	var dailyAttendances []*DailyAttendance
	if err := cursor.All(c, &dailyAttendances); err != nil {
		return nil, 0, err
	}

	return dailyAttendances, totalCount, nil

}
