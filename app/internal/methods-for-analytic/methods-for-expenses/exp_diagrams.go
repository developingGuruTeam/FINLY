package methods_for_expenses

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/url"
)

func GenerateWeeklyExpensePieChartURL(categorySummary map[string]uint64) (string, error) {
	labels := []string{}
	values := []int{}
	colors := []string{}
	var totalExpense uint64

	if len(categorySummary) == 0 {
		return "", fmt.Errorf("Нет данных для построения диаграммы")
	}

	categoryColors := map[string]string{
		"Бытовые траты":       "#4A90E2", // синий
		"Регулярные платежи":  "#FF4D4D", // красный
		"Одежда":              "#FFD700", // желтый
		"Здоровье":            "#4CAF50", // зеленый
		"Досуг и образование": "#FFA500", // оранжевый
		"Инвестиции":          "#9B59B6", // сиреневый
		"Прочее":              "#D3D3D3", // серый
	}

	// Общий расход
	for _, value := range categorySummary {
		totalExpense += value
	}

	// Преобразуем суммы в проценты
	for category, value := range categorySummary {
		if value > 0 {
			labels = append(labels, category)
			percentage := int(math.Round((float64(value) / float64(totalExpense)) * 100))
			values = append(values, percentage)
		}
	}

	// Присваиваем цвета
	for _, category := range labels {
		if color, exists := categoryColors[category]; exists {
			colors = append(colors, color)
		} else {
			colors = append(colors, "#CCCCCC") // цвет по умолчанию
		}
	}

	chartData := map[string]interface{}{
		"type": "doughnut", // Тип графика
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
					"display": false, // Отключаем легенду
				},
			},
		},
	}

	// Преобразуем в JSON
	jsonData, err := json.Marshal(chartData)
	if err != nil {
		return "", fmt.Errorf("Ошибка при создании данных для диаграммы: %v", err)
	}

	// Составляем URL
	baseURL := "https://quickchart.io/chart"
	params := url.Values{}
	params.Add("c", string(jsonData))

	return fmt.Sprintf("%s?%s", baseURL, params.Encode()), nil
}

// создание диаграммы
func GenerateExpensePieChartURL(categorySummary map[string]uint64) (string, error) {
	labels := []string{}
	values := []int{}
	colors := []string{}
	var totalExpense uint64

	if len(categorySummary) == 0 {
		return "", fmt.Errorf("Нет данных для построения диаграммы")
	}

	categoryColors := map[string]string{
		"Бытовые траты":       "#4A90E2",
		"Регулярные платежи":  "#FF4D4D",
		"Одежда":              "#FFD700",
		"Здоровье":            "#4CAF50",
		"Досуг и образование": "#FFA500",
		"Инвестиции":          "#9B59B6",
		"Прочее":              "#D3D3D3",
	}

	for _, value := range categorySummary {
		totalExpense += value
	}

	log.Printf("Общий расход: %d", totalExpense)

	for category, value := range categorySummary {
		if value > 0 {
			labels = append(labels, category)
			percentage := int(math.Round((float64(value) / float64(totalExpense)) * 100))
			values = append(values, percentage)
		}
	}

	log.Printf("Labels: %+v", labels)
	log.Printf("Values: %+v", values)

	for _, category := range labels {
		if color, exists := categoryColors[category]; exists {
			colors = append(colors, color)
		} else {
			colors = append(colors, "#CCCCCC")
		}
	}

	chartData := map[string]interface{}{
		"type": "doughnut",
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
					"display": false,
				},
			},
		},
	}

	jsonData, err := json.Marshal(chartData)
	if err != nil {
		return "", fmt.Errorf("Ошибка при создании данных для диаграммы: %v", err)
	}

	log.Printf("Данные для диаграммы (JSON): %s", jsonData)

	baseURL := "https://quickchart.io/chart"
	params := url.Values{}
	params.Add("c", string(jsonData))

	finalURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	log.Printf("URL для диаграммы: %s", finalURL)

	return finalURL, nil
}
