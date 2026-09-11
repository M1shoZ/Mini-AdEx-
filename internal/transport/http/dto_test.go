package http_test

import (
	"testing"

	transporthttp "mini-adex/internal/transport/http"
)

func TestAuctionRequestDTO_Validate(t *testing.T) {
	tests := []struct {
		name    string
		dto     transporthttp.AuctionRequestDTO
		wantErr bool
	}{
		{
			name: "валидный запрос",
			dto: transporthttp.AuctionRequestDTO{
				RequestID:  "req-1",
				Country:    "RU",
				DeviceType: "mobile",
				BidFloor:   1.0,
			},
			wantErr: false,
		},
		{
			name: "не хватает request_id",
			dto: transporthttp.AuctionRequestDTO{
				RequestID:  "",
				Country:    "RU",
				DeviceType: "mobile",
				BidFloor:   1.0,
			},
			wantErr: true,
		},
		{
			name: "неправильный формат кода страны",
			dto: transporthttp.AuctionRequestDTO{
				RequestID:  "req-1",
				Country:    "RUS", // Должно быть 2 символа
				DeviceType: "mobile",
				BidFloor:   1.0,
			},
			wantErr: true,
		},
		{
			name: "неизвестный device_type",
			dto: transporthttp.AuctionRequestDTO{
				RequestID:  "req-1",
				Country:    "RU",
				DeviceType: "smartwatch", // Не mobile/desktop/tv
				BidFloor:   1.0,
			},
			wantErr: true,
		},
		{
			name: "отрицательный bid_floor",
			dto: transporthttp.AuctionRequestDTO{
				RequestID:  "req-1",
				Country:    "RU",
				DeviceType: "mobile",
				BidFloor:   -0.5,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dto.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
