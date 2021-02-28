package controllers

import (
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Health() revel.Result {
	return c.RenderJSON(map[string]interface{}{
		"success": true,
		"status":  "Ok",
	})
}
