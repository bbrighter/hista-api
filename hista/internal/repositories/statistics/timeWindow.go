package statistics

import "fmt"

type timeWindow struct {
	name      string
	fromHours int
	toHours   int
}

func newTimeWindows() []timeWindow {
	return []timeWindow{
		{name: "hours_72", fromHours: 24, toHours: 72},
		{name: "hours_24", fromHours: 1, toHours: 24},
		{name: "hours_1", fromHours: 0, toHours: 1},
	}
}

func buildTimewindowSelects() []string {
	windows := newTimeWindows()
	var selects []string
	for _, w := range windows {
		s := fmt.Sprintf(`MAX(
	CASE WHEN condition_events.date BETWEEN meals.date + interval '%d hour' AND meals.date + interval '%d hour' THEN 1 
	ELSE 0 END
	) as %s`, w.fromHours, w.toHours, w.name)
		selects = append(selects, s)
	}
	return selects
}

func buildTimeWindowSum() []string {
	windows := newTimeWindows()
	var selects []string
	for _, w := range windows {
		s := fmt.Sprintf("sum(u.%s) as %s", w.name, w.name)
		selects = append(selects, s)
	}
	return selects
}
