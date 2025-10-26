package auth

import (
	db "1001api/bookarena/internal/database/generated"
	"1001api/bookarena/internal/modules/users"
	"1001api/bookarena/pkg"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type AuthService interface {
	Login(req LoginRequest) (*LoginResponse, error)
	Register(req RegisterRequest) (uuid.UUID, error)
	ValidateToken(tokenString string) (*Claims, error)
	GenerateToken(user db.GetUserByIdentifierRow) (string, error)
	GenerateRefreshToken(user db.GetUserByIdentifierRow) (string, error)
	RefreshToken(req RefreshRequest) (*LoginResponse, error)
	CheckEmailValidity(email string) bool
	ValidateUserExist(id uuid.UUID) (bool, error)
}

type service struct {
	userService   users.UserService
	jwtSecret     []byte
	tokenExpiry   time.Duration
	refreshExpiry time.Duration
	enckey        []byte
}

func NewAuthService(userService users.UserService, encKey string) AuthService {
	viper.SetDefault("JWT_EXPIRY", 15)
	viper.SetDefault("JWT_REFRESH_EXPIRY", 720)

	return &service{
		userService:   userService,
		jwtSecret:     []byte(viper.GetString("JWT_SECRET")),
		tokenExpiry:   time.Duration(viper.GetInt("JWT_EXPIRY")) * time.Minute,
		refreshExpiry: time.Duration(viper.GetInt("JWT_REFRESH_EXPIRY")) * time.Minute,
		enckey:        []byte(encKey),
	}
}

func (s *service) Login(req LoginRequest) (*LoginResponse, error) {
	if req.Identifier == "" || req.Password == "" {
		return nil, errors.New("username and password are required")
	}

	// Retrieve the user from the database based on the provided identifier (email or username).
	user, err := s.userService.GetUserByIdentifier(req.Identifier)
	if err != nil {
		log.Warn().Str("identifier", req.Identifier).Msg("Login attempt with invalid credentials")
		return nil, errors.New("invalid credentials")
	}

	// Account lock check: To prevent brute-force attacks, this check verifies if the user's
	// account is temporarily locked due to multiple failed login attempts.
	if cast.ToTime(user.LockedUntil).After(time.Now()) {
		log.Warn().Str("user_id", user.ID.String()).Msg("Login attempt on a locked account")
		return nil, errors.New("account is temporarily locked due to multiple failed login attempts, please wait a few minutes")
	}

	if err := pkg.VerifyPassword(user.PasswordHash.String, req.Password); err != nil {
		s.userService.IncrementFailedLogins(user.ID)
		log.Warn().Str("user_id", user.ID.String()).Msg("Invalid password provided for user")
		return nil, errors.New("invalid credentials")
	}

	// Upon successful authentication, a new access token and refresh token are generated for the user.
	accessToken, err := s.GenerateToken(user)
	if err != nil {
		log.Error().Err(err).Str("user_id", user.ID.String()).Msg("Failed to generate access token")
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}
	refreshToken, err := s.GenerateRefreshToken(user)
	if err != nil {
		log.Error().Err(err).Str("user_id", user.ID.String()).Msg("Failed to generate refresh token")
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Background post-login tasks: To ensure a responsive user experience, tasks such as
	// updating the last login timestamp and creating an audit log are performed in the
	// background. This avoids blocking the main execution thread and allows the user to
	// proceed without delay.
	go func(userID uuid.UUID) {
		s.userService.UpdateUserLastLogin(userID)
	}(user.ID)

	log.Info().Str("user_id", user.ID.String()).Msg("User logged in successfully")

	return &LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) Register(req RegisterRequest) (uuid.UUID, error) {
	// Securely hash the user's password before storing it.
	hashed, err := pkg.HashPassword(req.Password)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password during registration")
		return uuid.UUID{}, err
	}

	createdID, err := s.userService.CreateUser(users.CreateUserRequest{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashed),
		Role:         req.Role,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user during registration")
		return uuid.UUID{}, err
	}

	log.Info().Str("user_id", createdID.String()).Msg("User registered successfully")

	return createdID, nil
}

func (s *service) GenerateToken(user db.GetUserByIdentifierRow) (string, error) {
	// Set the token's expiration time based on the configured token expiry duration.
	expiresAt := time.Now().Add(s.tokenExpiry)

	// Create the JWT claims, which include user-specific data and standard JWT fields.
	// These claims are embedded in the token and can be used by the client application
	// to identify the user and their permissions.
	claims := &Claims{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     string(user.Email),
		Role:      user.RoleName,
		TokenType: pkg.TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    pkg.TokenIssuer,
			Subject:   user.ID.String(),
		},
	}

	// Create a new token with the specified claims and signing method. The HS256 algorithm
	// is used here, which is a common choice for symmetric key signing.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key to generate the final token string. If the signing
	// process fails, an error is logged, and an empty string is returned.
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		log.Error().Err(err).Str("user_id", user.ID.String()).Msg("Failed to sign access token")
		return "", err
	}

	log.Info().Str("user_id", user.ID.String()).Msg("Access token generated successfully")

	return tokenString, nil
}

