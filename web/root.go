package web

import (
	"github.com/labstack/echo/v4"

	"github.com/fmotalleb/north_outage/web/front"
)

func init() {
	RegisterEndpoint(
		func(web *echo.Echo) {
			fs, err := front.GetDist()
			if err != nil {
				panic("filesystem setup failed")
			}
			web.StaticFS("/", fs)
		},
	)
}
