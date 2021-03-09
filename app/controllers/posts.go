package controllers

import (
	"go-revel-blog/app/db"
	"go-revel-blog/app/entities"
	"go-revel-blog/app/forms"
	"go-revel-blog/app/models"
	"go-revel-blog/app/routes"
	"go-revel-blog/app/webutils"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/revel/revel"
	null "gopkg.in/guregu/null.v4"
)

type Posts struct {
	App
}

func (c *Posts) MarkdownHTML(text string) string {
	// md := []byte(text)
	// return string(markdown.ToHTML(md, nil, nil))
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)

	md := []byte(text)
	html := markdown.ToHTML(md, parser, nil)
	return string(html)
}

func (c *Posts) All() revel.Result {
	var result entities.Response
	paginationFilter, err := webutils.FilterFromQuery(c.Params)
	if err != nil {
		c.Log.Errorf("could not filter from params: %v", err)
		result = entities.Response{
			Success: false,
			Message: "Failed to parse page filters",
		}
		return c.Render(result)
	}

	newPost := models.Post{}
	data, err := newPost.All(c.Request.Context(), db.DB(), paginationFilter)
	if err != nil {
		c.Log.Errorf("could not get posts with id: %v", err)
		result = entities.Response{
			Success: false,
			Message: "Could not get posts",
		}
		return c.Render(result)
	}

	recordsCount, err := newPost.Count(c.Request.Context(), db.DB(), paginationFilter)
	if err != nil {
		c.Log.Errorf("could not get posts count: %v", err)
		result = entities.Response{
			Success: false,
			Message: "Could not get posts count",
		}
		return c.Render(result)
	}

	result = entities.Response{
		Success: true,
		Data: map[string]interface{}{
			"Posts":      data,
			"Pagination": models.NewPagination(recordsCount, paginationFilter.Page, paginationFilter.Per),
		},
	}
	return c.Render(result)
}

func (c *Posts) Get(id int64) revel.Result {
	var result entities.Response
	newPost := &models.Post{}
	post, err := newPost.ByID(c.Request.Context(), db.DB(), id)
	if err != nil {
		c.Log.Errorf("could not get post: %v", err)
		result = entities.Response{
			Success: false,
			Message: "Could not get the post",
			Data:    map[string]interface{}{"Post": post},
		}
		return c.Render(result)
	}

	post.Content = c.MarkdownHTML(post.Content)

	paginationFilter, err := webutils.FilterFromQuery(c.Params)
	if err != nil {
		c.Log.Errorf("could not filter from params: %v", err)
		result = entities.Response{
			Success: false,
			Message: "Failed to parse page filters",
		}
		return c.Render(result)
	}

	newComment := models.Comment{}
	comments, err := newComment.ForPost(c.Request.Context(), db.DB(), id, paginationFilter)
	if err != nil {
		c.Log.Errorf("could not get comments for post id %v: %v", id, err)
		result = entities.Response{
			Success: false,
			Message: "Could not get comments",
		}
		return c.Render(result)
	}

	recordsCount, err := newComment.CountForPost(c.Request.Context(), db.DB(), id, paginationFilter)
	if err != nil {
		c.Log.Errorf("could not get comment count for post %v: %v", id, err)
		result = entities.Response{
			Success: false,
			Message: "Could not get comments count",
		}
		return c.Render(result)
	}

	result = entities.Response{
		Success: true,
		Data: map[string]interface{}{
			"Post":       post,
			"Comments":   comments,
			"Pagination": models.NewPagination(recordsCount, paginationFilter.Page, paginationFilter.Per),
		},
	}
	return c.Render(result)
}

func (c *Posts) New() revel.Result {
	return c.Render()
}

func (c *Posts) Save(post *forms.Post) revel.Result {

	v := c.Validation
	post.Validate(v)

	if v.HasErrors() {
		v.Keep()
		c.FlashParams()
		return c.Redirect(routes.Posts.New())
	}

	newPost := &models.Post{
		UserID:  post.UserID,
		Title:   post.Title,
		Content: post.Content,
		Tag:     post.Tag,
	}
	err := newPost.Save(c.Request.Context(), db.DB())
	if err != nil {
		c.Log.Errorf("could not save post: %v", err)
		c.Flash.Error("internal server error -  could not save post")
		c.FlashParams()
		return c.Redirect(routes.Posts.New())
	}

	c.Flash.Success("post created - " + newPost.Title)
	return c.Redirect(routes.Posts.Get(newPost.ID))
}

func (c *Posts) SaveComment(id int64, comment *forms.Comment) revel.Result {
	v := c.Validation
	comment.Validate(v)

	if v.HasErrors() {
		v.Keep()
		c.FlashParams()
		return c.Redirect(routes.Posts.Get(id))
	}

	newComment := &models.Comment{
		UserID:  comment.UserID,
		PostID:  comment.PostID,
		Content: comment.Content,
	}
	err := newComment.Save(c.Request.Context(), db.DB())
	if err != nil {
		c.Log.Errorf("could not save comment for post %v: %v", id, err)
		c.Flash.Error("internal server error -  could not save comment")
		c.FlashParams()
		return c.Redirect(routes.Posts.Get(id))
	}

	c.Flash.Success("Comment added success ")
	return c.Redirect(routes.Posts.Get(id))
}

func (c *Posts) Edit(id int64) revel.Result {
	newPost := &models.Post{}
	post, err := newPost.ByID(c.Request.Context(), db.DB(), id)
	if err != nil {
		c.Log.Errorf("could not get post with id %+v: %v", id, err)
		return c.Redirect(routes.Posts.All())
	}
	return c.Render(post)
}

func (c *Posts) Update(id int64, post *forms.Post) revel.Result {
	v := c.Validation
	post.Validate(v)

	if v.HasErrors() {
		v.Keep()
		c.FlashParams()
		return c.Redirect(routes.Posts.Edit(id))
	}

	newPost := &models.Post{}
	existingPost, err := newPost.ByID(c.Request.Context(), db.DB(), id)
	if err != nil {
		c.Log.Errorf("could not get post with id %+v: %v", id, err)
		c.Flash.Error("internal server error -  could not save post")
		c.FlashParams()
		return c.Redirect(routes.Posts.Edit(id))
	}

	if existingPost == nil {
		c.Log.Errorf("could not get post with id %+v: %v", id, err)
		return c.Redirect(routes.Posts.All())
	}

	existingPost.Title = post.Title
	existingPost.Content = post.Content
	existingPost.Tag = post.Tag
	existingPost.UpdatedAt = null.TimeFrom(time.Now())

	c.Log.Infof("existingPost: %+v", existingPost)
	err = existingPost.Save(c.Request.Context(), db.DB())
	if err != nil {
		c.Log.Errorf("could not update post with id %+v: %v", id, err)
		c.Flash.Error("internal server error -  could not save post")
		c.FlashParams()
		return c.Redirect(routes.Posts.Edit(id))
	}

	return c.Redirect(routes.Posts.Get(id))
}

func (c *Posts) Delete(id int64) revel.Result {
	newPost := &models.Post{}
	_, err := newPost.Delete(c.Request.Context(), db.DB(), id)
	if err != nil {
		c.Log.Errorf("error newPost =[%v] delete: %v", id, err)
	}

	return c.Redirect(routes.Posts.All())
}

func (c *Posts) DeleteComment(postid, id int64) revel.Result {
	newComment := &models.Comment{}
	_, err := newComment.Delete(c.Request.Context(), db.DB(), id)
	if err != nil {
		c.Log.Errorf("error newComment =[%v] delete: %v", id, err)
	}

	return c.Redirect(routes.Posts.Get(postid))
}
