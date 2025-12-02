package attendance

import (
	"context"
	"fmt"
	"strconv"
	"worktime-service/helper"
	"worktime-service/internal/shared"
	"worktime-service/pkg/constants"

	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	service AttendanceService
}

func NewAttendanceHandler(service AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{
		service: service,
	}
}

func (h *AttendanceHandler) CheckIn(c *gin.Context) {

	var req CheckInRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	req.UserID = userID.(string)

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.service.CheckIn(ctx, &req)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Success", nil)

}

func (h *AttendanceHandler) CheckOut(c *gin.Context) {

	var req CheckOutRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	req.UserID = userID.(string)

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.service.CheckOut(ctx, &req)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Success", nil)

}

func (h *AttendanceHandler) GetMyAttendance(c *gin.Context) {

	month := c.Query("month")
	year := c.Query("year")
	userID := c.Query("user-id")

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	data, err := h.service.GetMyAttendance(ctx, userID, month, year)

	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Success", data)

}

func (h *AttendanceHandler) AttendanceStudent(c *gin.Context) {

	var req AttendanceStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	req.CreatedBy = userID.(string)

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.service.AttendanceStudent(ctx, &req)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Success", nil)

}

func (h *AttendanceHandler) GetMyAttendanceStudent(c *gin.Context) {

	month := c.Query("month")
	year := c.Query("year")
	userID := c.Query("user-id")

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	data, err := h.service.GetAttendanceStudent(ctx, userID, month, year)

	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Success", data)

}

func (h *AttendanceHandler) GetAllAttendances(c *gin.Context) {

	date := c.Query("date")
	userID := c.Query("user-id")
	page := c.Query("page")
	limit := c.Query("limit")

	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	if page == "" {
		pageInt = 1
	}
	if limit == "" {
		limitInt = 10
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	data, err := h.service.GetAllAttendances(ctx, userID, date, pageInt, limitInt)

	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Get All Attendances", data)

}

func (h *AttendanceHandler) GetStudentTemperature(c *gin.Context) {
	studentID := c.Query("student-id")
	if studentID == "" {
		helper.SendError(c, 400, fmt.Errorf("student_id is required"), helper.ErrInvalidRequest)
		return
	}
	
	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	res, err := h.service.GetStudentTemperature(ctx, studentID)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Success", res)
}

func (h *AttendanceHandler) GetStudentTemperatureChart(c *gin.Context) {
	studentID := c.Query("student-id")
	if studentID == "" {
		helper.SendError(c, 400, fmt.Errorf("student_id is required"), helper.ErrInvalidRequest)
		return
	}

	var req shared.GetStudentTemperatureChartRequest
	req.StudentID = studentID

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	res, err := h.service.GetStudentTemperatureChart(ctx, req)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Success", res)

}
