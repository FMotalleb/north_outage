package service

import (
	"context"

	"github.com/fmotalleb/go-tools/defaulter"
	"github.com/fmotalleb/go-tools/log"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/fmotalleb/north_outage/config"
	"github.com/fmotalleb/north_outage/database"
	"github.com/fmotalleb/north_outage/models"
)

func Serve(ctx context.Context) error {
	l := log.FromContext(ctx).Named("Serve")
	cfg, err := readConfig()
	if err != nil {
		return err
	}
	db, err := database.NewDB(cfg.DatabaseConnection)
	if err != nil {
		return err
	}
	if err = db.AutoMigrate(&models.Listener{}, &models.Event{}); err != nil {
		return err
	}
	l.Info("config initialized", zap.Any("cfg", cfg))
	return nil
}

func readConfig() (*config.Config, error) {
	cfg := &config.Config{}
	defaulter.ApplyDefaults(cfg, nil)
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
