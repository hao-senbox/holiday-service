package attendance

import (
	"worktime-service/helper"

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

	err := h.service.CheckIn(c, &req)
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

	err := h.service.CheckOut(c, &req)
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Success", nil)

}

func (h *AttendanceHandler) GetMyAttendance(c *gin.Context) {

	userID := c.Query("user-id")
	month := c.Query("month")
	year := c.Query("year")

	data, err := h.service.GetMyAttendance(c, userID, month, year)
	
	if err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, 200, "Success", data)

}