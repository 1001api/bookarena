package payments

import (
	db "1001api/bookarena/internal/database/generated"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentsRepository interface {
	CreatePayment(args db.CreatePaymentParams) (uuid.UUID, error)
	GetPaymentByID(id uuid.UUID) (db.GetPaymentByIDRow, error)
	ListPaymentsForAdmin(args db.ListPaymentsForAdminParams) ([]db.ListPaymentsForAdminRow, error)
	ListPaymentsForUser(args db.ListPaymentsForUserParams) ([]db.ListPaymentsForUserRow, error)
	UpdatePaymentStatus(args db.UpdatePaymentStatusParams) (uuid.UUID, error)
	DeletePayment(id uuid.UUID) (uuid.UUID, error)
}

type repository struct {
	queries *db.Queries
}

func NewPaymentsRepository(pool *pgxpool.Pool) PaymentsRepository {
	return &repository{queries: db.New(pool)}
}

func (r *repository) CreatePayment(args db.CreatePaymentParams) (uuid.UUID, error) {
	return r.queries.CreatePayment(context.Background(), args)
}

func (r *repository) GetPaymentByID(id uuid.UUID) (db.GetPaymentByIDRow, error) {
	return r.queries.GetPaymentByID(context.Background(), id)
}

func (r *repository) ListPaymentsForAdmin(args db.ListPaymentsForAdminParams) ([]db.ListPaymentsForAdminRow, error) {
	return r.queries.ListPaymentsForAdmin(context.Background(), args)
}

func (r *repository) ListPaymentsForUser(args db.ListPaymentsForUserParams) ([]db.ListPaymentsForUserRow, error) {
	return r.queries.ListPaymentsForUser(context.Background(), args)
}

func (r *repository) UpdatePaymentStatus(args db.UpdatePaymentStatusParams) (uuid.UUID, error) {
	return r.queries.UpdatePaymentStatus(context.Background(), args)
}

func (r *repository) DeletePayment(id uuid.UUID) (uuid.UUID, error) {
	return r.queries.DeletePayment(context.Background(), id)
}
