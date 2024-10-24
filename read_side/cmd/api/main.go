package main

import (
	"flag"
	"fmt"
	"go_kafka/pkg/logger"
	"go_kafka/read_side/config"
	"go_kafka/read_side/db"
	"go_kafka/read_side/handlers"
	"go_kafka/read_side/middlewares"
	"go_kafka/read_side/repositories"
	"go_kafka/read_side/services"

	kafkaConfigServer "go_kafka/read_side/server"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	envConfig := config.NewEnvConfig()

	db := db.Init(envConfig, db.DBMigrator)

	app := fiber.New(
		fiber.Config{
			AppName:      "GO_KAFKA",
			ServerHeader: "Fiber",
		})

	eventRepository := repositories.NewEventRepository(db)
	ticketRepository := repositories.NewTicketRepository(db)
	authRepository := repositories.NewAuthRepository(db)

	// Service
	authService := services.NewAuthService(authRepository)

	// Routing
	server := app.Group("/api")
	handlers.NewAuthHandler(server.Group("/auth"), authService)
	privateRoutes := server.Use(middlewares.AuthProtected(db))
	handlers.NewEventHandler(privateRoutes.Group("/event"), eventRepository)
	handlers.NewTicketHandler(privateRoutes.Group("/ticket"), ticketRepository)
	flag.Parse()

	cfg, err := config.InitConfig()

	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("ReaderService")

	s := kafkaConfigServer.NewServer(appLogger, cfg)

	go func() {
		appLogger.Fatal(s.Run())
	}()
	//fmt.Printf(s.Run().Error())

	app.Listen(fmt.Sprintf(":" + envConfig.ServerPort))
}
