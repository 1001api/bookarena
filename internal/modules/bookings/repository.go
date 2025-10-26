package bookings

import (
	db "1001api/bookarena/internal/database/generated"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingRepository interface {
	CreateBooking(args db.CreateBookingParams) (uuid.UUID, error)
	GetBookingByID(id uuid.UUID) (db.GetBookingByIDRow, error)
	ListBookings(args db.ListBookingsParams) ([]db.ListBookingsRow, error)
	ListBookingsByUser(args db.ListBookingsByUserParams) ([]db.ListBookingsByUserRow, error)
	ListBookingsByField(args db.ListBookingsByFieldParams) ([]db.ListBookingsByFieldRow, error)
	DeleteBooking(id uuid.UUID) (uuid.UUID, error)
	CheckBookingConflict(args db.CheckBookingConflictParams) (bool, error)
}

type repository struct {
	queries *db.Queries
}

func NewBookingRepository(pool *pgxpool.Pool) BookingRepository {
	return &repository{
		queries: db.New(pool),
	}
}

func (r *repository) CreateBooking(args db.CreateBookingParams) (uuid.UUID, error) {
	return r.queries.CreateBooking(context.Background(), args)
}

func (r *repository) GetBookingByID(id uuid.UUID) (db.GetBookingByIDRow, error) {
	return r.queries.GetBookingByID(context.Background(), id)
}

func (r *repository) ListBookings(args db.ListBookingsParams) ([]db.ListBookingsRow, error) {
	return r.queries.ListBookings(context.Background(), args)
}

func (r *repository) ListBookingsByUser(args db.ListBookingsByUserParams) ([]db.ListBookingsByUserRow, error) {
	return r.queries.ListBookingsByUser(context.Background(), args)
}

func (r *repository) ListBookingsByField(args db.ListBookingsByFieldParams) ([]db.ListBookingsByFieldRow, error) {
	return r.queries.ListBookingsByField(context.Background(), args)
}

func (r *repository) DeleteBooking(id uuid.UUID) (uuid.UUID, error) {
	return r.queries.DeleteBooking(context.Background(), id)
}

func (r *repository) CheckBookingConflict(args db.CheckBookingConflictParams) (bool, error) {
	return r.queries.CheckBookingConflict(context.Background(), args)
}
