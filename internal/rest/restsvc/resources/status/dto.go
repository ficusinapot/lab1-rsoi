package status

type StatusOutput struct {
	Body StatusResponse
}

type StatusResponse struct {
	Status string `json:"status"`
}

type GetStatusOutput struct {
	Body GetStatusResponse
}

type GetStatusResponse struct {
	DBAvailable     bool   `json:"db_available"`
	OperatingStatus string `json:"operating_status"`
}
