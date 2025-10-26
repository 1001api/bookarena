package routes

import (
	"1001api/bookarena/internal/controllers"
	"1001api/bookarena/internal/modules/auth"
	"1001api/bookarena/internal/modules/users"
	"1001api/bookarena/pkg"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func Routing(r fiber.Router, db *pgxpool.Pool) {
	var encKey = viper.GetString("ENC_KEY")

	if encKey == "" {
		log.Fatal().Msg("ENC_KEY is not set")
	}

	validator := validator.New()

	userRepo := users.NewUserRepository(db)

	userService := users.NewUserService(userRepo, encKey)
	authService := auth.NewAuthService(userService, encKey)

	userController := controllers.NewUserController(userService, validator)
	authController := controllers.NewAuthController(authService, validator)

	// Initialize root user
	if err := userService.InitializeRootUser(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize root user")
	}

	// API versioning
	versioning := r.Group("/api/v1")

	auth := versioning.Group("/auth")
	{
		auth.Post("/login", authController.Login)
		auth.Post("/register", authController.RegisterUser)
		auth.Post("/refresh", authController.Refresh)
	}

	userRoutes := versioning.Group("/users", JWTMiddleware(authService))
	{
		userRoutes.Get("/me", userController.GetCurrentUser)
		userRoutes.Patch("/me", userController.UpdateCurrentUser)
		userRoutes.Get("/list", RoleMiddleware(string(pkg.RoleAdmin), string(pkg.RoleSuperAdmin)), userController.GetMasterUser)
		userRoutes.Post("/create", RoleMiddleware(string(pkg.RoleAdmin), string(pkg.RoleSuperAdmin)), authController.RegisterAdmin)
		userRoutes.Get("/search", RoleMiddleware(string(pkg.RoleAdmin), string(pkg.RoleSuperAdmin)), userController.SearchUser)
		userRoutes.Get("/:id", RoleMiddleware(string(pkg.RoleAdmin), string(pkg.RoleSuperAdmin)), userController.GetUserByID)
		userRoutes.Patch("/:id", RoleMiddleware(string(pkg.RoleAdmin), string(pkg.RoleSuperAdmin)), userController.UpdateUser)
		userRoutes.Delete("/:id", RoleMiddleware(string(pkg.RoleAdmin), string(pkg.RoleSuperAdmin)), userController.DeleteUser)
	}
}

func JWTMiddleware(authService auth.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		accessToken := c.Cookies(pkg.AccessTokenName)

		if accessToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Validate token
		claims, err := authService.ValidateToken(accessToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		isExist, err := authService.ValidateUserExist(claims.UserID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to validate user exist")
		}

		if !isExist {
			log.Warn().Str("user_id", claims.UserID.String()).Msg("User not found but token is valid")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role")
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User role not found in context",
			})
		}

		role := cast.ToString(userRole)
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}
}
