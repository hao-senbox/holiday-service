package leave

import (
	"worktime-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *LeaveHandler) {
	
	{	leaveGroup := r.Group("/api/v1/leave").Use(middleware.Secured())

		leaveGroup.POST("/", handler.CreateRequestLeave)
		leaveGroup.DELETE("/delete-request", handler.DeleteRequestLeave)
		leaveGroup.GET("/my-request", handler.GetMyRequest)
		leaveGroup.GET("pending-request", handler.GetPendingRequest)
		leaveGroup.PUT("/:id", handler.UpdateRequestLeave)
		leaveGroup.GET("/statistical", handler.GetStatistical)

		// leaveGroup.GET("leave-balance/:user-id", handler.GetLeaveBalanceUser)
		leaveGroup.GET("/calendar", handler.GetAllLeaveCalendar)
		leaveGroup.GET("/calendar/:id", handler.GetDetailLeaveCalendar)
		leaveGroup.PUT("/calendar/:id", handler.EditMaxSlot)

		leaveGroup.GET("/setting", handler.GetSetting)
		leaveGroup.PUT("/setting/:id", handler.UpdateSetting)
	}
}	