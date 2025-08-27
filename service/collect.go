package service

import (
	"context"
	"slices"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/fmotalleb/go-tools/log"
	"github.com/robfig/cron/v3"

	"github.com/fmotalleb/north_outage/collector"
	"github.com/fmotalleb/north_outage/config"
	"github.com/fmotalleb/north_outage/database"
	"github.com/fmotalleb/north_outage/models"
)

func startCollectService(ctx context.Context, cfg *config.Config) error {
	l := log.FromContext(ctx).Named("Scheduler")
	scheduler := cron.New(cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)))
	j, err := scheduler.AddFunc(cfg.CollectCycle, collectSilent(ctx, cfg))
	if err != nil {
		l.Error("failed to register job", zap.Error(err))
		return err
	}
	l.Info(
		"collector job registered",
		zap.Int("id", int(j)),
	)
	go scheduler.Start()
	<-ctx.Done()
	if innerCtx := scheduler.Stop(); innerCtx != nil {
		<-innerCtx.Done()
	}
	return nil
}

func collectSilent(ctx context.Context, cfg *config.Config) func() {
	l := log.FromContext(ctx).Named("CollectCycle")
	return func() {
		l.Debug("collect cycle began")
		if err := collectCycle(ctx, cfg); err != nil {
			l.Error("unhandled exception in collector cycle", zap.Error(err))
		}
	}
}

func collectCycle(ctx context.Context, cfg *config.Config) error {
	l := log.FromContext(ctx).Named("EventsGC")
	var data []models.Event
	var err error
	ctx, cancel := context.WithTimeout(ctx, cfg.CollectTimeout)
	defer cancel()
	db := database.Get()
	if err = eventsGC(cfg.RotateAfter, db); err != nil {
		l.Warn("event garbage collector failed (neglected)", zap.Error(err))
	}
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
				tx = tx.Create(&ev)
				if err = tx.Error; err != nil {
					return err
				}
			}
			return nil
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func eventsGC(maxAge time.Duration, db *gorm.DB) error {
	events := db.Table("events")
	before := time.Now().Truncate(maxAge)
	err := events.Transaction(
		func(tx *gorm.DB) error {
			res := tx.Where("end <= ?", before).Delete(true)
			return res.Error
		},
	)
	if err != nil {
		return err
	}
	return nil
}
