package http

import (
	"errors"
	"fmt"
	"mini-adex/internal/domain"
	"slices"
)

var validDeviceTypes = []string{"mobile", "desktop", "tv"}

type AuctionRequestDTO struct {
	RequestID  string   `json:"request_id"`
	Country    string   `json:"country"`
	DeviceType string   `json:"device_type"`
	BidFloor   float64  `json:"bid_floor"`
	Categories []string `json:"categories"`
}

// Проверка корректности входящих данных
func (dto *AuctionRequestDTO) Validate() error {
	if dto.RequestID == "" {
		return errors.New("Поле RequestID обязательно для заполнения")
	}
	if len(dto.Country) != 2 {
		return errors.New("Поле Country должно соответствовать ISO 3166-1 alpha-2 (RU, US, DE)")
	}
	if !slices.Contains(validDeviceTypes, dto.DeviceType) {
		return fmt.Errorf("в поле DeviceType должно быть хотя бы одно из: %v", validDeviceTypes)
	}
	if dto.BidFloor < 0 {
		return errors.New("Поле BidFloor должно быть больше или равно 0")
	}
	return nil
}

func (dto *AuctionRequestDTO) ToDomain() domain.AuctionRequest {
	return domain.AuctionRequest{
		RequestID:  dto.RequestID,
		Country:    dto.Country,
		DeviceType: dto.DeviceType,
		BidFloor:   dto.BidFloor,
		Categories: dto.Categories,
	}
}

type AuctionResponseDTO struct {
	RequestID   string   `json:"request_id"`
	MatchedDSPs []string `json:"matched_dsps"`
	Sent        int      `json:"sent"`
	Succeeded   int      `json:"succeeded"`
	DurationMS  int64    `json:"duration_ms"`
}

func ToResponceDTO(res domain.AuctionResult) AuctionResponseDTO {
	matched := res.MatchedDSPs
	if matched == nil {
		matched = []string{}
	}

	return AuctionResponseDTO{
		RequestID:   res.RequestID,
		MatchedDSPs: res.MatchedDSPs,
		Sent:        res.Sent,
		Succeeded:   res.Succeeded,
		DurationMS:  res.DurationMs,
	}
}

type ErrorResponseDTO struct {
	Error string `json:"error"`
}
