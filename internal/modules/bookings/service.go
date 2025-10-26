package bookings

import (
	db "1001api/bookarena/internal/database/generated"
	"1001api/bookarena/pkg"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type BookingService interface {
	CreateBooking(userID uuid.UUID, req CreateBookingRequest) (uuid.UUID, error)
	GetBookingByID(id uuid.UUID) (db.GetBookingByIDRow, error)
	ListBookings(limit, page int) ([]db.ListBookingsRow, error)
	ListBookingsByUser(userID uuid.UUID, limit, page int) ([]db.ListBookingsByUserRow, error)
	ListBookingsByField(fieldID uuid.UUID, limit, page int) ([]db.ListBookingsByFieldRow, error)
	DeleteBooking(id uuid.UUID) (uuid.UUID, error)
	CheckBookingConflict(fieldID uuid.UUID, startTime, endTime time.Time) (bool, error)
}

type service struct {
	repo   BookingRepository
	encKey string
}

func NewBookingService(repo BookingRepository) BookingService {
	return &service{
		repo:   repo,
		encKey: viper.GetString("ENC_KEY"),
	}
}

func (s *service) CreateBooking(userID uuid.UUID, req CreateBookingRequest) (uuid.UUID, error) {
	// Check for booking conflict for current timeframe
	if _, err := s.CheckBookingConflict(req.FieldID, req.StartTime, req.EndTime); err != nil {
		return uuid.UUID{}, err
	}

	input := db.CreateBookingParams{
		UserID:    userID,
		FieldID:   req.FieldID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	return s.repo.CreateBooking(input)
}

func (s *service) GetBookingByID(id uuid.UUID) (db.GetBookingByIDRow, error) {
	bookings, err := s.repo.GetBookingByID(id)
	if err != nil {
		return db.GetBookingByIDRow{}, err
	}

	// Decode the hex string into bytes
	cipherBytes, err := hex.DecodeString(strings.TrimPrefix(bookings.UserEmail, "\\x"))
	if err != nil {
		log.Error().Err(err).Msg("failed to decode hex")
		return db.GetBookingByIDRow{}, err
	}

	// Decrypt the email using the encryption key (see .env & pkg/common.go)
	de, err := pkg.Decrypt(cipherBytes, []byte(s.encKey))
	if err != nil {
		log.Error().Err(err).Msg("failed to decrypt user email")
		return db.GetBookingByIDRow{}, err
	}

	bookings.UserEmail = string(de)

	return bookings, nil
}

func (s *service) ListBookings(limit, page int) ([]db.ListBookingsRow, error) {
	if limit < 1 {
		limit = 10
	}
	if page < 1 {
		page = 1
	}

	input := db.ListBookingsParams{
		LimitCount:  int32(limit),
		OffsetCount: int32((page - 1) * limit),
	}

	result, err := s.repo.ListBookings(input)
	if err != nil {
		return []db.ListBookingsRow{}, err
	}

	for i := range result {
		// Decode the hex string into bytes
		cipherBytes, err := hex.DecodeString(strings.TrimPrefix(result[i].UserEmail, "\\x"))
		if err != nil {
			log.Error().Err(err).Msg("failed to decode hex")
			return nil, err
		}

		// Decrypt the email using the encryption key (see .env & pkg/common.go)
		de, err := pkg.Decrypt(cipherBytes, []byte(s.encKey))
		if err != nil {
			log.Error().Err(err).Msg("failed to decrypt user email")
			return []db.ListBookingsRow{}, err
		}

		result[i].UserEmail = string(de)
	}

	return result, nil
}

func (s *service) ListBookingsByUser(userID uuid.UUID, limit, page int) ([]db.ListBookingsByUserRow, error) {
	if limit < 1 {
		limit = 10
	}
	if page < 1 {
		page = 1
	}

	input := db.ListBookingsByUserParams{
		UserID:      userID,
		LimitCount:  int32(limit),
		OffsetCount: int32((page - 1) * limit),
	}

	return s.repo.ListBookingsByUser(input)
}

func (s *service) ListBookingsByField(fieldID uuid.UUID, limit, page int) ([]db.ListBookingsByFieldRow, error) {
	if limit < 1 {
		limit = 10
	}
	if page < 1 {
		page = 1
	}

	input := db.ListBookingsByFieldParams{
		FieldUuid:   fieldID,
		LimitCount:  int32(limit),
		OffsetCount: int32((page - 1) * limit),
	}

	result, err := s.repo.ListBookingsByField(input)
	if err != nil {
		return []db.ListBookingsByFieldRow{}, err
	}

	for i := range result {
		// Decode the hex string into bytes
		cipherBytes, err := hex.DecodeString(strings.TrimPrefix(result[i].UserEmail, "\\x"))
		if err != nil {
			log.Error().Err(err).Msg("failed to decode hex")
			return nil, err
		}

		// Decrypt the email using the encryption key (see .env & pkg/common.go)
		de, err := pkg.Decrypt(cipherBytes, []byte(s.encKey))
		if err != nil {
			log.Error().Err(err).Msg("failed to decrypt user email")
			return []db.ListBookingsByFieldRow{}, err
		}

		result[i].UserEmail = string(de)
	}

	return result, nil
}

func (s *service) DeleteBooking(id uuid.UUID) (uuid.UUID, error) {
	return s.repo.DeleteBooking(id)
}

func (s *service) CheckBookingConflict(fieldID uuid.UUID, startTime, endTime time.Time) (bool, error) {
	if endTime.Before(startTime) {
		return false, errors.New("end time must be after start time")
	}

	input := db.CheckBookingConflictParams{
		FieldID:   fieldID,
		StartTime: startTime,
		EndTime:   endTime,
	}
	return s.repo.CheckBookingConflict(input)
}
