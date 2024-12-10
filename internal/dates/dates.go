package dates

import (
	"fmt"
	"hinoki-cli/internal/goal"
	"strconv"
	"strings"
	"time"
)

type Quarter int

const (
	Q1 Quarter = iota
	Q2
	Q3
	Q4
)

type DateKeyword = string

const (
	Today          DateKeyword = "today"
	TodayShort     DateKeyword = "t"
	Yesterday      DateKeyword = "yesterday"
	YesterdayShort DateKeyword = "ytd"
	Tomorrow       DateKeyword = "tomorrow"
	TomorrowShort  DateKeyword = "tmrw"
	Weekend        DateKeyword = "weekend"
	WeekendShort   DateKeyword = "wknd"
	Monday         DateKeyword = "monday"
	MondayShort    DateKeyword = "mon"
	Tuesday        DateKeyword = "tuesday"
	TuesdayShort   DateKeyword = "tue"
	Wednesday      DateKeyword = "wednesday"
	WednesdayShort DateKeyword = "wed"
	Thursday       DateKeyword = "thursday"
	ThursdayShort  DateKeyword = "thu"
	Friday         DateKeyword = "friday"
	FridayShort    DateKeyword = "fri"
	Saturday       DateKeyword = "saturday"
	SaturdayShort  DateKeyword = "sat"
	Sunday         DateKeyword = "sunday"
	SundayShort    DateKeyword = "sun"
	January        DateKeyword = "january"
	JanuaryShort   DateKeyword = "jan"
	February       DateKeyword = "february"
	FebruaryShort  DateKeyword = "feb"
	March          DateKeyword = "march"
	MarchShort     DateKeyword = "mar"
	April          DateKeyword = "april"
	AprilShort     DateKeyword = "apr"
	May            DateKeyword = "may"
	June           DateKeyword = "june"
	JuneShort      DateKeyword = "jun"
	July           DateKeyword = "july"
	JulyShort      DateKeyword = "jul"
	August         DateKeyword = "august"
	AugustShort    DateKeyword = "aug"
	September      DateKeyword = "september"
	SeptemberShort DateKeyword = "sep"
	October        DateKeyword = "october"
	OctoberShort   DateKeyword = "oct"
	November       DateKeyword = "november"
	NovemberShort  DateKeyword = "nov"
	December       DateKeyword = "december"
	DecemberShort  DateKeyword = "dec"
	Day            DateKeyword = "day"
	DayShort       DateKeyword = "d"
	Week           DateKeyword = "week"
	WeekShort      DateKeyword = "w"
	Month          DateKeyword = "month"
	MonthShort     DateKeyword = "m"
	Year           DateKeyword = "year"
	YearShort      DateKeyword = "y"
	QuarterKwrd    DateKeyword = "quarter"
	QuarterShort   DateKeyword = "q"
	Q1Kwrd         DateKeyword = "Q1"
	Q2Kwrd         DateKeyword = "Q2"
	Q3Kwrd         DateKeyword = "Q3"
	Q4Kwrd         DateKeyword = "Q4"
	Life           DateKeyword = "life"
	LifeShort      DateKeyword = "l"
)

func TimeframeDateString(t time.Time) string {
	return t.Format("2006-01-02")
}

func StartOfWeek(date time.Time) time.Time {
	// Calculate the number of days since Monday
	offset := int(date.Weekday()) - int(time.Monday)
	if offset < 0 {
		offset += 7 // Adjust for Sunday being the 0th day
	}

	// Subtract offset days from the given date
	return date.AddDate(0, 0, -offset).Truncate(24 * time.Hour)
}

func EndOfWeek(date time.Time) time.Time {
	// Calculate the number of days since Monday
	offset := int(time.Sunday) - int(date.Weekday())
	if offset < 0 {
		offset += 7 // Adjust for Sunday being the 0th day
	}

	// Subtract offset days from the given date
	return date.AddDate(0, 0, offset).Truncate(24 * time.Hour)
}

func StartOfQuarter(date time.Time) time.Time {
	switch date.Month() {
	case 1, 2, 3:
		return date.AddDate(0, 1-int(date.Month()), 0)
	case 4, 5, 6:
		return date.AddDate(0, 4-int(date.Month()), 0)
	case 7, 8, 9:
		return date.AddDate(0, 7-int(date.Month()), 0)
	case 10, 11, 12:
		return date.AddDate(0, 10-int(date.Month()), 0)
	}

	return date
}

func EndOfQuarter(date time.Time) time.Time {
	year := date.Year()
	month := date.Month()

	var quarterEndMonth time.Month
	var day int

	switch {
	case month >= 1 && month <= 3: // Q1
		quarterEndMonth = 3
		day = 31
	case month >= 4 && month <= 6: // Q2
		quarterEndMonth = 6
		day = 30
	case month >= 7 && month <= 9: // Q3
		quarterEndMonth = 9
		day = 30
	case month >= 10 && month <= 12: // Q4
		quarterEndMonth = 12
		day = 31
	}

	return time.Date(year, quarterEndMonth, day, 23, 59, 59, 0, date.Location())
}

