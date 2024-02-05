package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"sell/api/models"
	"sell/storage"
)

type categoryRepo struct {
	db *pgxpool.Pool
}

func NewCategoryRepo(db *pgxpool.Pool) storage.ICategory {
	return categoryRepo{db: db}
}

func (c categoryRepo) Create(ctx context.Context, category models.CreateCategory) (string, error) {
	id := uuid.New()
	query := `insert into categories (id, name, parent_id) values($1, $2, $3)`
	if _, err := c.db.Exec(ctx, query, id, category.Name, category.ParentID); err != nil {
		fmt.Println("error is while inserting data", err.Error())
		return "", err
	}
	return id.String(), nil
}

func (c categoryRepo) GetByID(ctx context.Context, id string) (models.Category, error) {
	category := models.Category{}
	query := `select id, name, parent_id, created_at, updated_at from categories where id = $1 and deleted_at is null`
	if err := c.db.QueryRow(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.ParentID,
		&category.CreatedAt,
		&category.UpdatedAt); err != nil {
		fmt.Println("error is while selecting by id", err.Error())
		return models.Category{}, err
	}
	return category, nil
}

func (c categoryRepo) GetList(ctx context.Context, request models.GetListRequest) (models.CategoryResponse, error) {
	var (
		page              = request.Page
		offset            = (page - 1) * request.Limit
		query, countQuery string
		count             = 0
		categories        = []models.Category{}
		search            = request.Search
	)
	countQuery = `select count(1) from categories where deleted_at is null `
	if search != "" {
		countQuery += fmt.Sprintf(` and name ilike '%%%s%%'`, search)
	}
	if err := c.db.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while scanning count", err.Error())
		return models.CategoryResponse{}, err
	}

	query = `select id, name, parent_id, created_at, updated_at from categories where deleted_at is null `
	if search != "" {
		query += fmt.Sprintf(` and name ilike '%%%s%%'`, search)
	}

	query += ` order by created_at desc LIMIT $1 OFFSET $2 `
	rows, err := c.db.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting all", err.Error())
		return models.CategoryResponse{}, err
	}

	for rows.Next() {
		category := models.Category{}
		if err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.ParentID,
			&category.CreatedAt,
			&category.UpdatedAt); err != nil {
			fmt.Println("error is while scanning category", err.Error())
			return models.CategoryResponse{}, err
		}
		categories = append(categories, category)
	}
	return models.CategoryResponse{
		Categories: categories,
		Count:      count,
	}, nil
}

func (c categoryRepo) Update(ctx context.Context, category models.UpdateCategory) (string, error) {
	query := `update categories set name = $1, parent_id = $2, updated_at = now() where id = $3`
	if _, err := c.db.Exec(ctx, query, &category.Name, &category.ParentID, &category.ID); err != nil {
		fmt.Println("error is while updating", err.Error())
		return "", err
	}
	return category.ID, nil
}

func (c categoryRepo) Delete(ctx context.Context, id string) error {
	query := `update categories set deleted_at = now() where id = $1`
	if _, err := c.db.Exec(ctx, query, id); err != nil {
		fmt.Println("error is while deleting", err.Error())
		return err
	}
	return nil
}
