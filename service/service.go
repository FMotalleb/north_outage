package service

import (
	"context"

	"github.com/fmotalleb/go-tools/defaulter"
	"github.com/fmotalleb/go-tools/log"
	"go.uber.org/zap"

	"github.com/fmotalleb/north_outage/config"
)

func Serve(ctx context.Context) error {
	l := log.FromContext(ctx).Named("Serve")
	cfg := &config.Config{}
	defaulter.ApplyDefaults(cfg, nil)
	l.Info("config initialized", zap.Any("cfg", cfg))
	return nil
}
