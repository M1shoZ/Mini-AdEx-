package service

import (
	"fmt"
	"mini-adex/internal/domain"
	"slices"
)

// TODO: искать заблокированные категории через мапу для увеличения эффективности

// Проверяет подходит ли DSP под параметры запроса
// При успешном сопоставлении возвращает true и пустую строку
// При НЕ успешном сопоставлении возвращает false и причину отказа
func MatchDSP(req domain.AuctionRequest, dsp domain.DSP) (bool, string) {
	// Проверка включен ли партнер
	if dsp.IsEnabled == false {
		return false, "Партнер не включен"
	}

	// Проверка гео
	if len(dsp.Countries) != 0 && !slices.Contains(dsp.Countries, req.Country) {
		return false, fmt.Sprintf("страна %q не поддерживается партнёром", req.Country)
	}

	// Проверка совместимости типа устройства
	if len(dsp.DeviceTypes) != 0 && !slices.Contains(dsp.DeviceTypes, req.DeviceType) {
		return false, fmt.Sprintf("Устройство типа %q не поддерживается партнёром", req.DeviceType)
	}

	// Проверка цены
	if req.BidFloor < dsp.MinBidFloor {
		return false, fmt.Sprintf("bid_floor %.2f меньше чем min_bid_floor партнёра %.2f", req.BidFloor, dsp.MinBidFloor)
	}

	// Проверка заблокированнных категорий
	for _, reqCategory := range req.Categories {
		if slices.Contains(dsp.BlockedCategories, reqCategory) {
			return false, fmt.Sprintf("Категория %q заблокирована партнёром", reqCategory)
		}
	}

	return true, ""
}

// Принимает список всех партнеров и возвращает только те, которые прошли претаргетинг,
// а также мапу причин отказа
func FiltreDSPs(req domain.AuctionRequest, dsps []domain.DSP) (matched []domain.DSP, rejectedReasons map[string]string) {
	rejectedReasons = make(map[string]string)
	for _, dsp := range dsps {
		matchedOk, reason := MatchDSP(req, dsp)
		if matchedOk {
			matched = append(matched, dsp)
		} else {
			rejectedReasons[dsp.Name] = reason
		}
	}
	return matched, rejectedReasons
}
