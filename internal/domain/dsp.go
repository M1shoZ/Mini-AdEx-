package domain

type DSP struct {
	ID                string
	Name              string
	Endpoint          string
	IsEnabled         bool
	Countries         []string
	DeviceTypes       []string
	MinBidFloor       float64
	BlockedCategories []string
}
