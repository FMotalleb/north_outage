package service

import (
	"context"
	"slices"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/go-co-op/gocron/v2"

	"github.com/fmotalleb/go-tools/log"

	"github.com/fmotalleb/north_outage/collector"
	"github.com/fmotalleb/north_outage/config"
	"github.com/fmotalleb/north_outage/database"
	"github.com/fmotalleb/north_outage/models"
)

func startCollectService(ctx context.Context, cfg *config.Config) error {
	l := log.FromContext(ctx).Named("Scheduler")
	s, err := gocron.NewScheduler()
	if err != nil {
		l.Error("failed to build scheduler", zap.Error(err))
		return err
	}

	j, err := s.NewJob(
		gocron.CronJob(
			cfg.CollectCycle,
			false,
		),
		gocron.NewTask(
			collectSilent,
			ctx,
			cfg,
		),
	)
	if err != nil {
		l.Error("failed to register job", zap.Error(err))
		return err
	}
	var nextRun time.Time
	if nextRun, err = j.NextRun(); err != nil {
		l.Warn("failed to calculate next run time", zap.Error(err))
	}
	l.Info(
		"collector job registered",
		zap.String("id", j.ID().String()),
		zap.String("name", j.Name()),
		zap.Time("next-run", nextRun),
	)

	s.Start()
	<-ctx.Done()
	if err = s.Shutdown(); err != nil {
		l.Error("failed to shutdown scheduler", zap.Error(err))
		return err
	}
	return nil
}

func collectSilent(ctx context.Context, cfg *config.Config) func() {
	l := log.FromContext(ctx).Named("EventsGC")
	return func() {
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
