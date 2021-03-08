package controllers

import (
	"go-revel-blog/app/db"
	"go-revel-blog/app/models"
	"time"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Health() revel.Result {
	return c.RenderJSON(map[string]interface{}{
		"success":     true,
		"message":     "Ok",
		"server_time": time.Now().String(),
	})
}

func (c App) getUser(username string) *models.User {
	user := &models.User{}
	_, err := c.Session.GetInto("user", user, false)
	if user.Username == username {
		return user
	}

	newUser := &models.User{}
	foundUser, err := newUser.GetByUsername(c.Request.Context(), db.DB(), username)
	if err != nil {
		return nil
	}

	c.Session["user"] = foundUser
	return foundUser
}

func (c App) connected() *models.User {
	if c.ViewArgs["user"] != nil {
		return c.ViewArgs["user"].(*models.User)
	}
	if username, ok := c.Session["username"]; ok {
		return c.getUser(username.(string))
	}
	return nil
}

func (c App) AddUser() revel.Result {
	if user := c.connected(); user != nil {
		c.ViewArgs["user"] = user
	}
	return nil
}
