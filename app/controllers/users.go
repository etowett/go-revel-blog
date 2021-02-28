package controllers

import (
	"go-revel-blog/app/db"
	"go-revel-blog/app/forms"
	"go-revel-blog/app/models"
	"go-revel-blog/app/routes"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	App
}

func (c Users) getUser(username string) *models.User {
	user := &models.User{}
	_, err := c.Session.GetInto("user", user, false)
	if user.Username == username {
		return user
	}

	newUser := models.User{}
	foundUser, err := newUser.GetByUsername(c.Request.Context(), db.DB(), username)
	if err != nil {
		return nil
	}

	c.Session["user"] = foundUser
	return foundUser
}

func (c Users) Register() revel.Result {
	return c.Render()
}

func (c Users) Save(user *forms.User) revel.Result {
	v := c.Validation
	user.Validate(v)
	if v.HasErrors() {
		v.Keep()
		c.FlashParams()
		return c.Redirect(routes.Users.Register())
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Errorf("error generate password hash: %v", err)
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("Could not generate password hash")
		return c.Redirect(routes.Users.Register())
	}
	newUser := models.User{
		Username:     user.Username,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		PasswordHash: string(passwordHash[:]),
	}

	err = newUser.Save(c.Request.Context(), db.DB())
	if err != nil {
		c.Log.Errorf("error user create: %v", err)
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("Could not save user")
		return c.Redirect(routes.Users.Register())
	}

	c.Session["username"] = newUser.Username
	c.Flash.Success("Welcome, " + newUser.FirstName)
	return c.Redirect(routes.Posts.All())
}

func (c Users) Login() revel.Result {
	return c.Render()
}

func (c Users) DoLogin(login *forms.Login) revel.Result {
	v := c.Validation
	login.Validate(v)

	if v.HasErrors() {
		v.Keep()
		c.FlashParams()
		return c.Redirect(routes.Users.Login())
	}

	user := c.getUser(login.Username)

	if user == nil {
		v.Keep()
		c.Flash.Error("Could not find user with that username")
		c.FlashParams()
		return c.Redirect(routes.Users.Login())
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(login.Password))
	if err != nil {
		v.Keep()
		c.Flash.Error("Invalid password provided")
		c.FlashParams()
		return c.Redirect(routes.Users.Login())
	}

	c.Session["username"] = login.Username
	if login.Remember {
		c.Session.SetNoExpiration()
	} else {
		c.Session.SetDefaultExpiration()
	}
	c.Flash.Success("Welcome, " + login.Username)
	return c.Redirect(routes.Posts.All())
}

func (c Users) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.Posts.All())
}
