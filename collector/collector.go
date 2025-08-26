package collector

import (
	"context"
	"strings"
	"time"

	"github.com/fmotalleb/go-tools/decoder"
	"github.com/fmotalleb/go-tools/log"
	"github.com/fmotalleb/scrapper-go/engine"
	"github.com/mshafiee/jalali"
	"github.com/spf13/cast"

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
	// c, _ := os.ReadFile("test.json")
	// var result map[string]any
	// json.Unmarshal(c, &result)
	l.Info("scrape finished")
	events := reshape(ctx, result)
	return events, nil
}

func reshape(ctx context.Context, input map[string]any) []models.Event {
	repeats := make(map[string]int)
	result := make([]models.Event, 0)
	collect := func(k string) []map[string]string {
		v := input[k]
		switch v := v.(type) {
		case []map[string]any:
			dst := make([]map[string]string, len(v))
			decoder.Decode(&dst, v)
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
			continue
		}

		for _, ref := range strings.Split(refs, ",") {
			data := collect(ref)
			for _, v := range data {
				ev, ok := normalize(k, v)
				if !ok {
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
	return result
}

const (
	keyAddr  = "آدرس"
	keyFrom  = "از ساعت"
	keyUntil = "تا ساعت"
	keyDate  = "تاریخ"
	// Who cares really
	// keyCause = "نوع خاموشی"
)

func normalize(city string, input map[string]string) (*models.Event, bool) {
	obj := new(models.Event)
	addr, ok := input[keyAddr]
	if !ok {
		return nil, false
	}
	fromStr, ok := input[keyFrom]
	if !ok {
		return nil, false
	}
	untilStr, ok := input[keyUntil]
	if !ok {
		return nil, false
	}
	dateStr, ok := input[keyDate]
	if !ok {
		return nil, false
	}
	date, ok := parseJalali(dateStr)
	if !ok {
		return nil, false
	}

	obj.Address = persianFixer(addr)
	obj.City = persianFixer(city)
	if v, ok := parseTime(fromStr); ok {
		obj.Start = date.ToTime().Add(v)
	} else {
		return nil, false
	}
	if v, ok := parseTime(untilStr); ok {
		obj.End = date.ToTime().Add(v)
	} else {
		obj.End = obj.Start.Add(2 * time.Hour)
	}
	obj.ResetHash()
	return obj, true
}

func parseJalali(input string) (jalali.JalaliTime, bool) {
	const size = 3
	seg := strings.Split(input, "/")
	if len(seg) != size {
		return jalali.Date(0, 0, 0, 0, 0, 0, 0, time.Local), false
	}
	y := cast.ToInt(seg[0])
	m := cast.ToInt(seg[1])
	d := cast.ToInt(seg[2])
	mn := jalali.Month(m)
	return jalali.Date(y, mn, d, 0, 0, 0, 0, time.Local), true
}

func parseTime(input string) (time.Duration, bool) {
	const size = 2
	seg := strings.Split(input, ":")
	if len(seg) != size {
		return time.Second * 0, false
	}
	h := cast.ToInt(seg[0])
	m := cast.ToInt(seg[1])
	t := time.Hour*time.Duration(h) + time.Minute*time.Duration(m)
	return t, true
}

var replacements = map[rune]rune{
	// Kaf
	'ك': 'ک',

	// Yeh & Alef Maksura
	'ي': 'ی',
	'ى': 'ی',
	'ئ': 'ی',

	// Heh variants
	'ة': 'ه',
	'ۀ': 'ه',
	'ہ': 'ه',
	'ھ': 'ه',

	// Waw variants
	'ؤ': 'و',

	// Digits (Arabic → Persian)
	'٠': '0',
	'١': '1',
	'٢': '2',
	'٣': '3',
	'٤': '4',
	'٥': '5',
	'٦': '6',
	'٧': '7',
	'٨': '8',
	'٩': '9',
}

func persianFixer(input string) string {
	var out strings.Builder
	for _, r := range input {
		// Remove Arabic diacritics
		if r >= 0x064B && r <= 0x0652 {
			continue
		}
		if rep, ok := replacements[r]; ok {
			out.WriteRune(rep)
		} else {
			out.WriteRune(r)
		}
	}
	return out.String()
}
