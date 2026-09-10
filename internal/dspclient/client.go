package dspclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mini-adex/internal/domain"
	"net/http"
	"time"
)

type HTTPClient struct {
	httpClient *http.Client
}

func NewHTTPClient(timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// Сериализует запрос и отправляет его на endpoint партнера (dsp)
func (c *HTTPClient) SendBidRequest(ctx context.Context, dsp domain.DSP, req domain.AuctionRequest) error {

	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("Ошибка при сериализации: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, dsp.Endpoint, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("Не удалось создать запрос: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("Не удалось выполнить запрос: %w", err)
	}

	defer resp.Body.Close()

	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("dsp вернул ответ %d, а не 200", resp.StatusCode)
	}

	return nil
}
