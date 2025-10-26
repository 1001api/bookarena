package main

import (
	"1001api/bookarena/internal/database"
	"1001api/bookarena/internal/routes"
	"os"
	"runtime"
	"time"
	_ "time/tzdata"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var (
	db *pgxpool.Pool
)

func init() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Msg("No .env file found, continuing with environment variables")
	}
	viper.AutomaticEnv()
}

func main() {
	// Jakarta Timezone
	jakarta, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load Jakarta location")
	}

	// Setup Zerolog
	multiWriter := zerolog.MultiLevelWriter(
		zerolog.ConsoleWriter{
			Out:          os.Stdout,
			TimeFormat:   time.RFC3339,
			TimeLocation: jakarta,
		},
	)
	log.Logger = zerolog.New(multiWriter).With().Timestamp().Logger()

	// connect to database
	db = database.ConnectPG()

	// connect to redis
	// redisCon := database.NewRedisCon()

	r := fiber.New()

	// Recover middleware
	r.Use(recover.New())

	// Global Limit
	r.Use(limiter.New(limiter.Config{
		Max:        60,               // 60 requests
		Expiration: 60 * time.Second, // 1 minute
		KeyGenerator: func(c *fiber.Ctx) string {
			ip := c.IP()
			if fwd := c.Get("X-Forwarded-For"); fwd != "" {
				ip = fwd
			}
			return ip
		},
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	// Request Logging
	r.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger:          &log.Logger,
		FieldsSnakeCase: true,
	}))

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     viper.GetString("CLIENT_DOMAIN"),
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Encrypt Cookie
	r.Use(encryptcookie.New(encryptcookie.Config{
		Key: viper.GetString("COOKIE_ENC_KEY"),
	}))

	// ROUTING
	routes.Routing(r, db)

	r.Get("/health", monitor.New(monitor.Config{
		Title: "Bookarena",
	}))

	port := viper.GetString("PORT")
	if port == "" {
		port = "8181"
	}

	log.Info().Str("port", port).Msg("Server successfully started")
	log.Info().Str("go_version", runtime.Version()).Msg("Go version")

	if err := r.Listen(":" + port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
