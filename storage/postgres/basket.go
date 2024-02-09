package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"sell/api/models"
	"sell/storage"
)

type basketRepo struct {
	DB *pgxpool.Pool
}

func NewBasketRepo(DB *pgxpool.Pool) storage.IBasketRepo {
	return &basketRepo{
		DB: DB,
	}
}

func (s *basketRepo) Create(ctx context.Context, basket models.CreateBasket) (string, error) {
	id := uuid.New().String()

	if _, err := s.DB.Exec(ctx, `INSERT INTO baskets 
		(id, sale_id, product_id, price, quantity)
			VALUES($1, $2, $3, $4, $5) `,
		id,
		basket.SaleID,
		basket.ProductID,
		basket.Price,
		basket.Quantity,
	); err != nil {
		log.Println("Error while inserting data:", err)
		return "", err
	}

	return id, nil
}

func (s *basketRepo) GetByID(ctx context.Context, id models.PrimaryKey) (models.Basket, error) {
	basket := models.Basket{}
	query := `SELECT id, sale_id, product_id, quantity, price, created_at, updated_at
				FROM baskets WHERE id = $1 and  deleted_at is null`
	err := s.DB.QueryRow(ctx, query, id.ID).Scan(
		&basket.ID,
		&basket.SaleID,
		&basket.ProductID,
		&basket.Quantity,
		&basket.Price,
		&basket.CreatedAt,
		&basket.UpdatedAt,
	)
	if err != nil {
		log.Println("Error while selecting basket by ID:", err)
		return models.Basket{}, err
	}
	return basket, nil
}

func (s *basketRepo) GetList(ctx context.Context, request models.GetListRequest) (models.BasketsResponse, error) {
	var (
		baskets = []models.Basket{}
		count   int
	)

	countQuery := `SELECT COUNT(*) FROM baskets where deleted_at is null`
	if request.Search != "" {
		countQuery += fmt.Sprintf(` and sale_id = '%s'`, request.Search)
	}

	err := s.DB.QueryRow(ctx, countQuery).Scan(&count)
	if err != nil {
		log.Println("Error while scanning count of repositories:", err)
		return models.BasketsResponse{}, err
	}

	query := `SELECT id, sale_id, product_id, quantity, price, created_at, updated_at
						FROM baskets where deleted_at is null`
	if request.Search != "" {
		query += fmt.Sprintf(` and sale_id = '%s'`, request.Search)
	}
	query += ` order by created_at desc LIMIT $1 OFFSET $2 `

	rows, err := s.DB.Query(ctx, query, request.Limit, (request.Page-1)*request.Limit)
	if err != nil {
		log.Println("Error while querying baskets:", err)
		return models.BasketsResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		basket := models.Basket{}
		err := rows.Scan(
			&basket.ID,
			&basket.SaleID,
			&basket.ProductID,
			&basket.Quantity,
			&basket.Price,
			&basket.CreatedAt,
			&basket.UpdatedAt,
		)
		if err != nil {
			log.Println("Error while scanning row of baskets:", err)
			return models.BasketsResponse{}, err
		}
		baskets = append(baskets, basket)
	}

	return models.BasketsResponse{
		Baskets: baskets,
		Count:   count,
	}, nil
}

func (s *basketRepo) Update(ctx context.Context, basket models.UpdateBasket) (string, error) {
	query := `UPDATE baskets SET sale_id = $1, product_id = $2, quantity = $3, price = $4, updated_at = NOW() WHERE id = $5`

	_, err := s.DB.Exec(ctx, query,
		&basket.SaleID,
		&basket.ProductID,
		&basket.Quantity,
		&basket.Price,
		&basket.ID,
	)
	if err != nil {
		log.Println("Error while updating Repository :", err)
		return "", err
	}

	return basket.ID, nil
}

func (s *basketRepo) Delete(ctx context.Context, id string) error {
	query := `UPDATE baskets SET deleted_at = NOW() WHERE id = $1`

	_, err := s.DB.Exec(ctx, query, id)
	if err != nil {
		log.Println("Error while deleting Repository :", err)
		return err
	}

	return nil
}
