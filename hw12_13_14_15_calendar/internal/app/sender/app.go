package sender

import (
	"context"
	"fmt"
	"os"

	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/interfaces"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/rmq"
	memorystorage "github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
	sql "github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/storage/sql"
)

type AppSender struct {
	logger  interfaces.Logger
	storage interfaces.EventStorage
}

func NewSenderApp(ctx context.Context, conf *config.SenderConfig) (*AppSender, error) {
	app := &AppSender{}

	// Инициализация логгера.
	logger := logger.New(conf.Logger.Level, os.Stdout)
	app.logger = logger

	// Инициализация хранилища данных.
	var eventStorage interfaces.EventStorage
	if conf.Storage.Type == "postgres" {
		psqlStorage := new(sql.Storage)
		if err := psqlStorage.Connect(
			ctx,
			conf.Database.Port,
			conf.Database.Host,
			conf.Database.Username,
			conf.Database.Password,
			conf.Database.Dbname,
		); err != nil {
			return nil, fmt.Errorf("cannot connect to PostgreSQL: %w", err)
		}
		err := psqlStorage.Migrate(ctx, conf.Storage.Migration)
		if err != nil {
			return nil, fmt.Errorf("migration did not work out: %w", err)
		}
		eventStorage = psqlStorage
	} else {
		eventStorage = memorystorage.New()
	}
	app.storage = eventStorage

	// Инициализация eventsConsMq.
	eventsConsMq, err := rmq.New(
		conf.RMQ.URI,
		conf.Queues.Events.ExchangeName,
		conf.Queues.Events.ExchangeType,
		conf.Queues.Events.QueueName,
		conf.Queues.Events.BindingKey,
		conf.RMQ.ReConnect.MaxElapsedTime,
		conf.RMQ.ReConnect.InitialInterval,
		conf.RMQ.ReConnect.Multiplier,
		conf.RMQ.ReConnect.MaxInterval,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize eventsConsMq RMQ for sender: %w", err)
	}

	if err := eventsConsMq.Init(ctx); err != nil {
		logger.Error("RMQ initialization failed: %v", err)
	}

	msgs, err := eventsConsMq.Consume(conf.Consumer.ConsumerTag)
	if err != nil {
		logger.Error("Failed to register a consumer: %v", err)
		return nil, nil
	}

	logger.Info("Sender is running...")

	for msg := range msgs {
		logger.Info("Received a message: %s", string(msg.Body))
		msg.Ack(false)
	}

	return app, nil
}
