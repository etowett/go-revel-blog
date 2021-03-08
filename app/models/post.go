package models

import (
	"context"
	"fmt"
	"go-revel-blog/app/db"
	"go-revel-blog/app/helpers"
	"strings"

	"github.com/revel/revel"
)

const (
	createPostSQL = `insert into posts (user_id, title, content, tag, created_at) values ($1, $2, $3, $4, $5) returning id`
	getPostSQL    = `select id, user_id, title, content, tag, created_at, updated_at from posts`
	getPostByID   = getPostSQL + ` where id=$1`
	updatePostSQL = `update posts set (title, content, tag, updated_at) = ($1, $2, $3, &4) where id = $5`
	countPostSQL  = `select count(id) from posts`
	deletePostSQL = `delete from posts where id=$1`
)

type (
	Post struct {
		SequentialIdentifier
		UserID  int64  `json:"user_id"`
		Title   string `json:"title"`
		Content string `json:"content"`
		Tag     string `json:"tag"`
		Timestamps
	}
)

func (p *Post) All(
	ctx context.Context,
	db db.SQLOperations,
	filter *Filter,
) ([]*Post, error) {
	posts := make([]*Post, 0)

	query, args := p.buildQuery(
		getPostSQL,
		filter,
	)

	rows, err := db.QueryContext(
		ctx,
		query,
		args...,
	)
	defer rows.Close()
	if err != nil {
		return posts, err
	}

	for rows.Next() {
		var post Post
		err = rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.Tag,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return posts, err
		}
		posts = append(posts, &post)
	}

	return posts, err
}

func (q *Post) Count(
	ctx context.Context,
	db db.SQLOperations,
	filter *Filter,
) (int, error) {
	query, args := q.buildQuery(
		countPostSQL,
		&Filter{
			Term: filter.Term,
		},
	)
	var recordsCount int
	err := db.QueryRowContext(ctx, query, args...).Scan(&recordsCount)
	return recordsCount, err
}

func (p *Post) Delete(
	ctx context.Context,
	db db.SQLOperations,
	id int64,
) (int64, error) {
	res, err := db.ExecContext(ctx, deletePostSQL, id)
	if err != nil {
		return 0, err
	}

	rowsDeleted, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsDeleted, nil
}

func (p *Post) ByID(
	ctx context.Context,
	db db.SQLOperations,
	id int64,
) (*Post, error) {
	row := db.QueryRowContext(ctx, getPostByID, id)
	return p.scan(row)
}

func (p *Post) Save(
	ctx context.Context,
	db db.SQLOperations,
) error {
	p.Timestamps.Touch()

	var err error
	if p.IsNew() {
		err := db.QueryRowContext(
			ctx,
			createPostSQL,
			p.UserID,
			p.Title,
			p.Content,
			p.Tag,
			p.Timestamps.CreatedAt,
		).Scan(&p.ID)
		return err
	}
	revel.AppLog.Errorf("updatePostSQL: %v", updatePostSQL)
	revel.AppLog.Errorf("post: %+v", p)
	_, err = db.ExecContext(
		ctx,
		updatePostSQL,
		p.Title,
		p.Content,
		p.Tag,
		p.Timestamps.UpdatedAt,
		p.ID,
	)
	return err
}

func (*Post) scan(
	row db.RowScanner,
) (*Post, error) {
	var p Post
	err := row.Scan(
		&p.ID,
		&p.UserID,
		&p.Title,
		&p.Content,
		&p.Tag,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	return &p, err
}

func (p *Post) buildQuery(
	query string,
	filter *Filter,
) (string, []interface{}) {
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	placeholder := helpers.NewPlaceholder()

	if filter.Term != "" {
		likeStmt := make([]string, 0)
		columns := []string{"title", "content"}
		for _, col := range columns {
			search := fmt.Sprintf(" (lower(%s) like '%%' || $%d || '%%')", col, placeholder.Touch())
			likeStmt = append(likeStmt, search)
			args = append(args, filter.Term)
		}
		conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(likeStmt, " or")))
	}

	if len(conditions) > 0 {
		query += " where" + strings.Join(conditions, " and")
	}

	if filter.Per > 0 && filter.Page > 0 {
		query += fmt.Sprintf(" order by id desc limit $%d offset $%d", placeholder.Touch(), placeholder.Touch())
		args = append(args, filter.Per, (filter.Page-1)*filter.Per)
	}

	return query, args
}
