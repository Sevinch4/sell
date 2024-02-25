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

type repositoryTransactionRepo struct {
	DB *pgxpool.Pool
}

func NewRepositoryTransactionRepo(DB *pgxpool.Pool) storage.IRepositoryTransactionRepo {
	return &repositoryTransactionRepo{
		DB: DB,
	}
}

func (s *repositoryTransactionRepo) Create(ctx context.Context, rtransaction models.CreateRepositoryTransaction) (string, error) {
	id := uuid.New().String()
	fmt.Println("id", id)
	fmt.Println("prod id", rtransaction.ProductID)

	if _, err := s.DB.Exec(ctx, `INSERT INTO repository_transactions
		(id, product_id, repository_transaction_type, price, quantity)
			VALUES($1, $2, $3, $4, $5)`,
		id,
		rtransaction.ProductID,
		rtransaction.RepositoryTransactionType,
		rtransaction.Price,
		rtransaction.Quantity,
	); err != nil {
		log.Println("Error while inserting data:", err)
		return "", err
	}

	return id, nil
}

func (s *repositoryTransactionRepo) GetByID(ctx context.Context, id models.PrimaryKey) (models.RepositoryTransaction, error) {
	rtransaction := models.RepositoryTransaction{}
	query := `SELECT id, product_id, repository_transaction_type, price, quantity, created_at, updated_at 
							FROM repository_transactions WHERE id = $1 and deleted_at is null
`

	err := s.DB.QueryRow(ctx, query, id.ID).Scan(
		&rtransaction.ID,
		&rtransaction.ProductID,
		&rtransaction.RepositoryTransactionType,
		&rtransaction.Price,
		&rtransaction.Quantity,
		&rtransaction.CreatedAt,
		&rtransaction.UpdatedAt,
	)
	if err != nil {
		log.Println("Error while selecting repository by ID:", err)
		return models.RepositoryTransaction{}, err
	}

	return rtransaction, nil
}

func (s *repositoryTransactionRepo) GetList(ctx context.Context, req models.GetListRequest) (models.RepositoryTransactionsResponse, error) {
	var (
		rtransactions []models.RepositoryTransaction
		count         int
	)

	countQuery := `SELECT COUNT(*) FROM repository_transactions where deleted_at is null `
	if req.Search != "" {
		fmt.Println("req.search", req.Search)
		countQuery += fmt.Sprintf(` and product_id = '%s'`, req.Search)
		fmt.Println("count query", countQuery)
	}

	err := s.DB.QueryRow(ctx, countQuery).Scan(&count)
	if err != nil {
		log.Println("Error while scanning count of repository_transactions:", err)
		return models.RepositoryTransactionsResponse{}, err
	}

	query := `SELECT id, product_id, repository_transaction_type, price, quantity, created_at, updated_at 
							FROM repository_transactions where deleted_at is null
`
	if req.Search != "" {
		query += fmt.Sprintf(` and product_id = '%s'`, req.Search)
	}
	query += ` order by created_at desc LIMIT $1 OFFSET $2 `

	fmt.Println("query", query)

	rows, err := s.DB.Query(ctx, query, req.Limit, (req.Page-1)*req.Limit)
	if err != nil {
		log.Println("Error while querying repository_transactions:", err)
		return models.RepositoryTransactionsResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		rtransaction := models.RepositoryTransaction{}
		err := rows.Scan(
			&rtransaction.ID,
			&rtransaction.ProductID,
			&rtransaction.RepositoryTransactionType,
			&rtransaction.Price,
			&rtransaction.Quantity,
			&rtransaction.CreatedAt,
			&rtransaction.UpdatedAt,
		)
		if err != nil {
			log.Println("Error while scanning row of repository_transactions:", err)
			return models.RepositoryTransactionsResponse{}, err
		}
		rtransactions = append(rtransactions, rtransaction)
	}

	return models.RepositoryTransactionsResponse{
		RepositoryTransactions: rtransactions,
		Count:                  count,
	}, nil
}

func (s *repositoryTransactionRepo) Update(ctx context.Context, transaction models.UpdateRepositoryTransaction) (string, error) {
	query := `UPDATE repository_transactions SET  product_id = $1, repository_transaction_type = $2, 
                                   price = $3, quantity = $4, updated_at = NOW() WHERE id = $5
`

	_, err := s.DB.Exec(ctx, query,
		&transaction.ProductID,
		&transaction.RepositoryTransactionType,
		&transaction.Price,
		&transaction.Quantity,
		&transaction.ID,
	)
	if err != nil {
		log.Println("Error while repository_transactions Repository :", err)
		return "", err
	}

	return transaction.ID, nil
}

func (s *repositoryTransactionRepo) Delete(ctx context.Context, id string) error {
	query := `UPDATE repository_transactions SET deleted_at = NOW() WHERE id = $1`

	_, err := s.DB.Exec(ctx, query, id)
	if err != nil {
		log.Println("Error while deleting repository_transactions ", err)
		return err
	}

	return nil
}
