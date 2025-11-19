package shared

type GetStudentTemperatureChartRequest struct {
	TermID    string `json:"term_id" bson:"term_id"`
	StudentID string `json:"student_id" bson:"student_id"`
}
