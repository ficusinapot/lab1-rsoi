package status

import "github.com/ficusinapot/ds/internal/models/entities"

func statusDataToResponse(statusData *entities.StatusData) GetStatusResponse {
	return GetStatusResponse{
		DBAvailable:     statusData.DBAvailable,
		OperatingStatus: string(statusData.OperatingStatus),
	}
}
