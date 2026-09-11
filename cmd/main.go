package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"mini-adex/internal/dspclient"
	"mini-adex/internal/repository"
	"mini-adex/internal/service"
	transporthttp "mini-adex/internal/transport/http"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// конфигурация через флаги
	port := flag.Int("port", 8080, "HTTP server port")
	dspTimeout := flag.Duration("dsp-timeout", 200*time.Millisecond, "Таймаут ожидания ответа DSP партнёра")
	flag.Parse()

	// инициализируем логгер
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// запускаем мок сервера
	startMockDSPServers(ctx, logger)

	repo := repository.NewMemoryDSPRepository()
	client := dspclient.NewHTTPClient(*dspTimeout)
	auctionService := service.NewAuctionService(repo, client, *dspTimeout, logger)
	handler := transporthttp.NewHandler(auctionService, logger)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	//Запуск сервера в отдельной горутине
	go func() {
		logger.Info("запуск",
			slog.Int("port", *port),
			slog.Duration("dsp_timeout", *dspTimeout),
		)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("не удалось запустить сервер", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Ожидаем сигнал остановки
	<-ctx.Done()
	logger.Info("Graceful shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("сервер был остановлен принудительно", slog.String("err", err.Error()))
	} else {
		logger.Info("сервер остановлен успешно!")
	}

}

// запускает легковесные моки DSP
func startMockDSPServers(ctx context.Context, logger *slog.Logger) {
	runMockServer(ctx, "127.0.0.1:9001", 0, logger, "DSP First")
	runMockServer(ctx, "127.0.0.1:9002", 0, logger, "DSP Second")
	runMockServer(ctx, "127.0.0.1:9003", 350*time.Millisecond, logger, "DSP Third") // медленный
}

func runMockServer(ctx context.Context, addr string, delay time.Duration, logger *slog.Logger, name string) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /bid", func(w http.ResponseWriter, r *http.Request) {

		time.Sleep(delay)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status: bid_received"}`))
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Warn("мок сервер остановлен", slog.String("name", name), slog.String("addr", addr))
		}
	}()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()
}
