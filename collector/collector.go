package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	cfgutil "github.com/fmotalleb/go-tools/config"
	"github.com/fmotalleb/go-tools/decoder"
	"github.com/fmotalleb/go-tools/log"
	sc "github.com/fmotalleb/scrapper-go/config"
	"github.com/fmotalleb/scrapper-go/engine"

	"github.com/fmotalleb/north_outage/config"
	"github.com/fmotalleb/north_outage/models"
)

func Collect(ctx context.Context) ([]models.Event, error) {
	l := log.FromContext(ctx)
	cfg, err := config.Get(ctx)
	if err != nil {
		return nil, err
	}
	result, err := engine.ExecuteConfig(ctx, cfg.CollectorConfig)
	if err != nil {
		return nil, err
	}
	l.Info("scrape finished")
	out, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	f, _ := os.Create("test.json")
	f.Write(out)
	// till here, it extracts data with city mapping
	return nil, nil
}

func parse(ctx context.Context, dst *sc.ExecutionConfig, path string) error {
	cfg, err := cfgutil.ReadAndMergeConfig(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to read and merge configs: %w", err)
	}
	// hooks.RegisterHook(template.StringTemplateEvaluate())
	decoder, err := decoder.Build(dst)
	if err != nil {
		return fmt.Errorf("create decoder: %w", err)
	}

	if err := decoder.Decode(cfg); err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	return nil
}
