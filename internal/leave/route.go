package leave

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, handler *LeaveHandler) {
	
	leaveGroup := r.Group("/api/v1/leave")
	{
		leaveGroup.POST("/", handler.CreateRequestLeave)
		leaveGroup.POST("/delete-request", handler.DeleteRequestLeave)
		leaveGroup.GET("/my-request/:user-id", handler.GetMyRequest)
		leaveGroup.GET("pending-request", handler.GetPendingRequest)
		leaveGroup.PUT("/:id", handler.UpdateRequestLeave)

		leaveGroup.GET("/calendar", handler.GetAllLeaveCalendar)
		leaveGroup.GET("/calendar/:id", handler.GetDetailLeaveCalendar)
		leaveGroup.PUT("/calendar/:id", handler.EditMaxSlot)
		
		leaveGroup.GET("/setting", handler.GetSetting)
		leaveGroup.PUT("/setting/:id", handler.UpdateSetting)
	}
}	