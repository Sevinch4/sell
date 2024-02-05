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

type staffTariffRepo struct {
	DB *pgxpool.Pool
}

func NewStaffTariffRepo(DB *pgxpool.Pool) storage.IStaffTariffRepo {
	return &staffTariffRepo{
		DB: DB,
	}
}

func (s *staffTariffRepo) Create(ctx context.Context, tariff models.CreateStaffTariff) (string, error) {
	id := uuid.New().String()

	if _, err := s.DB.Exec(ctx, `INSERT INTO staff_tariffs 
	(id, name, tariff_type, amount_for_cash, amount_for_card) 
		VALUES ($1, $2, $3, $4, $5)`,
		id,
		tariff.Name,
		tariff.TariffType,
		tariff.AmountForCash,
		tariff.AmountForCard,
	); err != nil {
		log.Println("Error while inserting data:", err)
		return "", err
	}

	return id, nil
}

func (s *staffTariffRepo) GetStaffTariffByID(ctx context.Context, id models.PrimaryKey) (models.StaffTariff, error) {
	staffTariff := models.StaffTariff{}
	query := `SELECT id, name, tariff_type, amount_for_cash, amount_for_card, created_at, updated_at 
								FROM staff_tariffs WHERE id = $1 and deleted_at is null`
	err := s.DB.QueryRow(ctx, query, id.ID).Scan(
		&staffTariff.ID,
		&staffTariff.Name,
		&staffTariff.TariffType,
		&staffTariff.AmountForCash,
		&staffTariff.AmountForCard,
		&staffTariff.CreatedAt,
		&staffTariff.UpdatedAt,
	)
	if err != nil {
		log.Println("Error while selecting staff tariff by ID:", err)
		return models.StaffTariff{}, err
	}
	return staffTariff, nil
}

func (s *staffTariffRepo) GetStaffTariffList(ctx context.Context, request models.GetListRequest) (models.StaffTariffResponse, error) {
	var (
		staffTariffs []models.StaffTariff
		count        int
	)

	countQuery := `SELECT COUNT(*) FROM staff_tariffs where deleted_at is null`
	if request.Search != "" {
		countQuery += fmt.Sprintf(` and name ILIKE '%s'`, request.Search)
	}

	err := s.DB.QueryRow(ctx, countQuery).Scan(&count)
	if err != nil {
		log.Println("Error while scanning count of staff tariffs:", err)
		return models.StaffTariffResponse{}, err
	}

	query := `SELECT id, name, tariff_type, amount_for_cash, amount_for_card, created_at, updated_at FROM staff_tariffs where deleted_at is null`
	if request.Search != "" {
		query += fmt.Sprintf(` and name ILIKE '%s'`, request.Search)
	}
	query += ` order by created_at desc LIMIT $1 OFFSET $2 `

	rows, err := s.DB.Query(ctx, query, request.Limit, (request.Page-1)*request.Limit)
	if err != nil {
		log.Println("Error while querying staff tariffs:", err)
		return models.StaffTariffResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		staffTariff := models.StaffTariff{}
		err := rows.Scan(
			&staffTariff.ID,
			&staffTariff.Name,
			&staffTariff.TariffType,
			&staffTariff.AmountForCash,
			&staffTariff.AmountForCard,
			&staffTariff.CreatedAt,
			&staffTariff.UpdatedAt,
		)
		if err != nil {
			log.Println("Error while scanning row of staff tariffs:", err)
			return models.StaffTariffResponse{}, err
		}
		staffTariffs = append(staffTariffs, staffTariff)
	}

	return models.StaffTariffResponse{
		StaffTariffs: staffTariffs,
		Count:        count,
	}, nil
}

func (s *staffTariffRepo) UpdateStaffTariff(ctx context.Context, tariff models.UpdateStaffTariff) (string, error) {
	query := `UPDATE staff_tariffs SET name = $1, tariff_type = $2, amount_for_cash = $3, 
                         amount_for_card = $4, updated_at = NOW() WHERE id = $5 and deleted_at is null`

	_, err := s.DB.Exec(ctx, query,
		tariff.Name,
		tariff.TariffType,
		tariff.AmountForCash,
		tariff.AmountForCard,
		tariff.ID,
	)
	if err != nil {
		log.Println("Error while updating Staff Tariff:", err)
		return "", err
	}

	return tariff.ID, nil
}

func (s *staffTariffRepo) DeleteStaffTariff(ctx context.Context, id string) error {
	query := `UPDATE staff_tariffs SET deleted_at = NOW() WHERE id = $1`

	_, err := s.DB.Exec(ctx, query, id)
	if err != nil {
		log.Println("Error while deleting Staff Tariff:", err)
		return err
	}

	return nil
}
