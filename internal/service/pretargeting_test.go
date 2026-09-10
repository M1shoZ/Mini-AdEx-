package service_test

import (
	"mini-adex/internal/domain"
	"mini-adex/internal/service"
	"testing"
)

// TODO Заменить названия "провален" на другое

func TestMatchDSP(t *testing.T) {
	testDSP := domain.DSP{
		ID:                "test-dsp-001",
		Name:              "TestDSP",
		IsEnabled:         true,
		Countries:         []string{"RU", "AM"},
		DeviceTypes:       []string{"mobile", "desktop"},
		MinBidFloor:       1.0,
		BlockedCategories: []string{"gambling", "adult"},
	}

	testReq := domain.AuctionRequest{
		RequestID:  "test-req-001",
		Country:    "RU",
		DeviceType: "mobile",
		BidFloor:   1.5,
		Categories: []string{"news", "sport"},
	}

	tests := []struct {
		name        string
		req         domain.AuctionRequest
		dsp         domain.DSP
		wantMatched bool ``
	}{
		{
			name:        "успешный: все условия выполнены",
			req:         testReq,
			dsp:         testDSP,
			wantMatched: true,
		},
		{
			name: "успешный: список стран и типов устройств",
			req:  testReq,
			dsp: domain.DSP{
				IsEnabled:         true,
				Countries:         nil, // Любая страна
				DeviceTypes:       nil, // Любое устройство
				MinBidFloor:       1.0,
				BlockedCategories: nil,
			},
			wantMatched: true,
		},
		{
			name: "успешный: bid_floor равен min_bid_floor партнёра",
			req: func() domain.AuctionRequest {
				r := testReq
				r.BidFloor = 1.0
				return r
			}(),
			dsp:         testDSP,
			wantMatched: true,
		},
		{
			name: "провален: DSP не включен",
			req:  testReq,
			dsp: func() domain.DSP {
				d := testDSP
				d.IsEnabled = false
				return d
			}(),
			wantMatched: false,
		},
		{
			name: "провален: страна недоступна",
			req: func() domain.AuctionRequest {
				r := testReq
				r.Country = "US"
				return r
			}(),
			dsp:         testDSP,
			wantMatched: false,
		},
		{
			name: "провален: тип устройства недоступен",
			req: func() domain.AuctionRequest {
				r := testReq
				r.DeviceType = "tv"
				return r
			}(),
			dsp:         testDSP,
			wantMatched: false,
		},
		{
			name: "провален: bid_floor меньше чем min_bid_floor партнёра",
			req: func() domain.AuctionRequest {
				r := testReq
				r.BidFloor = 0.5
				return r
			}(),
			dsp:         testDSP,
			wantMatched: false,
		},
		{
			name: "провален: среди категорий есть заблокированная",
			req: func() domain.AuctionRequest {
				r := testReq
				r.Categories = []string{"gambling", "news"} // "gambling" заблокирован у партнера
				return r
			}(),
			dsp:         testDSP,
			wantMatched: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatched, reason := service.MatchDSP(tt.req, tt.dsp)
			if gotMatched != tt.wantMatched {
				t.Errorf("MatchDSP() matched = %v, want %v; reason: %s", gotMatched, tt.wantMatched, reason)
			}
		})
	}
}
