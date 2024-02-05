package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"sell/api/models"
	"sell/storage"
	"strconv"
)

type transactionRepo struct {
	db *pgxpool.Pool
}

func NewTransactionRepo(db *pgxpool.Pool) storage.ITransactionStorage {
	return transactionRepo{db: db}
}

func (t transactionRepo) Create(ctx context.Context, trans models.CreateTransaction) (string, error) {
	id := uuid.New()
	query := `insert into transactions 
    					(id, sale_id, staff_id, transaction_type, source_type, amount, description) 
						values ($1, $2, $3, $4, $5, $6, $7)`
	if _, err := t.db.Exec(ctx, query, id,
		trans.SaleID,
		trans.StaffID,
		trans.TransactionType,
		trans.SourceType,
		trans.Amount,
		trans.Description); err != nil {
		fmt.Println("error is while inserting data", err.Error())
		return "", err
	}
	return id.String(), nil
}

func (t transactionRepo) GetByID(ctx context.Context, id string) (models.Transaction, error) {
	trans := models.Transaction{}
	query := `select id, sale_id, staff_id, transaction_type, source_type, amount,
       						description, created_at, updated_at
							from transactions where deleted_at is null and id = $1`
	if err := t.db.QueryRow(ctx, query, id).Scan(
		&trans.ID,
		&trans.SaleID,
		&trans.StaffID,
		&trans.TransactionType,
		&trans.SourceType,
		&trans.Amount,
		&trans.Description,
		&trans.CreatedAt,
		&trans.UpdatedAt); err != nil {
		fmt.Println("error is while selecting by id", err.Error())
		return models.Transaction{}, err
	}
	return trans, nil
}

func (t transactionRepo) GetList(ctx context.Context, request models.TransactionGetListRequest) (models.TransactionResponse, error) {
	var (
		page              = request.Page
		offset            = (page - 1) * request.Limit
		transactions      = []models.Transaction{}
		fromAmount        = request.FromAmount
		toAmount          = request.ToAmount
		count             = 0
		query, countQuery string
	)

	countQuery = `select count(1) from transactions where deleted_at is null `
	if fromAmount != 0 && toAmount != 0 {
		countQuery += fmt.Sprintf(` and amount between %f and %f`, fromAmount, toAmount)
	} else if fromAmount != 0 {
		countQuery += ` and amount >= ` + strconv.FormatFloat(fromAmount, 'f', 2, 64)
	} else {
		countQuery += ` and amount <= ` + strconv.FormatFloat(toAmount, 'f', 2, 64)

	}
	if err := t.db.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		fmt.Println("error is while scanning row", err.Error())
		return models.TransactionResponse{}, err
	}

	query = `select id, sale_id, staff_id, transaction_type, source_type, amount,
       						description, created_at, updated_at from transactions where deleted_at is null `

	if fromAmount != 0 && toAmount != 0 {
		query += fmt.Sprintf(` and amount between %f and %f  order by amount asc, `, fromAmount, toAmount)
	} else if fromAmount != 0 {
		query += ` and amount >= ` + strconv.FormatFloat(fromAmount, 'f', 2, 64) + `  order by amount asc, `
	} else {
		query += ` and amount <= ` + strconv.FormatFloat(toAmount, 'f', 2, 64) + ` order by amount asc, `

	}

	query += ` created_at desc LIMIT $1 OFFSET $2 `

	rows, err := t.db.Query(ctx, query, request.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting all from transactions", err.Error())
		return models.TransactionResponse{}, err
	}

	for rows.Next() {
		trans := models.Transaction{}
		if err = rows.Scan(
			&trans.ID,
			&trans.SaleID,
			&trans.StaffID,
			&trans.TransactionType,
			&trans.SourceType,
			&trans.Amount,
			&trans.Description,
			&trans.CreatedAt,
			&trans.UpdatedAt); err != nil {
			fmt.Println("error is while scanning rows", err.Error())
			return models.TransactionResponse{}, err
		}
		transactions = append(transactions, trans)
	}
	return models.TransactionResponse{
		Transactions: transactions,
		Count:        count,
	}, nil
}

func (t transactionRepo) Update(ctx context.Context, transaction models.UpdateTransaction) (string, error) {
	query := `update transactions set sale_id = $1, staff_id = $2, transaction_type = $3, source_type = $4, amount = $5,
								description = $6, updated_at = now() 
                    			where id = $7`
	if _, err := t.db.Exec(ctx, query,
		&transaction.SaleID,
		&transaction.StaffID,
		&transaction.TransactionType,
		&transaction.SourceType,
		&transaction.Amount,
		&transaction.Description,
		&transaction.ID); err != nil {
		fmt.Println("error is while updating transaction", err.Error())
		return "", err
	}
	return transaction.ID, nil
}

func (t transactionRepo) Delete(ctx context.Context, id string) error {
	query := `update transactions set deleted_at = now() where id = $1`
	if _, err := t.db.Exec(ctx, query, id); err != nil {
		fmt.Println("error is while deleting", err.Error())
		return err
	}
	return nil
}
