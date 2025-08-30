package service

import (
	"context"
	"fmt"
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
	parser := cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)
	c, err := parser.Parse(cfg.CollectCycle)
	if err != nil {
		return fmt.Errorf("failed to parse given cron string,cron=%s err=%w", cfg.CollectCycle, err)
	}
	now := time.Now()
	next := c.Next(now)
	timeTillNext := time.Until(next)

	scheduler := cron.New(cron.WithParser(parser))

	j, err := scheduler.AddFunc(cfg.CollectCycle, collectSilent(ctx, cfg))
	if err != nil {
		l.Error("failed to register job", zap.Error(err))
		return err
	}
	l.Info(
		"collector job registered",
		zap.Int("id", int(j)),
		zap.Time("next-run", next),
		zap.Duration("time-til-next", timeTillNext),
	)
	if cfg.CollectOnStart && timeTillNext > cfg.CollectOnStartThreshold {
		l.Info(
			"collect on start threshold reached, starting to collect",
			zap.Duration("threshold", cfg.CollectOnStartThreshold),
			zap.Duration("time-til-next", timeTillNext),
		)
		go collectSilent(ctx, cfg)()
	}
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
	l := log.FromContext(ctx).Named("CollectCycle")
	var data []models.Event
	var err error
	ctx, cancel := context.WithTimeout(ctx, cfg.CollectTimeout)
	defer cancel()
	db := database.Get()
	l.Info("booting the collector")
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
	if err = eventsGC(cfg.RotateAfter, db); err != nil {
		l := log.FromContext(ctx).Named("EventsGC")
		l.Warn("event garbage collector failed (neglected)", zap.Error(err))
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
