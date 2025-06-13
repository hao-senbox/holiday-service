package leave

import (
	"context"
	"fmt"
	"net/http"
	"worktime-service/helper"
	"worktime-service/pkg/constants"

	"github.com/gin-gonic/gin"
)

type LeaveHandler struct {
	leaveService LeaveService
}

func NewLeaveHandler(leaveService LeaveService) *LeaveHandler {
	return &LeaveHandler{
		leaveService: leaveService,
	}
}

func (h *LeaveHandler) GetMyRequest(c *gin.Context) {

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	userIDToken := userID.(string)

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}
	
	ctx := context.WithValue(c, constants.TokenKey, token)

	data, err := h.leaveService.GetMyRequest(ctx, userIDToken)
	
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Success", data)
}

func (h *LeaveHandler) GetAllLeaveCalendar(c *gin.Context) {

	date := c.Query("date")

	data, err := h.leaveService.GetAllLeaveCalendar(c, date)
	
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Success", data)
}

func (h *LeaveHandler) GetDetailLeaveCalendar(c *gin.Context) {

	id := c.Param("id")

	data, err := h.leaveService.GetDetailLeaveCalendar(c, id)
	
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Success", data)
} 

func (h *LeaveHandler) CreateRequestLeave(c *gin.Context) {

	var req CreateLeaveRequest
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

	err := h.leaveService.CreateRequestLeave(ctx, &req)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Success", nil)

}

func (h *LeaveHandler) GetSetting(c *gin.Context) {

	setting := h.leaveService.GetSettings(c)

	helper.SendSuccess(c, 200, "Success", setting)

}

func (h *LeaveHandler) UpdateSetting(c *gin.Context) {

	id := c.Param("id")

	var req SettingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	data, err := h.leaveService.UpdateSetting(c, &req, id)

	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Success", data)

}

func (h *LeaveHandler) EditMaxSlot(c *gin.Context) {

	id := c.Param("id")

	var req EditSlotRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err := h.leaveService.EditMaxSlot(c, &req, id)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", nil)

}

func (h *LeaveHandler) GetPendingRequest(c *gin.Context) {

	data, err := h.leaveService.GetPendingRequest(c)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}
 
	helper.SendSuccess(c, http.StatusOK, "Success", data)

}

func (h *LeaveHandler) DeleteRequestLeave(c *gin.Context) {

	var req DeleteLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
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

	err := h.leaveService.DeleteRequestLeave(ctx, &req)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", nil)

}

func (h *LeaveHandler) UpdateRequestLeave(c *gin.Context) {

	id := c.Param("id")

	var req UpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err := h.leaveService.UpdateRequestLeave(c, &req, id)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", nil)

}

func (h *LeaveHandler) GetStatistical (c *gin.Context) {

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	data, err := h.leaveService.GetStatistical(c, dateFrom, dateTo)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Success", data)

}

// func (h *LeaveHandler) GetLeaveBalanceUser (c *gin.Context) {

// 	id := c.Param("user-id")

// 	data, err := h.leaveService.GetLeaveBalanceUser(c, id)
// 	if err != nil {
// 		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
// 		return
// 	}

// 	helper.SendSuccess(c, http.StatusOK, "Success", data)
// }