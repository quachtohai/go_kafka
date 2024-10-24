package main

import (
	"flag"
	"go_kafka/pkg/logger"
	"go_kafka/write_side/config"
	"log"

	kafkaConfigServer "go_kafka/write_side/server"
)

func main() {
	// envConfig := config.NewEnvConfig()

	// db := db.Init(envConfig, db.DBMigrator)

	// app := fiber.New(
	// 	fiber.Config{
	// 		AppName:      "GO_KAFKA",
	// 		ServerHeader: "Fiber",
	// 	})

	// eventRepository := repositories.NewEventRepository(db)
	// ticketRepository := repositories.NewTicketRepository(db)
	// authRepository := repositories.NewAuthRepository(db)

	// // Service
	// authService := services.NewAuthService(authRepository)

	// // Routing
	// server := app.Group("/api")
	// handlers.NewAuthHandler(server.Group("/auth"), authService)
	// privateRoutes := server.Use(middlewares.AuthProtected(db))
	// handlers.NewEventHandler(privateRoutes.Group("/event"), eventRepository,kafka.NewProducer())
	// handlers.NewTicketHandler(privateRoutes.Group("/ticket"), ticketRepository)

	flag.Parse()

	cfg, err := config.InitConfig()

	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("WriterService")

	s := kafkaConfigServer.NewServer(appLogger, cfg)

	appLogger.Fatal(s.Run())

	//fmt.Printf(s.Run().Error())

}
