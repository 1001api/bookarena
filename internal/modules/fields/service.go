package fields

import (
	db "1001api/bookarena/internal/database/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type FieldService interface {
	CreateField(req CreateFieldRequest) (uuid.UUID, error)
	GetFieldByID(id uuid.UUID) (db.GetFieldByIDRow, error)
	ListFields(limit, page int) ([]db.ListFieldsRow, error)
	UpdateField(id uuid.UUID, req UpdateFieldRequest) (uuid.UUID, error)
	DeleteField(id uuid.UUID) (uuid.UUID, error)
}

type service struct {
	repo FieldRepository
}

func NewFieldService(repo FieldRepository) FieldService {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateField(req CreateFieldRequest) (uuid.UUID, error) {
	input := db.CreateFieldParams{
		ImageUrl: pgtype.Text{
			String: req.ImageUrl,
			Valid:  req.ImageUrl != "",
		},
		Name: req.Name,
		Type: req.Type,
		Description: pgtype.Text{
			String: req.Description,
			Valid:  req.Description != "",
		},
		Location: pgtype.Text{
			String: req.Location,
			Valid:  req.Location != "",
		},
		LocationLat: pgtype.Float8{
			Float64: req.LocationLat,
			Valid:   req.LocationLat != 0,
		},
		LocationLon: pgtype.Float8{
			Float64: req.LocationLon,
			Valid:   req.LocationLon != 0,
		},
		PricePerHour: req.PricePerHour,
	}

	return s.repo.CreateField(input)
}

func (s *service) ListFields(limit, page int) ([]db.ListFieldsRow, error) {
	if limit < 1 {
		limit = 10
	}

	if page < 1 {
		page = 1
	}

	arg := db.ListFieldsParams{
		LimitCount:  int32(limit),
		OffsetCount: int32((page - 1) * limit),
	}

	return s.repo.ListFields(arg)
}

func (s *service) GetFieldByID(id uuid.UUID) (db.GetFieldByIDRow, error) {
	return s.repo.GetFieldByID(id)
}

func (s *service) UpdateField(id uuid.UUID, req UpdateFieldRequest) (uuid.UUID, error) {
	arg := db.UpdateFieldParams{
		ID: id,
		ImageUrl: pgtype.Text{
			String: req.ImageUrl,
			Valid:  req.ImageUrl != "",
		},
		Name: req.Name,
		Type: req.Type,
		Description: pgtype.Text{
			String: req.Description,
			Valid:  req.Description != "",
		},
		Location: pgtype.Text{
			String: req.Location,
			Valid:  req.Location != "",
		},
		LocationLat: pgtype.Float8{
			Float64: req.LocationLat,
			Valid:   req.LocationLat != 0,
		},
		LocationLon: pgtype.Float8{
			Float64: req.LocationLon,
			Valid:   req.LocationLon != 0,
		},
		PricePerHour: req.PricePerHour,
	}

	return s.repo.UpdateField(arg)
}

func (s *service) DeleteField(id uuid.UUID) (uuid.UUID, error) {
	return s.repo.DeleteField(id)
}
