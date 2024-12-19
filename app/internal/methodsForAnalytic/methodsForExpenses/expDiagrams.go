package methodsForExpenses

import (
	"encoding/json"
	"fmt"
	"math"
	"net/url"
)

// GenerateExpensePieChartURL создает URL для диаграммы расходов с фиксированными цветами и процентами
func GenerateExpensePieChartURL(categorySummary map[string]uint64) (string, error) {
	labels := []string{}
	values := []int{}
	var totalExpense uint64

	categoryColors := map[string]string{
		"Бытовые траты":       "#4A90E2", // синий
		"Регулярные платежи":  "#FF4D4D", // красный
		"Одежда":              "#FFD700", // желтый
		"Здоровье":            "#4CAF50", // зеленый
		"Досуг и образование": "#FFA500", // оранжевый
		"Инвестиции":          "#9B59B6", // сиреневый
		"Прочее":              "#D3D3D3", // Серый
	}

	// общий расход
	for _, value := range categorySummary {
		totalExpense += value
	}

	// преобразуем суммы в целые %%
	for category, value := range categorySummary {
		if value > 0 { // Учитываем только ненулевые категории
			labels = append(labels, category)
			percentage := int(math.Round((float64(value) / float64(totalExpense)) * 100)) // округление без дробей
			values = append(values, percentage)
		}
	}

	if len(labels) == 0 {
		return "", fmt.Errorf("нет данных для построения диаграммы")
	}

	// присваиваем цвета для категорий
	colors := []string{}
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
					"formatter": "function(value) { return value + '%'; }", // Формат отображения: проценты
				},
			},
		},
	}

	// гоним в JSON
	jsonData, err := json.Marshal(chartData)
	if err != nil {
		return "", fmt.Errorf("ошибка при создании данных для диаграммы: %v", err)
	}

	// компануем URL для сервиса с графиками - QuickChart
	baseURL := "https://quickchart.io/chart"
	params := url.Values{}
	params.Add("c", string(jsonData))

	return fmt.Sprintf("%s?%s", baseURL, params.Encode()), nil
}
