package goal

import (
	"time"
)

type Timeframe string

const (
	Day     Timeframe = "day"
	Week    Timeframe = "week"
	Month   Timeframe = "month"
	Quarter Timeframe = "quarter"
	Year    Timeframe = "year"
	Life    Timeframe = "life"
)

func (t Timeframe) String() string {
	switch t {
	case Day:
		return "Day"
	case Week:
		return "Week"
	case Month:
		return "Month"
	case Quarter:
		return "Quarter"
	case Year:
		return "Year"
	case Life:
		return "Life"
	}

	return ""
}

type Goal struct {
	ID          string     `json:"id"`
	ParentId    *string    `json:"parent_id"`
	ParentTitle *string    `json:"parent_title"`
	CreatedAt   time.Time  `json:"createdAt" validate:"datetime=2006-01-02T15:04:05.999999"`
	UpdatedAt   time.Time  `json:"updatedAt" validate:"datetime=2006-01-02T15:04:05.999999"`
	Title       string     `json:"title"`
	IsDone      bool       `json:"isDone"`
	Timeframe   *Timeframe `json:"timeframe"`
	Date        *time.Time `json:"date"`
	IsArchived  bool       `json:"isArchived"`
}
