package service

import (
	"context"
	"slices"

	"github.com/fmotalleb/go-tools/log"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/fmotalleb/north_outage/collector"
	"github.com/fmotalleb/north_outage/config"
	"github.com/fmotalleb/north_outage/database"
	"github.com/fmotalleb/north_outage/models"
)

func Serve(ctx context.Context) error {
	l := log.FromContext(ctx).Named("Serve")
	cfg, err := config.Get(ctx)
	if err != nil {
		return err
	}
	ctx = config.Attach(ctx, cfg)
	db, err := database.NewDB(cfg.DatabaseConnection)
	if err != nil {
		return err
	}
	if err = db.AutoMigrate(&models.Listener{}, &models.Event{}); err != nil {
		return err
	}
	l.Info("config initialized", zap.Any("cfg", cfg))
	collectCycle(ctx, db)
	return nil
}

func collectCycle(ctx context.Context, db *gorm.DB) error {
	var data []models.Event
	var err error
	if data, err = collector.Collect(ctx); err != nil {
		return err
	}

	events := db.Table("events")
	var oldHash []string
	events.Select("hash").Find(&oldHash)
	err = db.Transaction(
		func(tx *gorm.DB) error {
			for _, ev := range data {
				if slices.Contains(oldHash, ev.Hash) {
					continue
				}
				tx.Create(&ev)
			}
			return nil
		},
	)
	if err != nil {
		return err
	}
	return nil
}
