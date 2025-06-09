package attendance

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, handler *AttendanceHandler) {

	attendanceGroup := r.Group("/api/v1/attendance")
	{
		attendanceGroup.POST("/checkin", handler.CheckIn)
		attendanceGroup.POST("/checkout", handler.CheckOut)
		attendanceGroup.GET("/my-attendance", handler.GetMyAttendance)
	}
}