package postgres

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"sell/config"
	"sell/storage"
	"strings"

	_ "github.com/golang-migrate/migrate/v4/database"          //database is needed for migration
	_ "github.com/golang-migrate/migrate/v4/database/postgres" //postgres is used for database
	_ "github.com/golang-migrate/migrate/v4/source/file"       //file is needed for migration url
)

type Store struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, cfg config.Config) (storage.IStorage, error) {
	url := fmt.Sprintf(
		`postgres://%s:%s@%s:%s/%s?sslmode=disable`,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB)

	fmt.Println("url", url)
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		fmt.Println("error is while parsing config", err.Error())
		return nil, err
	}
	poolConfig.MaxConns = 100

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		fmt.Println("error is while connecting to db", err.Error())
		return nil, err
	}

	m, err := migrate.New("file://migrations/postgres/", url)
	if err != nil {
		fmt.Println("error is while migrating ", err.Error())
		return nil, err
	}

	if err = m.Up(); err != nil {
		if !strings.Contains(err.Error(), "no change") {
			version, dirty, err := m.Version()
			if err != nil {
				fmt.Println("error is while checking version and dirty", err.Error())
				return nil, err
			}

			if dirty {
				version--
				if err = m.Force(int(version)); err != nil {
					fmt.Println("error is while forcing", err.Error())
					return nil, err
				}
			}
			fmt.Println("ERROR in migrating", err.Error())
			return nil, err
		}
	}
	return &Store{
		Pool: pool,
	}, nil
}

func (s *Store) Close() {
	s.Pool.Close()
}

func (s *Store) StaffTariff() storage.IStaffTariffRepo {
	return NewStaffTariffRepo(s.Pool)
}

func (s *Store) Category() storage.ICategory {
	return NewCategoryRepo(s.Pool)
}

func (s *Store) Product() storage.IProducts {
	return NewProductRepo(s.Pool)
}

func (s *Store) Branch() storage.IBranchStorage {
	return NewBranchRepo(s.Pool)
}

func (s *Store) Sale() storage.ISaleStorage {
	return NewSaleRepo(s.Pool)
}

func (s *Store) Transaction() storage.ITransactionStorage {
	return NewTransactionRepo(s.Pool)

}

func (s *Store) Staff() storage.IStaffRepo {
	return NewStaffRepo(s.Pool)
}

func (s *Store) Repository() storage.IRepositoryRepo {
	return NewRepositoryRepo(s.Pool)
}

func (s *Store) Basket() storage.IBasketRepo {
	return NewBasketRepo(s.Pool)
}

func (s *Store) RTransaction() storage.IRepositoryTransactionRepo {
	return NewRepositoryTransactionRepo(s.Pool)
}
