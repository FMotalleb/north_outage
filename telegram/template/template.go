package template

import (
	"strings"
	"time"

	"github.com/fmotalleb/go-tools/template"
	"github.com/go-telegram/bot/models"
	"github.com/mshafiee/jalali"
	"github.com/spf13/cast"
)

var funcs = map[string]any{
	"toJalali": toJalali,
	"jFormat":  jFormat,
}

func EvaluateTemplate(tmplt string, data map[string]any, update *models.Update) (string, error) {
	if data == nil {
		data = make(map[string]any)
	}
	data["msg"] = update.Message
	data["name"] = getName(&update.Message.Chat)
	out, err := template.EvaluateTemplateWithFuncs(tmplt, data, funcs)
	return out, err
}

func getName(c *models.Chat) string {
	sb := new(strings.Builder)
	if c.FirstName != "" {
		sb.WriteString(c.FirstName)
		if c.LastName != "" {
			sb.WriteRune(' ')
			sb.WriteString(c.LastName)
		}
	} else if c.LastName != "" {
		sb.WriteString(c.LastName)
	}
	if sb.String() == "" {
		return c.Title
	}
	return sb.String()
}

func toJalali(t any) jalali.JalaliTime {
	realValue := cast.ToTime(t)
	return jalali.JalaliFromTime(realValue)
}

func jFormat(format string, t time.Time) string {
	return toJalali(t).Format(format)
}
