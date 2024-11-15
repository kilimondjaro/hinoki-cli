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
