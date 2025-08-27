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

func Collect(ctx context.Context) ([]models.Event, error) {
	l := log.FromContext(ctx).Named("scraper")
	cfg, err := config.Get(ctx)
	if err != nil {
		return nil, err
	}
	result, err := engine.ExecuteConfig(ctx, cfg.CollectorConfig)
	if err != nil {
		l.Error("scrapper fatal error", zap.Error(err))
		return nil, err
	}
	l.Info("scrape finished")
	events, err := reshape(result)
	if err != nil {
		errs, ok := err.(interface {
			Unwrap() []error
		})
		var field zap.Field
		if ok {
			field = zap.Errors("errors", errs.Unwrap())
		} else {
			field = zap.Error(err)
		}
		l.Error("reshape faced some errors (neglected)", field)
	}
	return events, nil
}

func reshape(input map[string]any) ([]models.Event, error) {
	repeats := make(map[string]int)
	result := make([]models.Event, 0)
	allErrs := make([]error, 0)
	collect := func(k string) []map[string]string {
		v := input[k]
		switch v := v.(type) {
		case []map[string]any:
			dst := make([]map[string]string, len(v))
			if err := decoder.Decode(&dst, v); err != nil {
				allErrs = append(allErrs, err)
			}
			return dst
		default:
			return []map[string]string{}
		}
	}

	for k, v := range input {
		if !strings.HasPrefix(k, "map.") {
			continue
		}
		k = strings.TrimPrefix(k, "map.")
		refs, ok := v.(string)
		if !ok {
			er := fmt.Errorf("unexpected mapping type, required `string` got `%T`", v)
			allErrs = append(allErrs, er)
			continue
		}

		for _, ref := range strings.Split(refs, ",") {
			data := collect(ref)
			for _, v := range data {
				ev, ok := normalize(k, v)
				if !ok {
					er := fmt.Errorf("failed to normalize input, %v", v)
					allErrs = append(allErrs, er)
					continue
				}
				if repeats[ev.Hash] > 0 {
					continue
				}
				repeats[ev.Hash]++
				result = append(result, *ev)
			}
		}
	}
	var err error
	if len(allErrs) != 0 {
		err = errors.Join(allErrs...)
	}
	return result, err
}
