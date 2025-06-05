package leave

type LeaveStatistical struct {
	TotalRequested         int            `bson:"total_requested" json:"total_requested"`
	TotalConfirmed         int            `bson:"total_confirmed" json:"total_confirmed"`
	TotalPending           int            `bson:"total_pending" json:"total_pending"`
	TotalRejected          int            `bson:"total_rejected" json:"total_rejected"`
	ImmediateRequests      int            `bson:"immediate_requests" json:"immediate_requests"`
	ApproveRate            float64        `bson:"approve_rate" json:"approve_rate"`
	RequestsByMonth        []MonthlyStats `bson:"requests_by_month" json:"requests_by_month"`
	RequestByWeekDay       []WeeklyStats  `bson:"request_by_week_day" json:"request_by_week_day"`
	TopRequestUsers        []UserStats    `bson:"top_request_users" json:"top_request_users"`
	AverageRequestsPerUser float64        `bson:"average_requests_per_user" json:"average_requests_per_user"`
}

type MonthlyStats struct {
	Month    string `json:"month"`
	Count    int    `json:"count"`
	Approved int    `json:"approved"`
	Rejected int    `json:"rejected"`
}

type WeeklyStats struct {
	Weekday string `json:"weekday"`
	Count   int    `json:"count"`
}

type UserStats struct {
	UserID           string  `json:"user_id"`
	UserName         string  `json:"user_name"`
	TotalRequests    int     `json:"total_requests"`
	ApprovedRequests int     `json:"approved_requests"`
	ApprovalRate     float64 `json:"approval_rate"`
}