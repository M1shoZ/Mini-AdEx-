package service

import (
	"context"
	"log/slog"
	"mini-adex/internal/domain"
	"sync"
	"sync/atomic"
	"time"
)

type DSPRepo interface {
	GetAll(ctx context.Context) ([]domain.DSP, error)
}

type DSPClient interface {
	SendBidRequest(ctx context.Context, dsp domain.DSP, req domain.AuctionRequest) error
}

// Управляет проведением аукциона
type AuctionService struct {
	repo       DSPRepo
	client     DSPClient
	dspTimeout time.Duration
	logger     *slog.Logger
}

func NewAuctionService(repo DSPRepo, client DSPClient, dspTimeout time.Duration, logger *slog.Logger) *AuctionService {
	if logger == nil {
		logger = slog.Default()
	}
	return &AuctionService{
		repo:       repo,
		client:     client,
		dspTimeout: dspTimeout,
		logger:     logger,
	}
}

func (s *AuctionService) RunAuction(ctx context.Context, req domain.AuctionRequest) (domain.AuctionResult, error) {
	startTime := time.Now()

	allDSPs, err := s.repo.GetAll(ctx)
	if err != nil {
		return domain.AuctionResult{}, err
	}

	matchedDSPs, rejectedReasons := FiltreDSPs(req, allDSPs)

	s.logger.InfoContext(ctx, "pretargeted completed",
		slog.String("request_id", req.RequestID),
		slog.Int("total_dsps", len(allDSPs)),
		slog.Int("matched_dsds", len(matchedDSPs)),
		slog.Int("rejected_dsds", len(rejectedReasons)),
		slog.Any("rejected_reasons", rejectedReasons),
	)

	matchedNames := make([]string, 0, len(matchedDSPs))
	for _, dsp := range matchedDSPs {
		matchedNames = append(matchedNames, dsp.Name)
	}

	// если никто не подошел отдаем результат без сетевых вызовов
	if len(matchedDSPs) == 0 {
		return domain.AuctionResult{
			RequestID:   req.RequestID,
			MatchedDSPs: matchedNames,
			Sent:        0,
			Succeeded:   0,
			DurationMs:  time.Since(startTime).Milliseconds(),
		}, nil
	}

	// параллельная отправка запросов отобранным DSP (Fan-Out) с общим таймаутом
	auctionCtx, cancel := context.WithTimeout(ctx, s.dspTimeout)
	defer cancel()

	var (
		wg        sync.WaitGroup
		succeeded atomic.Int64
	)

	for _, dsp := range matchedDSPs {
		wg.Add(1)

		go func(targetDSP domain.DSP) {
			defer wg.Done()

			err := s.client.SendBidRequest(auctionCtx, targetDSP, req)
			if err != nil {
				//  Ошибку отдельного DSP логируем в Debug, чтобы не засорять логи
				s.logger.DebugContext(auctionCtx, "dsp request failed",
					slog.String("request_id", req.RequestID),
					slog.String("dsp_name", targetDSP.Name),
					slog.String("error", err.Error()),
				)
				return
			}
			succeeded.Add(1)
		}(dsp)
	}

	wg.Wait()

	duration := time.Since(startTime)

	result := domain.AuctionResult{
		RequestID:   req.RequestID,
		MatchedDSPs: matchedNames,
		Sent:        len(matchedDSPs),
		Succeeded:   int(succeeded.Load()),
		DurationMs:  duration.Milliseconds(),
	}

	s.logger.InfoContext(ctx, "auction finished",
		slog.String("request_id", req.RequestID),
		slog.Int("sent", result.Sent),
		slog.Int("succeeded", result.Succeeded),
		slog.Int64("duration_ms", result.DurationMs),
	)

	return result, nil
}
