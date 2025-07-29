package main

import (
	"fmt"
	"github.com/Rickykn/rental-car/config"
	"github.com/Rickykn/rental-car/db"
	"github.com/Rickykn/rental-car/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"strconv"
)

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("grpc: main failed to load and parse config: %s", err)
		return
	}

	dbCon := db.Config{
		AppInfo:     cfg.AppInfo.Name,
		Username:    cfg.PostgreSQL.Username,
		Password:    cfg.PostgreSQL.Password,
		Database:    cfg.PostgreSQL.Database,
		Host:        cfg.PostgreSQL.Host,
		SSLMode:     cfg.PostgreSQL.SSLMode,
		Port:        cfg.PostgreSQL.Port,
		ConnMaxOpen: cfg.PostgreSQL.MaxOpenConns,
		ConnMaxIdle: cfg.PostgreSQL.MaxIdleConns,
		Logging:     cfg.PostgreSQL.Logging,
	}

	newDB, err := db.NewDB(dbCon)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${ip} - ${method} ${path} - ${status} - ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	routes.SetupRoutes(app, newDB)

	// Start server
	err = app.Listen(fmt.Sprintf(":%s", strconv.Itoa(cfg.AppInfo.Port)))
	if err != nil {
		panic(err)
	}
}
