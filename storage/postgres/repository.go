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

type repositoryRepo struct {
	DB *pgxpool.Pool
}

func NewRepositoryRepo(DB *pgxpool.Pool) storage.IRepositoryRepo {
	return &repositoryRepo{
		DB: DB,
	}
}

func (s *repositoryRepo) Create(ctx context.Context, repository models.CreateRepository) (string, error) {
	id := uuid.New()

	if _, err := s.DB.Exec(ctx, `INSERT INTO repositories 
    (id, product_id, branch_id, count) 
        VALUES ($1, $2, $3, $4)`,
		id,
		repository.ProductID,
		repository.BranchID,
		repository.Count,
	); err != nil {
		log.Println("Error while inserting data:", err)
		return "", err
	}

	return id.String(), nil
}

func (s *repositoryRepo) GetByID(ctx context.Context, id models.PrimaryKey) (models.Repository, error) {
	repository := models.Repository{}
	query := `SELECT id, product_id, branch_id, count, created_at, updated_at 
							FROM repositories WHERE id = $1 and deleted_at is null
`
	err := s.DB.QueryRow(ctx, query, id.ID).Scan(
		&repository.ID,
		&repository.ProductID,
		&repository.BranchID,
		&repository.Count,
		&repository.CreatedAt,
		&repository.UpdatedAt,
	)
	if err != nil {
		log.Println("Error while selecting repository by ID:", err)
		return models.Repository{}, err
	}
	return repository, nil
}

func (s *repositoryRepo) GetList(ctx context.Context, request models.GetListRequest) (models.RepositoriesResponse, error) {
	var (
		repositories = []models.Repository{}
		count        int
	)

	countQuery := `SELECT COUNT(*) FROM repositories where deleted_at is null`
	if request.Search != "" {
		countQuery += fmt.Sprintf(` and branch_id = '%s'`, request.Search)
	}

	err := s.DB.QueryRow(ctx, countQuery).Scan(&count)
	if err != nil {
		log.Println("Error while scanning count of repositories:", err)
		return models.RepositoriesResponse{}, err
	}

	query := `SELECT id, product_id, branch_id, count, created_at, updated_at FROM repositories where deleted_at is null`
	if request.Search != "" {
		query += fmt.Sprintf(` and branch_id = '%s'`, request.Search)
	}
	query += ` order by created_at desc LIMIT $1 OFFSET $2 `

	rows, err := s.DB.Query(ctx, query, request.Limit, (request.Page-1)*request.Limit)
	if err != nil {
		log.Println("Error while querying repositories:", err)
		return models.RepositoriesResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		repository := models.Repository{}
		err := rows.Scan(
			&repository.ID,
			&repository.ProductID,
			&repository.BranchID,
			&repository.Count,
			&repository.CreatedAt,
			&repository.UpdatedAt,
		)
		if err != nil {
			log.Println("Error while scanning row of repositories:", err)
			return models.RepositoriesResponse{}, err
		}
		repositories = append(repositories, repository)
	}

	return models.RepositoriesResponse{
		Repositories: repositories,
		Count:        count,
	}, nil
}

func (s *repositoryRepo) Update(ctx context.Context, repository models.UpdateRepository) (string, error) {
	query := `UPDATE repositories SET branch_id = $1, product_id = $2, count = $3, updated_at = NOW() WHERE id = $4`

	_, err := s.DB.Exec(ctx, query,
		&repository.BranchID,
		&repository.ProductID,
		&repository.Count,
		&repository.ID,
	)
	if err != nil {
		log.Println("Error while updating Repository :", err)
		return "", err
	}

	return repository.ID, nil
}

func (s *repositoryRepo) Delete(ctx context.Context, id string) error {
	query := `UPDATE repositories SET deleted_at = NOW() WHERE id = $1`

	_, err := s.DB.Exec(ctx, query, id)
	if err != nil {
		log.Println("Error while deleting Repository :", err)
		return err
	}

	return nil
}
