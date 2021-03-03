package controllers

import (
	"go-revel-blog/app/db"
	"go-revel-blog/app/forms"
	"go-revel-blog/app/models"
	"go-revel-blog/app/routes"
	"go-revel-blog/app/webutils"
	"time"

	"github.com/revel/revel"
	null "gopkg.in/guregu/null.v4"
)

type Categories struct {
	App
}

func (c Categories) All() revel.Result {
	paginationFilter, err := webutils.FilterFromQuery(c.Params)
	if err != nil {
		c.Log.Errorf("could not filter from params: %v", err)
		return c.Render(response("error", "could not get Categories", "failed"))
	}

	newCategory := models.Category{}
	data, err := newCategory.All(c.Request.Context(), db.DB(), paginationFilter)
	if err != nil {
		c.Log.Errorf("could not get newCategory: %v", err)
		return c.Render(response(data, "could not get newCategory", "failed"))
	}

	recordsCount, err := newCategory.Count(c.Request.Context(), db.DB(), paginationFilter)
	if err != nil {
		c.Log.Errorf("could not count newCategory: %v", err)
		return c.Render(response(data, "could not count newCategory", "failed"))
	}

	ret := map[string]interface{}{
		"Categories": data,
		"Pagination": models.NewPagination(recordsCount, paginationFilter.Page, paginationFilter.Per),
	}
	result := response(ret, "get Categories successful", "success")
	return c.Render(result)
}

func (c Categories) Get(id int64) revel.Result {
	newCat := &models.Category{}
	data, err := newCat.GetByID(c.Request.Context(), db.DB(), id)
	c.Log.Infof("id %v, data %v, err %v", id, data, err)
	if err != nil {
		if err.Error() == "record not found" {
			result := response(data, "data not found", "fail")
			return c.Render(result)
		}
		result := response(err.Error(), "error get category", "failed")
		return c.Render(result)
	}

	result := response(data, "get data successfull", "success")
	return c.Render(result)
}

func (c Categories) Create() revel.Result {
	return c.Render()
}

func (c Categories) Save(category *forms.Category) revel.Result {

	v := c.Validation
	category.Validate(v)

	if v.HasErrors() {
		v.Keep()
		c.FlashParams()
		return c.Redirect(routes.Categories.Create())
	}

	newCat := &models.Category{
		UserID:      category.UserID,
		Name:        category.Name,
		Description: category.Description,
	}
	err := newCat.Save(c.Request.Context(), db.DB())
	if err != nil {
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("Could not save category: %v", err)
		return c.Redirect(routes.Categories.Create())
	}

	c.Flash.Success("Category created - " + newCat.Name)
	return c.Redirect(routes.Categories.Get(newCat.ID))
}

func (c Categories) Edit(id int64) revel.Result {
	newCategory := &models.Category{}
	category, err := newCategory.GetByID(c.Request.Context(), db.DB(), id)
	if err != nil {
		return c.NotFound("Category does not exist")
	}
	return c.Render(category)
}

func (c Categories) Update(id int64, category *forms.Category) revel.Result {
	v := c.Validation
	category.Validate(v)

	if v.HasErrors() {
		v.Keep()
		c.FlashParams()
		return c.Redirect(routes.Categories.Edit(id))
	}

	newCategory := &models.Category{}
	existingCategory, err := newCategory.GetByID(c.Request.Context(), db.DB(), id)
	if err != nil {
		c.Flash.Error("could not get existing cat:", err)
		c.FlashParams()
		return c.Redirect(routes.Categories.Edit(id))
	}

	if existingCategory == nil {
		return c.NotFound("Category does not exist")
	}

	existingCategory.Name = category.Name
	existingCategory.Description = category.Description
	existingCategory.UpdatedAt = null.TimeFrom(time.Now())

	err = existingCategory.Save(c.Request.Context(), db.DB())
	if err != nil {
		v.Keep()
		c.FlashParams()
		return c.Redirect(routes.Categories.Edit(id))
	}

	return c.Redirect(routes.Categories.All())
}

func (c Categories) Delete(id int64) revel.Result {
	newCategory := &models.Category{}
	_, err := newCategory.Delete(c.Request.Context(), db.DB(), id)
	if err != nil {
		c.Log.Errorf("error newCategory delete: %v", err)
	}

	return c.Redirect(routes.Categories.All())
}
