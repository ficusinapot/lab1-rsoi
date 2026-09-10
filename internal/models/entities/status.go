package entities

type OperatingStatus string

const (
	FullyOperational OperatingStatus = "FullyOperational"
	NonOperational   OperatingStatus = "NonOperational"
)

type StatusData struct {
	DBAvailable     bool
	OperatingStatus OperatingStatus
}
