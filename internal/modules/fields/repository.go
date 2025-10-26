package fields

import (
	db "1001api/bookarena/internal/database/generated"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FieldRepository interface {
	CreateField(arg db.CreateFieldParams) (uuid.UUID, error)
	GetFieldByID(id uuid.UUID) (db.GetFieldByIDRow, error)
	ListFields(arg db.ListFieldsParams) ([]db.ListFieldsRow, error)
	UpdateField(arg db.UpdateFieldParams) (uuid.UUID, error)
	DeleteField(id uuid.UUID) (uuid.UUID, error)
}

type repository struct {
	queries *db.Queries
}

func NewFieldRepository(pool *pgxpool.Pool) FieldRepository {
	return &repository{
		queries: db.New(pool),
	}
}

func (r *repository) CreateField(arg db.CreateFieldParams) (uuid.UUID, error) {
	return r.queries.CreateField(context.Background(), arg)
}

func (r *repository) GetFieldByID(id uuid.UUID) (db.GetFieldByIDRow, error) {
	return r.queries.GetFieldByID(context.Background(), id)
}

func (r *repository) ListFields(arg db.ListFieldsParams) ([]db.ListFieldsRow, error) {
	return r.queries.ListFields(context.Background(), arg)
}

func (r *repository) UpdateField(arg db.UpdateFieldParams) (uuid.UUID, error) {
	return r.queries.UpdateField(context.Background(), arg)
}

func (r *repository) DeleteField(id uuid.UUID) (uuid.UUID, error) {
	return r.queries.DeleteField(context.Background(), id)
}
