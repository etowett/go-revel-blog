package controllers

import (
	"go-revel-blog/app/db"
	"go-revel-blog/app/forms"
	"go-revel-blog/app/models"
	"go-revel-blog/app/routes"
	"net/http"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	App
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

func (c *Users) Get(id int64) revel.Result {
	newUser := models.User{}
	user, err := newUser.GetByID(c.Request.Context(), db.DB(), id)
	if err != nil {
		if err.Error() == "record not found" {
			result := response(err.Error(), "data not found", "fail")
			c.Response.SetStatus(http.StatusNotFound)
			return c.Render(result)
		}
		result := response(err.Error(), "error get user", "failed")
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.Render(result)
	}

	c.Response.SetStatus(http.StatusFound)
	result := response(user, "get user successfull", "success")
	return c.Render(result)
}