func NextWeekday(date time.Time, targetWeekday time.Weekday) time.Time {
	daysUntilTarget := (int(targetWeekday) - int(date.Weekday()) + 7) % 7
	return date.AddDate(0, 0, daysUntilTarget)
}

func PrevWeekday(date time.Time, targetWeekday time.Weekday) time.Time {
	daysUntilTarget := (int(targetWeekday) - int(date.Weekday()) - 7) % 7
	return date.AddDate(0, 0, daysUntilTarget-7)
}

func CurrentWeekday(date time.Time, targetWeekday time.Weekday) time.Time {
	return date.AddDate(0, 0, int(targetWeekday)-int(date.Weekday()))
}

func NextMonth(date time.Time, targetMonth time.Month) time.Time {
	monthsUntilTarget := (int(targetMonth) - int(date.Month()) + 12) % 12
	return date.AddDate(0, monthsUntilTarget, 0)
}

func PrevMonth(date time.Time, targetMonth time.Month) time.Time {
	monthsUntilTarget := (int(targetMonth) - int(date.Month()) - 12) % 12
	return date.AddDate(0, monthsUntilTarget-12, 0)
}

func CurrentMonth(date time.Time, targetMonth time.Month) time.Time {
	return time.Date(date.Year(), targetMonth, date.Day(), 0, 0, 0, 0, date.Location())
}

func QuarterByNumber(date time.Time, q Quarter) time.Time {
	switch q {
	case Q1:
		return time.Date(date.Year(), 1, date.Day(), 0, 0, 0, 0, date.Location())
	case Q2:
		return time.Date(date.Year(), 4, date.Day(), 0, 0, 0, 0, date.Location())
	case Q3:
		return time.Date(date.Year(), 7, date.Day(), 0, 0, 0, 0, date.Location())
	case Q4:
		return time.Date(date.Year(), 10, date.Day(), 0, 0, 0, 0, date.Location())
	}

	return date
}

func weekdayKeywordToWeekday(keyword DateKeyword) (time.Weekday, bool) {
	switch keyword {
	case Monday, MondayShort:
		return time.Monday, true
	case Tuesday, TuesdayShort:
		return time.Tuesday, true
	case Wednesday, WednesdayShort:
		return time.Wednesday, true
	case Thursday, ThursdayShort:
		return time.Thursday, true
	case Friday, FridayShort:
		return time.Friday, true
	case Saturday, SaturdayShort:
		return time.Saturday, true
	case Sunday, SundayShort:
		return time.Sunday, true
	default:
		return time.Monday, false
	}
}

func monthKeywordToMonth(keyword DateKeyword) (time.Month, bool) {
	switch keyword {
	case January, JanuaryShort:
		return time.January, true
	case February, FebruaryShort:
		return time.February, true
	case March, MarchShort:
		return time.March, true
	case April, AprilShort:
		return time.April, true
	case May:
		return time.May, true
	case June, JuneShort:
		return time.June, true
	case July, JulyShort:
		return time.July, true
	case August, AugustShort:
		return time.August, true
	case September, SeptemberShort:
		return time.September, true
	case October, OctoberShort:
		return time.October, true
	case November, NovemberShort:
		return time.November, true
	case December, DecemberShort:
		return time.December, true
	default:
		return time.Month(0), false
	}
}

func dayKeywordToDayNumber(keyword DateKeyword) (int, bool) {
	if num, err := toInt(keyword); err == nil && num > 0 && num <= 31 {
		return num, true
	}
	return 0, false
}

func yearKeywordToYearNumber(keyword DateKeyword) (int, bool) {
	if num, err := toInt(keyword); err == nil && num > 0 && num > 1900 {
		return num, true
	}
	return 0, false
}

func quarterKeywordToQuarterNumber(keyword DateKeyword) (Quarter, bool) {
	switch keyword {
	case Q1Kwrd:
		return Q1, true
	case Q2Kwrd:
		return Q2, true
	case Q3Kwrd:
		return Q3, true
	case Q4Kwrd:
		return Q4, true
	default:
		return 0, false
	}
}

