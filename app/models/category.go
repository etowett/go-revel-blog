package models

import (
	"context"
	"fmt"
	"go-revel-blog/app/db"
	"go-revel-blog/app/helpers"
	"strings"
)

const (
	createCategorySQL = `insert into categories (user_id, name, description created_at) values ($1, $2, $3, $4) returning id`
	getCategorySQL    = `select id, user_id, name, description created_at, updated_at from users`
	getCategoryByID   = getUsersSQL + ` where id=$1`
	updateCategorySQL = `update categories set (name, description, updated_at) = ($1, $2, $3) where id = $4`
	countCategorySQL  = `select count(id) from categories`
	deleteCategorySQL = `delete from categories where id=$1`
)

type (
	Category struct {
		SequentialIdentifier
		UserID      string `json:"user_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Timestamps
	}
)

func (category *Category) String() string {
	return fmt.Sprintf("Category(%s)", category.Name)
}

func (c *Category) All(
	ctx context.Context,
	db db.SQLOperations,
	filter *Filter,
) ([]*Category, error) {
	cats := make([]*Category, 0)

	query, args := c.buildQuery(
		getCategorySQL,
		filter,
	)

	rows, err := db.QueryContext(
		ctx,
		query,
		args...,
	)
	defer rows.Close()
	if err != nil {
		return cats, err
	}

	for rows.Next() {
		var cat Category
		err = rows.Scan(
			&cat.ID,
			&cat.UserID,
			&cat.Name,
			&cat.Description,
			&cat.CreatedAt,
			&cat.UpdatedAt,
		)
		if err != nil {
			return cats, err
		}
		cats = append(cats, &cat)
	}

	return cats, err
}

func (q *Category) Count(
	ctx context.Context,
	db db.SQLOperations,
	filter *Filter,
) (int, error) {
	query, args := q.buildQuery(
		countCategorySQL,
		&Filter{
			Term: filter.Term,
		},
	)
	var recordsCount int
	err := db.QueryRowContext(ctx, query, args...).Scan(&recordsCount)
	return recordsCount, err
}

func (q *Category) Delete(
	ctx context.Context,
	db db.SQLOperations,
	id int64,
) (int64, error) {
	res, err := db.ExecContext(ctx, deleteCategorySQL, id)
	if err != nil {
		return 0, err
	}

	rowsDeleted, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsDeleted, nil
}

func (c *Category) GetByID(
	ctx context.Context,
	db db.SQLOperations,
	id int64,
) (*Category, error) {
	row := db.QueryRowContext(ctx, getUserByID, id)
	return c.scan(row)
}

func (c *Category) Save(
	ctx context.Context,
	db db.SQLOperations,
) error {
	c.Timestamps.Touch()

	var err error
	if c.IsNew() {
		err := db.QueryRowContext(
			ctx,
			createCategorySQL,
			c.UserID,
			c.Name,
			c.Description,
			c.Timestamps.CreatedAt,
		).Scan(&c.ID)
		return err
	}
	_, err = db.ExecContext(
		ctx,
		updateCategorySQL,
		c.Name,
		c.Description,
		c.Timestamps.UpdatedAt,
		c.ID,
	)
	return err
}

func (c *Category) scan(
	row db.RowScanner,
) (*Category, error) {
	var category Category
	err := row.Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	return &category, err
}

func (c *Category) buildQuery(
	query string,
	filter *Filter,
) (string, []interface{}) {
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	placeholder := helpers.NewPlaceholder()

	if filter.Term != "" {
		likeStmt := make([]string, 0)
		columns := []string{"name", "description"}
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
