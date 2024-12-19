package methodsForIncomeAnalys

import (
	"encoding/json"
	"fmt"
	"math"
	"net/url"
)

// создание диаграммы
func GenerateIncomePieChartURL(categorySummary map[string]uint64, totalIncome uint64) (string, error) {
	labels := []string{}
	values := []int{}
	colors := []string{}

	if len(categorySummary) == 0 {
		return "", fmt.Errorf("Нет данных для построения диаграмм")
	}

	categoryColors := map[string]string{
		"Заработная плата":    "#4A90E2", // Синий
		"Побочный доход":      "#FF4D4D", // Красный
		"Доход от бизнеса":    "#FFD700", // Желтый
		"Гос. выплаты":        "#4CAF50", // Зеленый
		"Продажа имущества":   "#FFA500", // Оранжевый
		"Доход от инвестиций": "#9B59B6", // Сиреневый
		"Прочее":              "#A52A2A", // серый
	}

	// преобразуем суммы в целые %% учитывая только не нули
	for category, value := range categorySummary {
		if value > 0 {
			labels = append(labels, category)
			percentage := int(math.Round((float64(value) / float64(totalIncome)) * 100)) // округление без дробей
			values = append(values, percentage)
		}
	}

	// присваиваем цвета для категорий
	for _, category := range labels {
		if color, exists := categoryColors[category]; exists {
			colors = append(colors, color)
		} else {
			colors = append(colors, "#CCCCCC") // цвет по умолчанию
		}
	}
	// то, из чего получается диаграмма
	chartData := map[string]interface{}{
		"type": "doughnut", // тип
		"data": map[string]interface{}{
			"labels": labels,
			"datasets": []map[string]interface{}{
				{
					"data":            values,
					"backgroundColor": colors,
				},
			},
		},
		"options": map[string]interface{}{
			"plugins": map[string]interface{}{
				"legend": map[string]interface{}{
					"display": false, // Полностью отключаем легенду
				},
				"datalabels": map[string]interface{}{
					"formatter": "function(value) { return value + '%'; }",
				},
			},
		},
	}

	jsonData, err := json.Marshal(chartData)
	if err != nil {
		return "", fmt.Errorf("ошибка создания данных для графика: %v", err)
	}

	baseURL := "https://quickchart.io/chart"
	params := url.Values{}
	params.Add("c", string(jsonData))

	return fmt.Sprintf("%s?%s", baseURL, params.Encode()), nil
}
