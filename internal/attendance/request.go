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
	DurianLunch int    `json:"durian_lunch" bson:"durian_lunch"`
}
