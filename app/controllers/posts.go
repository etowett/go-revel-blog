package controllers

import "github.com/revel/revel"

type Posts struct {
	App
}

func (c Posts) All() revel.Result {
	return c.Render()
}
