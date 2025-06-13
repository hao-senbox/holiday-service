package attendance

import (
	"worktime-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *AttendanceHandler) {

	attendanceGroup := r.Group("/api/v1/attendance").Use(middleware.Secured())
	{
		attendanceGroup.POST("/checkin", handler.CheckIn)
		attendanceGroup.POST("/checkout", handler.CheckOut)
		attendanceGroup.GET("/my-attendance", handler.GetMyAttendance)
	}
}