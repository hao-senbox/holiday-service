package attendance

import (
	"context"
	"fmt"
	"worktime-service/helper"
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

func (h *AttendanceHandler)  CheckIn(c *gin.Context) {
	
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

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	UserIDToken := userID.(string)

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}
	
	ctx := context.WithValue(c, constants.TokenKey, token)

	data, err := h.service.GetMyAttendance(ctx, UserIDToken, month, year)
	
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Success", data)

}