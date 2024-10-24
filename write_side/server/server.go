package server

import (
	"context"
	"fmt"
	kafkaClient "go_kafka/pkg/kafka"
	"go_kafka/pkg/logger"
	"go_kafka/write_side/config"
	"go_kafka/write_side/db"
	kafkaConsumer "go_kafka/write_side/delivery/kafka"
	"go_kafka/write_side/handlers"
	"go_kafka/write_side/middlewares"
	"go_kafka/write_side/repositories"
	"go_kafka/write_side/services"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type server struct {
	log       logger.Logger
	cfg       *config.KafkaConfig
	kafkaConn *kafka.Conn
}

func NewServer(log logger.Logger, cfg *config.KafkaConfig) *server {

	return &server{log: log, cfg: cfg}
}
func (s *server) Run() error {
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

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	kafkaProducer := kafkaClient.NewProducer(s.log, s.cfg.Kafka.Brokers)
	defer kafkaProducer.Close()

	eventMessageProcessor := kafkaConsumer.NewEventMessageProcessor(s.log, s.cfg)
	s.log.Info("Starting Writer Kafka consumers")
	server := app.Group("/api")
	handlers.NewAuthHandler(server.Group("/auth"), authService)
	privateRoutes := server.Use(middlewares.AuthProtected(db))
	handlers.NewEventHandler(privateRoutes.Group("/event"), eventRepository, kafkaProducer)
	handlers.NewTicketHandler(privateRoutes.Group("/ticket"), ticketRepository)
	cg := kafkaClient.NewConsumerGroup(s.cfg.Kafka.Brokers, s.cfg.Kafka.GroupID, s.log)
	go cg.ConsumeTopic(ctx, s.getConsumerGroupTopics(), kafkaConsumer.PoolSize, eventMessageProcessor.ProcessMessages)
	if err := s.connectKafkaBrokers(ctx); err != nil {
		return errors.Wrap(err, "s.connectKafkaBrokers")

	}
	defer s.kafkaConn.Close() // nolint: errcheck

	if s.cfg.Kafka.InitTopics {
		s.initKafkaTopics(ctx)
	}
	app.Listen(fmt.Sprintf(":" + envConfig.ServerPort))
	<-ctx.Done()
	return nil
}