func (s *service) GenerateRefreshToken(user db.GetUserByIdentifierRow) (string, error) {
	// Set the token's expiration time based on the configured refresh token expiry duration.
	expiresAt := time.Now().Add(s.refreshExpiry)

	// Create the JWT claims for the refresh token. These claims are similar to the access
	// token's claims but are explicitly marked with a "refresh" token type to distinguish
	// them from access tokens.
	claims := &Claims{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     string(user.Email),
		Role:      user.RoleName,
		TokenType: pkg.TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    pkg.TokenIssuer,
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key to generate the final token string.
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		log.Error().Err(err).Str("user_id", user.ID.String()).Msg("Failed to sign refresh token")
		return "", err
	}

	log.Info().Str("user_id", user.ID.String()).Msg("Refresh token generated successfully")

	return tokenString, nil
}

// RefreshToken handles the process of issuing a new access and refresh token pair
// using a valid refresh token. It first validates the provided refresh token, and if
// successful, generates new tokens for the user. This allows users to maintain their
// session without needing to re-authenticate.
func (s *service) RefreshToken(req RefreshRequest) (*LoginResponse, error) {
	// Validate the refresh token to ensure it is authentic and not expired.
	claims, err := s.ValidateToken(req.RefreshToken)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to validate refresh token")
		return nil, err
	}

	// Ensure that the token provided is a refresh token. This check prevents access tokens
	// from being used to refresh a session, which would be a security risk.
	if claims.TokenType != "refresh" {
		log.Warn().Str("token_type", claims.TokenType).Msg("Invalid token type for refreshing session")
		return nil, errors.New("invalid token type")
	}

	// Retrieve the user associated with the refresh token. This ensures that the user
	// still exists in the system before issuing new tokens.
	user, err := s.userService.GetUserByIdentifier(claims.UserID.String())
	if err != nil {
		log.Warn().Str("user_id", claims.UserID.String()).Msg("User not found for refresh token")
		return nil, errors.New("user not found")
	}

	// Generate a new access token for the user. This token will have a new expiration
	// time and will be used for subsequent API requests.
	accessToken, err := s.GenerateToken(user)
	if err != nil {
		log.Error().Err(err).Str("user_id", user.ID.String()).Msg("Failed to generate new access token during refresh")
		return nil, err
	}

	// Generate a new refresh token. This allows the user to continue refreshing their
	// session in the future without needing to log in again.
	newRefreshToken, err := s.GenerateRefreshToken(user)
	if err != nil {
		log.Error().Err(err).Str("user_id", user.ID.String()).Msg("Failed to generate new refresh token during refresh")
		return nil, err
	}

	log.Info().Str("user_id", user.ID.String()).Msg("Token refreshed successfully")

	return &LoginResponse{
		Token:        accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// ValidateToken verifies the authenticity and validity of a JWT token. It checks the
// token's signature, ensures that the signing method is as expected, and extracts
// the claims if the token is valid. This function is a critical part of securing
// endpoints and protecting user data.
func (s *service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Warn().Str("signing_method", fmt.Sprintf("%v", token.Header["alg"])).Msg("Unexpected signing method")
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return s.jwtSecret, nil
		},
	)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to parse or validate token")
		return nil, err
	}

	// If the token is valid, the claims are extracted and returned. This allows the
	// application to access the user's information and permissions stored within the token.
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		log.Info().Str("user_id", claims.UserID.String()).Msg("Token validated successfully")
		return claims, nil
	}

	// If the token is not valid for any other reason, an "invalid token" error is returned.
	// This is a generic error message to avoid leaking specific details about the failure.
	log.Warn().Msg("Invalid token provided")
	return nil, errors.New("invalid token")
}

func (s *service) CheckEmailValidity(email string) bool {
	start := time.Now()

	if len(email) > 254 || strings.Contains(email, "..") {
		log.Warn().
			Str("email", pkg.MaskEmail(email)).
			Dur("duration", time.Since(start)).
			Msg("  Invalid email: exceeds length or contains consecutive dots")
		return false
	}

	if !govalidator.IsEmail(email) {
		log.Warn().
			Str("email", pkg.MaskEmail(email)).
			Dur("duration", time.Since(start)).
			Msg("  Invalid email format")
		return false
	}

	if !govalidator.IsExistingEmail(email) {
		log.Warn().
			Str("email", pkg.MaskEmail(email)).
			Dur("duration", time.Since(start)).
			Msg("Email domain does not exist")
		return false
	}

	return true
}

func (s *service) ValidateUserExist(id uuid.UUID) (bool, error) {
	return s.userService.CheckUserExistsByID(id)
}
