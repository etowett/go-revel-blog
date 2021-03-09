package models

import (
	"context"
	"fmt"
	"go-revel-blog/app/db"
	"go-revel-blog/app/helpers"
	"strings"
)

const (
	createCommentSQL     = `insert into comments (user_id, post_id, content, created_at) values ($1, $2, $3, $4) returning id`
	getCommentSQL        = `select id, user_id, post_id, content, created_at, updated_at from comments`
	getCommentByID       = getCommentSQL + ` where id=$1`
	updateCommentSQL     = `update comments set (content, updated_at) = ($1, $2) where id = $3`
	countPostCommentsSQL = `select count(id) from comments`
	deleteCommentSQL     = `delete from comments where id=$1`
)

type (
	Comment struct {
		SequentialIdentifier
		UserID  int64  `json:"user_id"`
		PostID  int64  `json:"post_id"`
		Content string `json:"content"`
		Timestamps
	}
)

func (c *Comment) ForPost(
	ctx context.Context,
	db db.SQLOperations,
	postID int64,
	filter *Filter,
) ([]*Comment, error) {
	comments := make([]*Comment, 0)

	query, args := c.buildQuery(
		getCommentSQL,
		postID,
		filter,
	)

	rows, err := db.QueryContext(
		ctx,
		query,
		args...,
	)
	defer rows.Close()
	if err != nil {
		return comments, err
	}

	for rows.Next() {
		var comment Comment
		err = rows.Scan(
			&comment.ID,
			&comment.UserID,
			&comment.PostID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return comments, err
		}
		comments = append(comments, &comment)
	}

	return comments, err
}

func (c *Comment) CountForPost(
	ctx context.Context,
	db db.SQLOperations,
	postID int64,
	filter *Filter,
) (int, error) {
	query, args := c.buildQuery(
		countPostCommentsSQL,
		postID,
		&Filter{
			Term: filter.Term,
		},
	)
	var recordsCount int
	err := db.QueryRowContext(ctx, query, args...).Scan(&recordsCount)
	return recordsCount, err
}

func (c *Comment) Delete(
	ctx context.Context,
	db db.SQLOperations,
	id int64,
) (int64, error) {
	res, err := db.ExecContext(ctx, deleteCommentSQL, id)
	if err != nil {
		return 0, err
	}

	rowsDeleted, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsDeleted, nil
}

func (c *Comment) buildQuery(
	query string,
	postID int64,
	filter *Filter,
) (string, []interface{}) {
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	placeholder := helpers.NewPlaceholder()

	conditions = append(conditions, fmt.Sprintf(" post_id=$%d", placeholder.Touch()))
	args = append(args, postID)

	if len(conditions) > 0 {
		query += " where" + strings.Join(conditions, " and")
	}

	if filter.Per > 0 && filter.Page > 0 {
		query += fmt.Sprintf(" order by id desc limit $%d offset $%d", placeholder.Touch(), placeholder.Touch())
		args = append(args, filter.Per, (filter.Page-1)*filter.Per)
	}

	return query, args
}

func (c *Comment) Save(
	ctx context.Context,
	db db.SQLOperations,
) error {
	c.Timestamps.Touch()

	var err error
	if c.IsNew() {
		err := db.QueryRowContext(
			ctx,
			createCommentSQL,
			c.UserID,
			c.PostID,
			c.Content,
			c.Timestamps.CreatedAt,
		).Scan(&c.ID)
		return err
	}
	_, err = db.ExecContext(
		ctx,
		updateCommentSQL,
		c.UserID,
		c.PostID,
		c.Content,
		c.Timestamps.UpdatedAt,
		c.ID,
	)
	return err
}
