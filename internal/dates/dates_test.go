package dates

import (
	"hinoki-cli/internal/goal"
	"testing"
	"time"
)

func TestParseDate_ThuToFri(t *testing.T) {
	layout := "2006-01-02"
	date, _ := time.Parse(layout, "2024-11-21") // Thu
	result, timeframe, err := ParseDate(date, "Fri")
	expected := "2024-11-22"
	expectedTimeframe := goal.Day

	resultStr := result.Format(layout)
	if err == nil && (timeframe != expectedTimeframe || resultStr != expected) {
		t.Errorf("res %s; want %s", resultStr, expected)
		t.Errorf("res %s; want %s", timeframe, expectedTimeframe)
	}
}

func TestParseDate_ThuToThu(t *testing.T) {
	layout := "2006-01-02"
	date, _ := time.Parse(layout, "2024-11-21") // Thu
	result, timeframe, err := ParseDate(date, "Thu")
	expected := "2024-11-21"
	expectedTimeframe := goal.Day

	resultStr := result.Format(layout)
	if err == nil && (timeframe != expectedTimeframe || resultStr != expected) {
		t.Errorf("res %s; want %s", resultStr, expected)
		t.Errorf("res %s; want %s", timeframe, expectedTimeframe)
	}
}

func TestParseDate_ThuToMon(t *testing.T) {
	layout := "2006-01-02"
	date, _ := time.Parse(layout, "2024-11-21") // Thu
	result, timeframe, err := ParseDate(date, "Monday")
	expected := "2024-11-25"
	expectedTimeframe := goal.Day

	resultStr := result.Format(layout)
	if err == nil && (timeframe != expectedTimeframe || resultStr != expected) {
		t.Errorf("res %s; want %s", resultStr, expected)
		t.Errorf("res %s; want %s", timeframe, expectedTimeframe)
	}
}

func TestParseDate_27(t *testing.T) {
	layout := "2006-01-02"
	date, _ := time.Parse(layout, "2024-11-21") // Thu
	result, timeframe, err := ParseDate(date, "27")
	expected := "2024-11-27"
	expectedTimeframe := goal.Day

	resultStr := result.Format(layout)
	if err == nil && (timeframe != expectedTimeframe || resultStr != expected) {
		t.Errorf("res %s; want %s", resultStr, expected)
		t.Errorf("res %s; want %s", timeframe, expectedTimeframe)
	}
}

func TestParseDate_35(t *testing.T) {
	layout := "2006-01-02"
	date, _ := time.Parse(layout, "2024-11-21") // Thu
	_, _, err := ParseDate(date, "35")

	if err == nil {
		t.Errorf("Should have thrown an error")
	}
}

func TestIsOverdue_NilInputs(t *testing.T) {
	// Test with nil date
	var nilDate *time.Time
	timeframe := goal.Day
	if IsOverdue(nilDate, &timeframe) {
		t.Errorf("IsOverdue with nil date should return false")
	}

	// Test with nil timeframe
	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	var nilTimeframe *goal.Timeframe
	if IsOverdue(&date, nilTimeframe) {
		t.Errorf("IsOverdue with nil timeframe should return false")
	}

	// Test with both nil
	if IsOverdue(nilDate, nilTimeframe) {
		t.Errorf("IsOverdue with both nil should return false")
	}
}

func TestIsOverdue_Day(t *testing.T) {
	timeframe := goal.Day

	// Test with a date that's definitely in the past (overdue)
	pastDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	if !IsOverdue(&pastDate, &timeframe) {
		t.Errorf("Past date should be overdue")
	}

	// Test with a date that's definitely in the future (not overdue)
	futureDate := time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)
	if IsOverdue(&futureDate, &timeframe) {
		t.Errorf("Future date should not be overdue")
	}
}

func TestIsOverdue_Week(t *testing.T) {
	timeframe := goal.Week

	// Test with a definitely past week (overdue)
	pastWeekDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	if !IsOverdue(&pastWeekDate, &timeframe) {
		t.Errorf("Past week should be overdue")
	}

	// Test with a definitely future week (not overdue)
	futureWeekDate := time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)
	if IsOverdue(&futureWeekDate, &timeframe) {
		t.Errorf("Future week should not be overdue")
	}
}

func TestIsOverdue_Month(t *testing.T) {
	timeframe := goal.Month

	// Test past month (overdue)
	pastMonthDate := time.Date(2000, 6, 15, 0, 0, 0, 0, time.UTC)
	if !IsOverdue(&pastMonthDate, &timeframe) {
		t.Errorf("Past month should be overdue")
	}

	// Test future month (not overdue)
	futureMonthDate := time.Date(2100, 1, 15, 0, 0, 0, 0, time.UTC)
	if IsOverdue(&futureMonthDate, &timeframe) {
		t.Errorf("Future month should not be overdue")
	}
}

func TestIsOverdue_Quarter(t *testing.T) {
	timeframe := goal.Quarter

	// Test past quarter (overdue)
	pastQuarter := time.Date(2000, 5, 15, 0, 0, 0, 0, time.UTC)
	if !IsOverdue(&pastQuarter, &timeframe) {
		t.Errorf("Past quarter should be overdue")
	}

	// Test future quarter (not overdue)
	futureQuarter := time.Date(2100, 1, 15, 0, 0, 0, 0, time.UTC)
	if IsOverdue(&futureQuarter, &timeframe) {
		t.Errorf("Future quarter should not be overdue")
	}
}

func TestIsOverdue_Year(t *testing.T) {
	timeframe := goal.Year

	// Test past year (overdue)
	pastYear := time.Date(2000, 6, 15, 0, 0, 0, 0, time.UTC)
	if !IsOverdue(&pastYear, &timeframe) {
		t.Errorf("Past year should be overdue")
	}

	// Test future year (not overdue)
	futureYear := time.Date(2100, 6, 15, 0, 0, 0, 0, time.UTC)
	if IsOverdue(&futureYear, &timeframe) {
		t.Errorf("Future year should not be overdue")
	}
}

func TestIsOverdue_Life(t *testing.T) {
	timeframe := goal.Life

	// Life goals are never overdue, regardless of date
	pastDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	if IsOverdue(&pastDate, &timeframe) {
		t.Errorf("Life goals should never be overdue")
	}

	futureDate := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	if IsOverdue(&futureDate, &timeframe) {
		t.Errorf("Life goals should never be overdue, even with future dates")
	}
}
