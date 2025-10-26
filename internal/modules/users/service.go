package users

import (
	db "1001api/bookarena/internal/database/generated"
	"1001api/bookarena/pkg"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type UserService interface {
	InitializeRootUser() error
	GetUsers(arg db.GetUsersParams) ([]UserProfileResponse, error)
	UpdateUserLastLogin(id uuid.UUID) error
	IncrementFailedLogins(id uuid.UUID) error
	CheckUserExists(email, username string) (bool, error)
	CheckUserExistsByID(id uuid.UUID) (bool, error)
	CreateUser(params CreateUserRequest) (uuid.UUID, error)
	UpdateUser(targetID uuid.UUID, req UpdateUserRequest) error
	DeleteUser(id uuid.UUID) error
	RestoreUser(id uuid.UUID) error
	SearchUser(query string, page, limit int32) ([]UserProfileResponse, error)
	GetUserByIdentifier(identifier string) (db.GetUserByIdentifierRow, error)
	GetUserByID(id uuid.UUID) (UserProfileResponse, error)
}

type service struct {
	enckey []byte
	repo   UserRepository
}

func NewUserService(
	repo UserRepository,
	encKey string,
) UserService {
	return &service{
		repo:   repo,
		enckey: []byte(encKey),
	}
}

func (s *service) decryptFields(encEmail []byte) ([]byte, error) {
	email, err := pkg.Decrypt(encEmail, s.enckey)
	if err != nil {
		return nil, err
	}
	return email, nil
}

func (s *service) InitializeRootUser() error {
	username := viper.GetString("ROOT_USERNAME")
	password := viper.GetString("ROOT_PASSWORD")
	email := viper.GetString("ROOT_EMAIL")
	passwordHash, err := pkg.HashPassword(password)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash root password")
		return err
	}

	// check if root user exists
	exists, err := s.CheckUserExists(email, username)
	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		log.Error().Err(err).Msg("Failed to check root user")
		return err
	}

	// skip if exist
	if exists {
		return nil
	}

	if username == "" || passwordHash == "" || email == "" {
		return errors.New("ROOT_USERNAME, ROOT_PASSWORD, ROOT_EMAIL are required")
	}

	_, err = s.CreateUser(CreateUserRequest{
		Username:     username,
		PasswordHash: passwordHash,
		Email:        email,
		Role:         string(pkg.RoleSuperAdmin),
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create root user")
		return err
	}

	return nil
}

func (s *service) GetUsers(arg db.GetUsersParams) ([]UserProfileResponse, error) {
	result, err := s.repo.GetUsers(arg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get users")
		return nil, err
	}

	var res []UserProfileResponse

	for _, user := range result {
		de, err := s.decryptFields(user.Email)
		if err != nil {
			log.Error().Err(err).Msg("Failed to decrypt user fields")
			return nil, err
		}

		res = append(res, UserProfileResponse{
			ID:          user.ID,
			Username:    user.Username,
			Email:       string(de),
			Role:        user.RoleName,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			LastLoginAt: user.LastLoginAt,
		})
	}

	return res, nil
}

func (s *service) UpdateUserLastLogin(id uuid.UUID) error {
	return s.repo.UpdateLastLogin(id)
}

func (s *service) IncrementFailedLogins(id uuid.UUID) error {
	return s.repo.IncrementFailedLogins(id)
}

func (s *service) CheckUserExists(email, username string) (bool, error) {
	hashedEmail := pkg.HashValue(email)

	return s.repo.CheckUserExists(db.CheckUserExistsParams{
		EmailHash: hashedEmail,
		Username:  username,
	})
}

