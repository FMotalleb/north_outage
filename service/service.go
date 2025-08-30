package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/fmotalleb/go-tools/log"
	"go.uber.org/zap"

	"github.com/fmotalleb/north_outage/config"
	"github.com/fmotalleb/north_outage/database"
	"github.com/fmotalleb/north_outage/models"
	"github.com/fmotalleb/north_outage/telegram"
	"github.com/fmotalleb/north_outage/web"
)

func Serve(ctx context.Context) error {
	l := log.FromContext(ctx).Named("Serve")
	cfg, err := config.Get(ctx)
	if err != nil {
		return err
	}
	ctx = config.Attach(ctx, cfg)
	db, err := database.Connect(cfg.DatabaseConnection)
	if err != nil {
		return err
	}
	if err = db.AutoMigrate(&models.Listener{}, &models.Event{}); err != nil {
		return err
	}
	l.Info("config initialized", zap.Any("cfg", cfg))
	wg := new(sync.WaitGroup)
	wg.Go(
		func() {
			err := web.Start(ctx, cfg)
			if err != nil {
				l.Error("api server collapsed", zap.Error(err))
				panic(fmt.Errorf("api server unrecoverable exception: %w", err))
			}
		},
	)
	wg.Go(
		func() {
			err := startCollector(ctx, cfg)
			if err != nil {
				l.Error("scheduler service collapsed", zap.Error(err))
				panic(fmt.Errorf("scheduler service unrecoverable exception: %w", err))
			}
		},
	)
	wg.Go(
		func() {
			err := telegram.Run(ctx, cfg)
			if err != nil {
				l.Error("telegram service collapsed", zap.Error(err))
				panic(fmt.Errorf("telegram service unrecoverable exception: %w", err))
			}
		},
	)
	wg.Wait()
	return nil
}
