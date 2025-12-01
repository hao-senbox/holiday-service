package shared

type GetStudentTemperatureChartRequest struct {
	OrgID     string `json:"org_id" bson:"org_id"`
	StudentID string `json:"student_id" bson:"student_id"`
}
