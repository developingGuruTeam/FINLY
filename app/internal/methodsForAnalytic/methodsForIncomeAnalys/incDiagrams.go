package methodsForIncomeAnalys

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// создание диаграммы
func GenerateIncomePieChartURL(categorySummary map[string]uint64, totalIncome uint64) (string, error) {
	if len(categorySummary) == 0 {
		return "", fmt.Errorf("нет данных для построения диаграммы")
	}

	// цвета
	categoryColors := map[string]string{
		"Заработная плата":    "#4A90E2", // Синий
		"Побочный доход":      "#FF4D4D", // Красный
		"Доход от бизнеса":    "#FFD700", // Желтый
		"Гос. выплаты":        "#4CAF50", // Зеленый
		"Продажа имущества":   "#FFA500", // Оранжевый
		"Доход от инвестиций": "#9B59B6", // Сиреневый
		"Прочее":              "#A52A2A", // Коричневый
	}

	labels := []string{}
	values := []int{}
	colors := []string{}

	for category, value := range categorySummary {
		percentage := int((float64(value) / float64(totalIncome)) * 100)
		labels = append(labels, category)
		values = append(values, percentage)

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
					"formatter": "function(value) { return value + '%'; }", // Формат отображения: проценты
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