func (s *service) CreateUser(params CreateUserRequest) (uuid.UUID, error) {
	exists, err := s.CheckUserExists(params.Email, params.Username)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check user exists")
		return uuid.UUID{}, err
	}

	if exists {
		log.Error().Msg("Username or email already exists")
		return uuid.UUID{}, errors.New("username or email already exists")
	}

	emailHash := pkg.HashValue(params.Email)
	emailEnc, err := pkg.Encrypt([]byte(params.Email), s.enckey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to encrypt email")
		return uuid.UUID{}, err
	}

	userID, err := s.repo.CreateUser(db.CreateUserParams{
		EmailHash:    emailHash,
		EmailEnc:     emailEnc,
		Username:     params.Username,
		PasswordHash: params.PasswordHash,
		Role:         params.Role,
	})
	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		return uuid.UUID{}, err
	}

	return userID, nil
}

func (s *service) UpdateUser(targetID uuid.UUID, req UpdateUserRequest) error {
	current, err := s.GetUserByID(targetID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user for update")
		return err
	}

	exists, err := s.CheckUserExists(req.Email, req.Username)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check user exists:")
		return err
	}

	if exists && (req.Email != current.Email || req.Username != current.Username) {
		return errors.New("username or email already exists")
	}

	emailHash := pkg.HashValue(req.Email)

	emailEnc, err := pkg.Encrypt([]byte(req.Email), s.enckey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to encrypt email")
		return err
	}

	if err := s.repo.UpdateUser(db.UpdateUserParams{
		Username:  req.Username,
		EmailHash: emailHash,
		EmailEnc:  emailEnc,
		ID:        targetID,
	}); err != nil {
		log.Error().Err(err).Msg("Failed to update user")
		return err
	}

	return nil
}

func (s *service) DeleteUser(id uuid.UUID) error {
	if err := s.repo.DeleteUser(id); err != nil {
		log.Error().Err(err).Msg("Failed to delete user")
		return err
	}

	return nil
}

func (s *service) RestoreUser(id uuid.UUID) error {
	if err := s.repo.RestoreUser(id); err != nil {
		log.Error().Err(err).Msg("Failed to restore user")
		return err
	}

	return nil
}

func (s *service) SearchUser(query string, page, limit int32) ([]UserProfileResponse, error) {
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 20
	}

	results, err := s.repo.SearchUser(db.SearchUserParams{
		Query:       query,
		OffsetCount: (page - 1) * limit,
		LimitCount:  limit,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to search user")
		return nil, err
	}

	var res []UserProfileResponse

	for _, user := range results {
		de, err := s.decryptFields(user.Email)
		if err != nil {
			log.Error().Err(err).Msg("Failed to decrypt user fields")
			return nil, err
		}

		res = append(res, UserProfileResponse{
			ID:          user.ID,
			Username:    user.Username,
			Email:       string(de),
			Role:        user.RoleName,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			LastLoginAt: user.LastLoginAt,
		})
	}

	return res, nil
}

func (s *service) GetUserByIdentifier(identifier string) (db.GetUserByIdentifierRow, error) {
	parsedUUID, _ := uuid.Parse(identifier)

	user, err := s.repo.GetUserByIdentifier(db.GetUserByIdentifierParams{
		ID:        parsedUUID,
		EmailHash: pkg.HashValue(identifier),
		Username:  identifier,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user by identifier")
		return db.GetUserByIdentifierRow{}, err
	}

	de, err := s.decryptFields(user.Email)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decrypt user fields")
		return db.GetUserByIdentifierRow{}, err
	}

	user.Email = de

	return user, nil
}

func (s *service) GetUserByID(id uuid.UUID) (UserProfileResponse, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user by ID")
		return UserProfileResponse{}, err
	}

	de, err := s.decryptFields(user.Email)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decrypt user fields")
		return UserProfileResponse{}, err
	}

	return UserProfileResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       string(de),
		Role:        user.RoleName,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		LastLoginAt: user.LastLoginAt,
	}, nil
}

func (s *service) CheckUserExistsByID(id uuid.UUID) (bool, error) {
	return s.repo.CheckUserExistsByID(id)
}
