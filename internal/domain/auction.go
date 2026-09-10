package domain

type AuctionRequest struct {
	RequestID  string   `json:"request_id"`
	Country    string   `json:"country"`
	DeviceType string   `json:"device_type"`
	BidFloor   float64  `json:"bid_floor"`
	Categories []string `json:"categories"`
}

type AuctionResult struct {
	RequestID   string   `json:"request_id"`
	MatchedDSPs []string `json:"matched_dsps"`
	Sent        int      `json:"sent"`
	Succeeded   int      `json:"succeeded"`
	DurationMs  int64    `json:"duration_ms"`
}
