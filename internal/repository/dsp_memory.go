package repository

import (
	"context"
	"mini-adex/internal/domain"
)

type MemoryDSPRepo struct {
	dsps []domain.DSP
}

func NewMemoryDSPRepository() *MemoryDSPRepo {
	return &MemoryDSPRepo{
		dsps: []domain.DSP{
			{
				ID:                "11111111-1111-1111-1111-111111111111",
				Name:              "DSP First",
				Endpoint:          "http://127.0.0.1:9001/bid",
				IsEnabled:         true,
				Countries:         []string{"RU", "AM"},
				DeviceTypes:       []string{"mobile", "desktop"},
				MinBidFloor:       0.5,
				BlockedCategories: []string{"gambling"},
			},
			{
				ID:                "22222222-2222-2222-2222-222222222222",
				Name:              "DSP Second",
				Endpoint:          "http://127.0.0.1:9002/bid",
				IsEnabled:         true,
				Countries:         []string{"US", "DE"},
				DeviceTypes:       []string{"desktop"},
				MinBidFloor:       2.0,
				BlockedCategories: []string{"adult"},
			},
			{
				ID:                "33333333-3333-3333-3333-333333333333",
				Name:              "DSP Third",
				Endpoint:          "http://127.0.0.1:9003/bid",
				IsEnabled:         true,
				Countries:         nil, // Любые страны
				DeviceTypes:       nil, // Любые девайсы
				MinBidFloor:       0.1,
				BlockedCategories: nil,
			},
			{
				ID:                "44444444-4444-4444-4444-444444444444",
				Name:              "DSP Fourth",
				Endpoint:          "http://127.0.0.1:9004/bid",
				IsEnabled:         false, // Отключен
				Countries:         []string{"RU"},
				DeviceTypes:       []string{"mobile"},
				MinBidFloor:       0.1,
				BlockedCategories: nil,
			},
		},
	}
}

func (r *MemoryDSPRepo) GetAll(ctx context.Context) ([]domain.DSP, error) {
	result := make([]domain.DSP, len(r.dsps))
	copy(result, r.dsps)
	return result, nil
}
