package bookings

import (
	db "1001api/bookarena/internal/database/generated"
	"1001api/bookarena/internal/modules/fields"
	"1001api/bookarena/pkg"
	"encoding/hex"
	"errors"
	"math"
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
	fieldService fields.FieldService
	repo         BookingRepository
	encKey       string
}

func NewBookingService(fieldService fields.FieldService, repo BookingRepository) BookingService {
	return &service{
		fieldService: fieldService,
		repo:         repo,
		encKey:       viper.GetString("ENC_KEY"),
	}
}

func (s *service) CreateBooking(userID uuid.UUID, req CreateBookingRequest) (uuid.UUID, error) {
	// check if time is valid
	if req.StartTime.Before(time.Now()) {
		return uuid.UUID{}, errors.New("start time must be in the future")
	}

	if req.EndTime.Before(req.StartTime) {
		return uuid.UUID{}, errors.New("end time must be after start time")
	}

	// Check for booking conflict for current timeframe
	conflict, err := s.CheckBookingConflict(req.FieldID, req.StartTime, req.EndTime)
	if err != nil {
		log.Error().Err(err).Msg("failed to check booking conflict")
		return uuid.UUID{}, err
	}

	if conflict {
		return uuid.UUID{}, errors.New("booking conflict")
	}

	// Check if field exists
	field, err := s.fieldService.GetFieldByID(req.FieldID)
	if err != nil {
		return uuid.UUID{}, errors.New("no rows in result set")
	}

	// calculate total price per hour
	totalHour := req.EndTime.Sub(req.StartTime).Hours()
	pricePerHour := field.PricePerHour

	// round total hour up to nearest integer
	// this to ensure that if the total hour is less than 1, it will be rounded up to 1
	totalHour = math.Ceil(totalHour)

	totalPrice := int64(totalHour * float64(pricePerHour))

	input := db.CreateBookingParams{
		UserID:     userID,
		FieldID:    req.FieldID,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		TotalPrice: totalPrice,
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
