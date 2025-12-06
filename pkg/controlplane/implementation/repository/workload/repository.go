package workload

import (
	"github.com/jackc/pgx/v5/pgxpool"
	dbutils "github.com/wernsiet/morchy/pkg/db.utils"
)

type Repository struct {
	db      dbutils.DB
	queries queries
}

func NewRepo(dbPool *pgxpool.Pool) *Repository {
	return &Repository{
		db: dbPool,
	}
}
