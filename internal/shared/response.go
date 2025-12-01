package shared

type StudentTemperatureChartResponse struct {
	Date         string    `json:"date"`
	Unit         string    `json:"unit"`
	Labels       string  `json:"labels"`
	Temperatures float64 `json:"temperatures"`
}
