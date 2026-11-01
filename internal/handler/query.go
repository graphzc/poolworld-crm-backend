package handler

import (
	"strings"

	"poolworld-pattaya-backend/internal/entity"
)

var validDays = map[string]bool{
	string(entity.DayOfWeekSunday):    true,
	string(entity.DayOfWeekMonday):    true,
	string(entity.DayOfWeekTuesday):   true,
	string(entity.DayOfWeekWednesday): true,
	string(entity.DayOfWeekThursday):  true,
	string(entity.DayOfWeekFriday):    true,
	string(entity.DayOfWeekSaturday):  true,
}

// parseDays splits a comma-separated list of day codes (e.g. "MON,wed"),
// dropping unknown values and duplicates.
func parseDays(raw string) []string {
	var days []string
	seen := map[string]bool{}
	for _, part := range strings.Split(raw, ",") {
		day := strings.ToUpper(strings.TrimSpace(part))
		if validDays[day] && !seen[day] {
			seen[day] = true
			days = append(days, day)
		}
	}
	return days
}
