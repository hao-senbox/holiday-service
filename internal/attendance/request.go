package attendance

type CheckInRequest struct {
	UserID  string `json:"user_id" bson:"user_id"`
	Emotion string `json:"emotion" bson:"emotion"`
	Notes   string `json:"notes" bson:"notes"`
}

type CheckOutRequest struct {
	UserID      string `json:"user_id" bson:"user_id"`
	Emotion     string `json:"emotion" bson:"emotion"`
	Notes       string `json:"notes" bson:"notes"`
	DurianLunch int    `json:"duration_lunch" bson:"duration_lunch"`
}

type AttendanceStudentRequest struct {
	UserID      string `json:"user_id" bson:"user_id"`
	Types       string `json:"types" bson:"types"`
	Notes       string `json:"note" bson:"note"`
	Temperature float64 `json:"temperature" bson:"temperature"`
	Date        string `json:"date" bson:"date"`
	CreatedBy   string `json:"created_by" bson:"created_by"`
}
