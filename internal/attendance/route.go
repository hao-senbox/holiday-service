package attendance

import (
	"worktime-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *AttendanceHandler) {

	adminGroup := r.Group("/api/v1/admin")
	gatewayGroup := r.Group("/api/v1/gateway")
	gatewayGroup.Use(middleware.Secured())
	{
		gatewayGroup.GET("/student-temperature", handler.GetStudentTemperature)
	}
	
	adminGroup.Use(middleware.Secured())
	{
		attendanceGroup := adminGroup.Group("/attendance")
		{
			attendanceGroup.GET("/student-temperature-chart", handler.GetStudentTemperatureChart)
		}
	}

	attendanceGroup := r.Group("/api/v1/attendance").Use(middleware.Secured())
	{
		attendanceGroup.POST("/checkin", handler.CheckIn)
		attendanceGroup.POST("/checkout", handler.CheckOut)
		attendanceGroup.GET("/my-attendance", handler.GetMyAttendance)
		attendanceGroup.GET("", handler.GetAllAttendances)
		attendanceGroup.POST("/student", handler.AttendanceStudent)
		attendanceGroup.GET("/student", handler.GetMyAttendanceStudent)
	}

	
}
