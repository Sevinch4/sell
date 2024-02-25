package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"sell/api/models"
	"sell/storage"
)

type incomeRepo struct {
	db *pgxpool.Pool
}

func NewIncomeRepo(db *pgxpool.Pool) storage.IIncomeStorage {
	return &incomeRepo{db: db}
}

func (i *incomeRepo) Create(ctx context.Context, income models.CreateIncome) (string, error) {
	id := uuid.New()
	query := `insert into incomes(id, branch_id, price) values($1, $2, 0)`
	if _, err := i.db.Exec(ctx, query, id, income.BranchID); err != nil {
		fmt.Println("error is while inserting income", err.Error())
		return "", err
	}
	return id.String(), nil
}

func (i *incomeRepo) GetByID(ctx context.Context, id string) (models.Income, error) {
	income := models.Income{}

	query := `select id, branch_id, price, created_at, updated_at from incomes where id = $1 and deleted_at is null`

	if err := i.db.QueryRow(ctx, query, id).Scan(
		&income.ID,
		&income.BranchID,
		&income.Price,
		&income.CreatedAt,
		&income.UpdatedAt,
	); err != nil {
		fmt.Println("error is while selecting incomes", err.Error())
		return models.Income{}, err
	}
	return income, nil
}

func (i *incomeRepo) GetList(ctx context.Context, request models.IncomeGetListRequest) (models.IncomeResponse, error) {
	var (
		query, countQuery string
		filter            string
		pagination        string
		count             int
		page              = request.Page
		offset            = (page - 1) * request.Limit
		incomes           []models.Income
	)

	if request.BranchID != "" {
		filter += fmt.Sprintf(` and branch_id = '%s'`, request.BranchID)
	}
	countQuery = `select count(1) from incomes where deleted_at is null ` + filter
	if err := i.db.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting count ", err.Error())
		return models.IncomeResponse{}, err
	}

	pagination = ` ORDER BY created_at desc LIMIT $1 OFFSET $2 `
	query = `select id, branch_id, price, created_at, updated_at from incomes where deleted_at is null ` + filter + pagination

	rows, err := i.db.Query(ctx, query, request.Limit, offset)
	fmt.Println("limit", request.Limit)
	fmt.Println("rows", rows)
	if err != nil {
		fmt.Println("error is while selecting all from incomes", err.Error())
		return models.IncomeResponse{}, err
	}
	fmt.Println("queyr", query)

	for rows.Next() {
		fmt.Println("here")
		income := models.Income{}
		if err := rows.Scan(
			&income.ID,
			&income.BranchID,
			&income.Price,
			&income.CreatedAt,
			&income.UpdatedAt,
		); err != nil {
			fmt.Println("error is while scanning incomes", err.Error())
			return models.IncomeResponse{}, err
		}
		fmt.Println("income", income)
		incomes = append(incomes, income)
	}

	fmt.Println("incomes", incomes)
	return models.IncomeResponse{
		Incomes: incomes,
		Count:   count,
	}, nil
}

func (i *incomeRepo) Update(ctx context.Context, income models.UpdateIncome) (string, error) {
	query := `update incomes set branch_id = $1, price = $2, updated_at = now() where id = $3`
	rowsAffected, err := i.db.Exec(ctx, query, &income.BranchID, &income.Price, &income.ID)
	if err != nil {
		fmt.Println("error is while updating incomes", err.Error())
		return "", err
	}

	if r := rowsAffected.RowsAffected(); r == 0 {
		fmt.Println("error is while rowsAffected", err.Error())
		return "", err
	}

	return income.ID, err
}

func (i *incomeRepo) Delete(ctx context.Context, id string) error {
	query := `update incomes set deleted_at = now() where id = $1`
	if _, err := i.db.Exec(ctx, query, id); err != nil {
		fmt.Println("error is while deleting income", err.Error())
		return err
	}
	return nil
}
