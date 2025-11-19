package response

type TermResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
