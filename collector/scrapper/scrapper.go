package scrapper

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/fmotalleb/go-tools/decoder"
	"github.com/fmotalleb/go-tools/log"
	"github.com/fmotalleb/scrapper-go/engine"
	"go.uber.org/zap"

	"github.com/fmotalleb/north_outage/config"
	"github.com/fmotalleb/north_outage/models"
)

// Run executes the configured collector engine and returns a deduplicated list of events.
func Run(ctx context.Context) ([]models.Event, error) {
	logger := log.FromContext(ctx).Named("collector")

	cfg, err := config.Get(ctx)
	if err != nil {
		return nil, err
	}

	rawResult, err := engine.ExecuteConfig(ctx, cfg.CollectorConfig)
	if err != nil {
		logger.Error("collector execution failed", zap.Error(err))
		return nil, err
	}

	logger.Info("collector finished successfully")

	events, reshapeErr := transformResult(rawResult)
	if reshapeErr != nil {
		// Preserve detailed error reporting but keep the main flow working.
		if multiErr, ok := reshapeErr.(interface{ Unwrap() []error }); ok {
			logger.Error("transform produced some errors (ignored)", zap.Errors("errors", multiErr.Unwrap()))
		} else {
			logger.Error("transform produced some errors (ignored)", zap.Error(reshapeErr))
		}
	}

	return events, nil
}

// transformResult converts engine output into []models.Event, deduplicating by Hash.
func transformResult(data map[string]any) ([]models.Event, error) {
	seen := make(map[string]struct{})
	events := make([]models.Event, 0)
	var errs []error

	// helper to decode a reference key into a []map[string]string
	decodeRef := func(ref string) []map[string]string {
		val, ok := data[ref]
		if !ok {
			return nil
		}
		list, ok := val.([]map[string]any)
		if !ok {
			return nil
		}

		dst := make([]map[string]string, len(list))
		if err := decoder.Decode(&dst, list); err != nil {
			errs = append(errs, fmt.Errorf("failed to decode %q: %w", ref, err))
		}
		return dst
	}

	for key, value := range data {
		if !strings.HasPrefix(key, "map.") {
			continue
		}

		mapName := strings.TrimPrefix(key, "map.")
		refString, ok := value.(string)
		if !ok {
			errs = append(errs, fmt.Errorf("invalid mapping type for %q: expected string, got %T", key, value))
			continue
		}

		for _, ref := range strings.Split(refString, ",") {
			for _, row := range decodeRef(ref) {
				ev, ok := normalize(mapName, row)
				if !ok {
					errs = append(errs, fmt.Errorf("failed to normalize entry for %q: %v", mapName, row))
					continue
				}
				if _, exists := seen[ev.Hash]; exists {
					continue
				}
				seen[ev.Hash] = struct{}{}
				events = append(events, *ev)
			}
		}
	}

	if len(errs) > 0 {
		return events, errors.Join(errs...)
	}
	return events, nil
}
