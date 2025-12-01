package shared

type StudentTemperatureChartResponse struct {
	Title        string    `json:"title"`
	Unit         string    `json:"unit"`
	Labels       []string  `json:"labels"`
	Temperatures []float64 `json:"temperatures"`
}
