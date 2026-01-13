package actions

import "time"

type DatedItem interface {
	GetDate() time.Time
	CreateWithDate(date time.Time) interface{}
}

func PopulateMissingData[T DatedItem](allDays []time.Time, items []T) []T {
	result := make([]T, 0, len(allDays))
	itemMap := make(map[string]T)

	for _, item := range items {
		itemMap[item.GetDate().Format("2006-01-02")] = item
	}

	var lastKnown T
	var hasLastKnown bool

	for _, day := range allDays {
		dateStr := day.Format("2006-01-02")
		if item, exists := itemMap[dateStr]; exists {
			result = append(result, item)
			lastKnown = item
			hasLastKnown = true
		} else {
			if hasLastKnown {
				newItem := lastKnown.CreateWithDate(day).(T)
				result = append(result, newItem)
			}
		}
	}

	return result
}