func ParseDate(current time.Time, date string) (time.Time, goal.Timeframe, error) {
	date = strings.ToLower(date)
	parts := strings.Split(date, " ")

	if len(parts) == 0 {
		return time.Now(), goal.Day, fmt.Errorf("invalid date: %s", date)
	}

	direction := 0

	if parts[0] == "next" || parts[0] == "n" {
		direction = 1
		date = parts[1]
	}

	if parts[0] == "prev" || parts[0] == "p" {
		direction = -1
		date = parts[1]
	}

	if len(parts) == 2 {
		day, isDay := dayKeywordToDayNumber(parts[0])
		month, isMonth := monthKeywordToMonth(parts[1])

		if isDay && isMonth {
			return time.Date(current.Year(), month, day, 0, 0, 0, 0, current.Location()), goal.Day, nil
		}
	}

	if len(parts) == 2 {
		month, isMonth := monthKeywordToMonth(parts[0])
		year, isYear := yearKeywordToYearNumber(parts[1])

		if isYear && isMonth {
			return time.Date(year, month, 1, 0, 0, 0, 0, current.Location()), goal.Month, nil
		}
	}

	if len(parts) == 3 {
		day, isDay := dayKeywordToDayNumber(parts[0])
		month, isMonth := monthKeywordToMonth(parts[1])
		year, isYear := yearKeywordToYearNumber(parts[2])

		if isYear && isMonth && isDay {
			return time.Date(year, month, day, 0, 0, 0, 0, current.Location()), goal.Day, nil
		}
	}

	weekdayFn := CurrentWeekday
	monthFn := CurrentMonth
	if direction > 0 {
		weekdayFn = NextWeekday
		monthFn = NextMonth
	}
	if direction < 0 {
		weekdayFn = PrevWeekday
		monthFn = PrevMonth
	}

	if day, ok := dayKeywordToDayNumber(date); ok {
		current = current.AddDate(0, direction, 0)
		date := time.Date(current.Year(), current.Month(), day, 0, 0, 0, 0, current.Location())
		return date, goal.Day, nil
	}

	if weekday, ok := weekdayKeywordToWeekday(date); ok {
		return weekdayFn(current, weekday), goal.Day, nil
	}

	if month, ok := monthKeywordToMonth(date); ok {
		return monthFn(current, month), goal.Month, nil
	}

	if year, ok := yearKeywordToYearNumber(date); ok {
		date := time.Date(year, current.Month(), current.Day(), 0, 0, 0, 0, current.Location())
		return date, goal.Year, nil
	}

	if quarter, ok := quarterKeywordToQuarterNumber(date); ok {
		return QuarterByNumber(current, quarter).AddDate(direction, 0, 0), goal.Quarter, nil
	}

	switch date {
	case Today, TodayShort:
		return current, goal.Day, nil
	case Day, DayShort:
		return current.AddDate(0, 0, 1+direction), goal.Day, nil
	case Yesterday, YesterdayShort:
		return current.AddDate(0, 0, -1), goal.Day, nil
	case Tomorrow, TomorrowShort:
		return current.AddDate(0, 0, 1), goal.Day, nil
	case Weekend, WeekendShort:
		return current.AddDate(0, 0, int(time.Saturday)-int(current.Weekday())+7*direction), goal.Day, nil
	case Week, WeekShort:
		return current.AddDate(0, 0, 7*direction), goal.Week, nil
	case Month, MonthShort:
		return current.AddDate(0, direction, 0), goal.Month, nil
	case QuarterKwrd, QuarterShort:
		return current.AddDate(0, 3*direction, 0), goal.Quarter, nil
	case Year, YearShort:
		return current.AddDate(direction, 0, 0), goal.Year, nil
	case Life, LifeShort:
		return current, goal.Life, nil
	}

	return time.Now(), goal.Day, fmt.Errorf("invalid date: %s", date)
}

func toInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func japaneseDayWeek(date time.Time) string {
	switch date.Weekday() {
	case time.Monday:
		return "月"
	case time.Tuesday:
		return "火"
	case time.Wednesday:
		return "水"
	case time.Thursday:
		return "木"
	case time.Friday:
		return "金"
	case time.Saturday:
		return "土"
	case time.Sunday:
		return "日"
	}

	return ""
}

func englishDayWeek(date time.Time) string {
	switch date.Weekday() {
	case time.Monday:
		return "Mon"
	case time.Tuesday:
		return "Tue"
	case time.Wednesday:
		return "Wed"
	case time.Thursday:
		return "Thu"
	case time.Friday:
		return "Fri"
	case time.Saturday:
		return "Sat"
	case time.Sunday:
		return "Sun"
	}
	return ""
}

func DateString(t time.Time, timeslice goal.Timeframe) string {
	switch timeslice {
	case goal.Day:
		return fmt.Sprintf("%s (%s)", t.Format("2 January 2006"), englishDayWeek(t))
	case goal.Week:
		_, week := t.ISOWeek()
		return fmt.Sprintf("%s – %s %s (%d)", StartOfWeek(t).Format("02"), EndOfWeek(t).Format("02"), t.Format("January 2006"), week)
	case goal.Month:
		return t.Format("January 2006")
	case goal.Quarter:
		return fmt.Sprintf("Q%d %d", int(t.Month()-1)/3+1, t.Year())
	case goal.Year:
		return t.Format("2006")
	case goal.Life:
		//return "is what happens when you’re busy making other plans"
		return ""
	}
	return t.Format("02 January 2006")
}

func ChangePeriod(t time.Time, timeframe goal.Timeframe, by int) time.Time {
	switch timeframe {
	case goal.Day:
		return t.AddDate(0, 0, by)
	case goal.Week:
		return t.AddDate(0, 0, 7*by)
	case goal.Month:
		return t.AddDate(0, by, 0)
	case goal.Quarter:
		return t.AddDate(0, 3*by, 0)
	case goal.Year:
		return t.AddDate(by, 0, 0)
	}
	return t
}
