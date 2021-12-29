package dbm

import (
	"context"
	"time"
)

type (
	DailyItem struct {
		ID        int `json:"id"`
		ItemType  int
		ItemComm  string
		CreatedAt int `json:"createdAt" gorm:"autoUpdateTime:milli"`
		Min       int
		Comment   string
		UpdatedAt time.Time `json:"updatedAt"`
	}

	DailyItemStore interface {
		Create(ctx context.Context, di *DailyItem) error
		FindTodayItem(ctx context.Context, timestamp int) ([]DailyItem, error)
	}
)
