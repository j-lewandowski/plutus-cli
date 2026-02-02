package sync

import (
	"plutus-cli/internal/db"
	"time"
)

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
		item, exists := itemMap[dateStr]

		if exists {
			result = append(result, item)
			lastKnown = item
			hasLastKnown = true
			continue
		}

		if !hasLastKnown {
			continue
		}

		newItem := lastKnown.CreateWithDate(day).(T)

		if ip, ok := any(newItem).(db.IndexPrice); ok {
			ip.IsReal = false
			newItem = any(ip).(T)
		}

		result = append(result, newItem)
	}

	return result
}

func DaysUntilToday(startDate time.Time) []time.Time {
	var days []time.Time

	now := time.Now()

	cursor := startDate

	for !cursor.After(now) {
		days = append(days, cursor)
		cursor = cursor.AddDate(0, 0, 1)
	}

	return days
}
