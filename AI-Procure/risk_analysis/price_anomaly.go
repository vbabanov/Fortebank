package risk_analysis

import (
	"ai_procure/models" // Для доступа к структурам ContractAmount и ContractSubject
	"fmt"
	"sort"
)

// PriceAnomalyResult содержит результат проверки цены
type PriceAnomalyResult struct {
	IsAnomaly bool
	Reason    string
	Boundary  float64 // Верхняя/Нижняя граница, которую цена нарушила
}

// CalculateQuartiles находит Q1 и Q3 в отсортированном слайсе цен.
func calculateQuartiles(prices []float64) (q1, q3 float64) {
	// В Go нужно вручную найти индексы Q1 (25%) и Q3 (75%).
	// (Детали реализации, включая интерполяцию, опускаем для концепта.)
	return prices[len(prices)/4], prices[len(prices)*3/4]
}

// CheckPriceAnomaly проверяет ContractAmount текущего тендера на наличие выбросов.
func CheckPriceAnomaly(currentPrice float64, segmentedHistory []models.HistoryRecord) PriceAnomalyResult {

	// 1. Сбор и сортировка исторических цен
	var prices []float64
	for _, record := range segmentedHistory {
		prices = append(prices, record.ContractAmount)
	}

	if len(prices) < 5 {
		// Недостаточно данных для статистического анализа
		return PriceAnomalyResult{IsAnomaly: false, Reason: "Недостаточно исторических данных для анализа цен."}
	}

	sort.Float64s(prices)

	// 2. Расчет IQR и Границ
	q1, q3 := calculateQuartiles(prices)
	iqr := q3 - q1

	// Множитель 1.5 для "мягких" выбросов
	upperBound := q3 + 1.5*iqr
	lowerBound := q1 - 1.5*iqr

	// 3. Правила Риска
	if currentPrice >= upperBound {
		return PriceAnomalyResult{
			IsAnomaly: true,
			Reason:    fmt.Sprintf("Высокий риск завышения цены. Текущая цена (%.2f) превышает верхнюю границу (%.2f) на %.2f%%.", currentPrice, upperBound, (currentPrice/upperBound-1)*100),
			Boundary:  upperBound,
		}
	}

	if currentPrice <= lowerBound {
		return PriceAnomalyResult{
			IsAnomaly: true,
			Reason:    fmt.Sprintf("Низкий риск (демпинг/нереалистичность). Текущая цена (%.2f) ниже нижней границы (%.2f).", currentPrice, lowerBound),
			Boundary:  lowerBound,
		}
	}

	return PriceAnomalyResult{IsAnomaly: false, Reason: "Цена соответствует историческому диапазону."}
}
