package service

import (
	"context"
	"mini-adex/internal/domain"
)

type DSPRepo interface {
	GetAll(ctx context.Context) ([]domain.DSP, error)
}

type ClientRepo interface {
	SendBidRequest(ctx context.Context, dsp domain.DSP, req domain.AuctionRequest) error
}
