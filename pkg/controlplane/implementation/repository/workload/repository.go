package workload

import (
	"github.com/wernsiet/morchy/pkg/controlplane/domain/workload"
	dbutils "github.com/wernsiet/morchy/pkg/db.utils"
)

type Repository struct {
	db      dbutils.DB
	queries queries
}

func NewRepo(dbPool dbutils.DB) *Repository {
	return &Repository{
		db: dbPool,
	}
}

type WorkloadRepoFactory struct{}

func (w WorkloadRepoFactory) New(db dbutils.DB) workload.Repository {
	return NewRepo(db)
}
