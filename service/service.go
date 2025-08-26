package service

import (
	"context"

	"github.com/fmotalleb/go-tools/log"
	"go.uber.org/zap"

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
	var data []models.Event
	if data, err = collector.Collect(ctx); err != nil {
		return err
	}
	tx := db.Create(data)
	err = tx.Apply(db.Config)
	println(err)
	return nil
}
