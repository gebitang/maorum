package models

import (
	"context"
	"gebitang.com/maorum/models/dbm"
)

type dailyItemStore struct{}

func (p *dailyItemStore) Create(ctx context.Context, item *dbm.DailyItem) error {
	if err := GetDB(ctx).Create(item).Error; err != nil {
		return err
	}
	return nil
}

func (p *dailyItemStore) FindTodayItem(ctx context.Context, timestamp int) ([]dbm.DailyItem, error) {
	var res []dbm.DailyItem
	err := GetDB(ctx).Model(&dbm.DailyItem{}).Where(" created_at > ? ", timestamp).Find(&res).Error
	return res, err
}
