package risk_analysis

import (
	"AI-Procure/models"
	"AI-Procure/utils"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	NameSimilarityThreshold   = 0.85
	AdressSimilarityThreshold = 0.90
)

func CalculateFuzzyScore(s1, s2 string) float64 {
	s1 = strings.TrimSpace(s1)
	s2 = strings.TrimSpace(s2)
	if len(s1) == 0 || len(s2) == 0 {
		return 0.0
	}
	if strings.Contains(s1, s2) || strings.Contains(s2, s1) {
		return 0.95
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Float64() * 0.2
}

type AffiliationCheckResult struct {
	IsAffiliated bool
	Reason       string
	Score        float64
}

func CheckAffiliation(tenderer models.CompanyProfile, partnerHistory []models.HistoryRecord) AffiliationCheckResult {
	for _, history := range partnerHistory {
		if tenderer.BIN_IIN != "" && tenderer.BIN_IIN == history.Company.BIN_IIN {
			return AffiliationCheckResult{
				IsAffiliated: true,
				Reason:       fmt.Sprintf("Точное совпадение БИН/ИИН с компанией %s. Аффилированность подтверждена.", history.Company.Name),
				Score:        1.0,
			}
		}
	}

	for _, history := range partnerHistory {
		normNameCurrent := utils.NormalizeString(tenderer.Name)
		normNameHistory := utils.NormalizeString(history.Company.Name)
		normAddressCurrent := utils.NormalizeString(tenderer.Adress)
		normAddressHistory := utils.NormalizeString(history.Company.Adress)

		// 1  сходство названий

		nameScore := CalculateFuzzyScore(normNameCurrent, normNameHistory)
		if nameScore >= NameSimilarityThreshold {
			return AffiliationCheckResult{
				IsAffiliated: true,
				Reason:       fmt.Sprintf("Нечеткое совпадение названий: %s b %s. Сходство %.2f%%. Высокий риск аффилированности/смены названия.", tenderer.Name, history.Company.Name, nameScore*100),
				Score:        nameScore,
			}
		}

		// 2  общий адрес
		addressScore := CalculateFuzzyScore(normAddressCurrent, normAddressHistory)
		if addressScore >= AdressSimilarityThreshold {
			return AffiliationCheckResult{
				IsAffiliated: true,
				Reason:       fmt.Sprintf("Общий адрес: %.2f%% сходство адресовю Риск аффилированости (один офис)ю", addressScore*100),
				Score:        addressScore,
			}
		}

		// 3  негативная история
		if history.IsSuspicious && (nameScore >= 0.75 || addressScore >= 0.80) {
			return AffiliationCheckResult{
				IsAffiliated: true,
				Reason:       fmt.Sprintf("Подозрительная связь: Текущий поставщик схож с компанией %s, которая замешана в подозрительном контрактею Сходство имен/фдресов > 75% %.", history.Company.Name),
				Score:        0, 95,
			}
		}
	}

	return AffiliationCheckResult{
		IsAffiliated: false,
		Reason:       "Аффилированность не обнаружена.",
		Score:        0.0,
	}
}
