package controllers

import (
	"go-revel-blog/app/db"
	"go-revel-blog/app/models"
	"go-revel-blog/app/webutils"

	"github.com/revel/revel"
)

type Posts struct {
	App
}

func (c Posts) All() revel.Result {
	paginationFilter, err := webutils.FilterFromQuery(c.Params)
	if err != nil {
		c.Log.Errorf("could not filter from params: %v", err)
		return c.Render(response("error", "could not get posts", "failed"))
	}

	newPost := models.Post{}
	data, err := newPost.All(c.Request.Context(), db.DB(), paginationFilter)
	if err != nil {
		c.Log.Errorf("could not get newPost: %v", err)
		return c.Render(response(data, "could not get newPost", "failed"))
	}

	recordsCount, err := newPost.Count(c.Request.Context(), db.DB(), paginationFilter)
	if err != nil {
		c.Log.Errorf("could not count newPost: %v", err)
		return c.Render(response(data, "could not count newPost", "failed"))
	}

	ret := map[string]interface{}{
		"Posts":      data,
		"Pagination": models.NewPagination(recordsCount, paginationFilter.Page, paginationFilter.Per),
	}
	result := response(ret, "get posts successful", "success")
	return c.Render(result)
}

func (c Posts) Get(id int64) revel.Result {
	return c.Render()
}
