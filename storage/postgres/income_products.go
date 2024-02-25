package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"sell/api/models"
	"sell/storage"
)

type incomeProductRepo struct {
	db *pgxpool.Pool
}

func NewIncomeProductsRepo(db *pgxpool.Pool) storage.IIncomeProductsStorage {
	return incomeProductRepo{db: db}
}

func (i incomeProductRepo) Create(ctx context.Context, incomeProduct models.CreateIncomeProduct) (string, error) {
	id := uuid.New()
	query := `insert into income_products(id, income_id, product_id, price, count) values($1, $2, $3, $4, $5)`
	if _, err := i.db.Exec(ctx, query,
		id,
		incomeProduct.IncomeID,
		incomeProduct.ProductID,
		incomeProduct.Price,
		incomeProduct.Count,
	); err != nil {
		fmt.Println("error is while inserting income products", err.Error())
		return "", err
	}
	return id.String(), nil
}

func (i incomeProductRepo) GetByID(ctx context.Context, id string) (models.IncomeProduct, error) {
	incomeProduct := models.IncomeProduct{}

	query := `select id, income_id, product_id, price, count, created_at::text, updated_at::text from income_products where id = $1 and deleted_at is null`

	if err := i.db.QueryRow(ctx, query, id).Scan(
		&incomeProduct.ID,
		&incomeProduct.IncomeID,
		&incomeProduct.ProductID,
		&incomeProduct.Price,
		&incomeProduct.Count,
		&incomeProduct.CreatedAt,
		&incomeProduct.UpdatedAt,
	); err != nil {
		fmt.Println("error is while selecting income products", err.Error())
		return models.IncomeProduct{}, err
	}
	return incomeProduct, nil
}

func (i incomeProductRepo) GetList(ctx context.Context, request models.IncomeProductRequest) (models.IncomeProductsResponse, error) {
	var (
		query, countQuery string
		filter            string
		pagination        string
		count             int
		page              = request.Page
		offset            = (page - 1) * request.Limit
		incomeProducts    []models.IncomeProduct
	)

	if request.ProductID != "" {
		filter += fmt.Sprintf(` and product_id = '%s'`, request.ProductID)
	}

	if request.IncomeID != "" {
		filter += fmt.Sprintf(` and income_id = '%s'`, request.IncomeID)
	}

	countQuery = `select count(1) from income_products where deleted_at is null ` + filter
	if err := i.db.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting count ", err.Error())
		return models.IncomeProductsResponse{}, err
	}

	pagination = ` ORDER BY created_at desc LIMIT $1 OFFSET $2 `
	query = `select id, income_id, product_id, price, count, created_at::text, updated_at::text from income_products where deleted_at is null ` + filter + pagination

	rows, err := i.db.Query(ctx, query, request.Limit, offset)
	fmt.Println("limit", request.Limit)
	fmt.Println("rows", rows)
	if err != nil {
		fmt.Println("error is while selecting all from incomes", err.Error())
		return models.IncomeProductsResponse{}, err
	}

	for rows.Next() {
		incomeProduct := models.IncomeProduct{}
		if err := rows.Scan(
			&incomeProduct.ID,
			&incomeProduct.IncomeID,
			&incomeProduct.ProductID,
			&incomeProduct.Price,
			&incomeProduct.Count,
			&incomeProduct.CreatedAt,
			&incomeProduct.UpdatedAt,
		); err != nil {
			fmt.Println("error is while scanning incomes", err.Error())
			return models.IncomeProductsResponse{}, err
		}
		incomeProducts = append(incomeProducts, incomeProduct)
	}

	return models.IncomeProductsResponse{
		IncomeProducts:        incomeProducts,
		CountOfIncomeProducts: count,
	}, nil

}

func (i incomeProductRepo) Update(ctx context.Context, income models.UpdateIncomeProduct) (string, error) {
	query := `update income_products set income_id = $1, product_id = $2, price = $3, count = $4 , updated_at = now() where id = $5`
	rowsAffected, err := i.db.Exec(ctx, query, &income.IncomeID, &income.ProductID, &income.Price, &income.Count, &income.ID)
	if err != nil {
		fmt.Println("error is while updating income products", err.Error())
		return "", err
	}

	if r := rowsAffected.RowsAffected(); r == 0 {
		fmt.Println("error is while rowsAffected", err.Error())
		return "", err
	}

	return income.ID, err
}

func (i incomeProductRepo) Delete(ctx context.Context, id string) error {
	query := `update income_products set deleted_at = now() where id = $1`
	if _, err := i.db.Exec(ctx, query, id); err != nil {
		fmt.Println("error is while deleting income", err.Error())
		return err
	}
	return nil
}
