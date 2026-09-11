package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"mini-adex/internal/domain"
	"net/http"
)

type AuctionService interface {
	RunAuction(ctx context.Context, req domain.AuctionRequest) (domain.AuctionResult, error)
}

type Handler struct {
	service AuctionService
	logger  *slog.Logger
}

func NewHandler(service AuctionService, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// регистрируем эндпоинты в роутере стандартной библиотеки
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auction", h.handleAuction)
}

func (h *Handler) handleAuction(w http.ResponseWriter, r *http.Request) {
	var reqDTO AuctionRequestDTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&reqDTO); err != nil {
		h.writeError(w, http.StatusBadRequest, "некорректные данные: "+err.Error())
		return
	}

	if err := reqDTO.Validate(); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Запускаем аукцион через бизнес-слой
	result, err := h.service.RunAuction(r.Context(), reqDTO.ToDomain())
	if err != nil {
		// если клиент сам оборвал связь
		if errors.Is(err, context.Canceled) {
			return
		}
		h.logger.Error("Не удалось провести аукцион", slog.String("error", err.Error()))
		h.writeError(w, http.StatusInternalServerError, "Ошибка сервера")
		return
	}

	// Возвращаем ответ
	h.writeJSON(w, http.StatusOK, ToResponceDTO(result))

}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("не удалось записать ответ", slog.String("error", err.Error()))
	}
}

func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, ErrorResponseDTO{Error: message})
}
