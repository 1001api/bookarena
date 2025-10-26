package payments

import (
	db "1001api/bookarena/internal/database/generated"
	"1001api/bookarena/pkg"
	"encoding/hex"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type PaymentsService interface {
	CreatePayment(args db.CreatePaymentParams) (uuid.UUID, error)
	GetPaymentByID(id uuid.UUID) (db.GetPaymentByIDRow, error)
	ListPaymentsForAdmin(limit, page int) ([]db.ListPaymentsForAdminRow, error)
	ListPaymentsForUser(currentUserUUID uuid.UUID, limit, page int) ([]db.ListPaymentsForUserRow, error)
	UpdatePaymentStatus(args db.UpdatePaymentStatusParams) (uuid.UUID, error)
	DeletePayment(id uuid.UUID) (uuid.UUID, error)
}

type service struct {
	repository PaymentsRepository
	encKey     string
}

func NewPaymentsService(repository PaymentsRepository) PaymentsService {
	return &service{
		repository: repository,
		encKey:     viper.GetString("ENC_KEY"),
	}
}

func (s *service) CreatePayment(args db.CreatePaymentParams) (uuid.UUID, error) {
	return s.repository.CreatePayment(args)
}

func (s *service) GetPaymentByID(id uuid.UUID) (db.GetPaymentByIDRow, error) {
	payment, err := s.repository.GetPaymentByID(id)
	if err != nil {
		return db.GetPaymentByIDRow{}, err
	}

	// Decode the hex string into bytes
	cipherBytes, err := hex.DecodeString(strings.TrimPrefix(payment.UserEmail, "\\x"))
	if err != nil {
		log.Error().Err(err).Msg("failed to decode hex")
		return db.GetPaymentByIDRow{}, err
	}

	// Decrypt the email using the encryption key (see .env & pkg/common.go)
	de, err := pkg.Decrypt(cipherBytes, []byte(s.encKey))
	if err != nil {
		log.Error().Err(err).Msg("failed to decrypt user email")
		return db.GetPaymentByIDRow{}, err
	}

	payment.UserEmail = string(de)

	return payment, nil
}

func (s *service) ListPaymentsForAdmin(limit, page int) ([]db.ListPaymentsForAdminRow, error) {
	if limit < 1 {
		limit = 10
	}
	if page < 1 {
		page = 1
	}

	input := db.ListPaymentsForAdminParams{
		LimitCount:  int32(limit),
		OffsetCount: int32((page - 1) * limit),
	}

	result, err := s.repository.ListPaymentsForAdmin(input)
	if err != nil {
		return []db.ListPaymentsForAdminRow{}, err
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
			return []db.ListPaymentsForAdminRow{}, err
		}

		result[i].UserEmail = string(de)
	}

	return result, nil
}

func (s *service) ListPaymentsForUser(currentUserUUID uuid.UUID, limit, page int) ([]db.ListPaymentsForUserRow, error) {
	if limit < 1 {
		limit = 10
	}
	if page < 1 {
		page = 1
	}

	input := db.ListPaymentsForUserParams{
		UserUuid:    currentUserUUID,
		LimitCount:  int32(limit),
		OffsetCount: int32((page - 1) * limit),
	}

	return s.repository.ListPaymentsForUser(input)
}

func (s *service) UpdatePaymentStatus(args db.UpdatePaymentStatusParams) (uuid.UUID, error) {
	return s.repository.UpdatePaymentStatus(args)
}

func (s *service) DeletePayment(id uuid.UUID) (uuid.UUID, error) {
	return s.repository.DeletePayment(id)
}
