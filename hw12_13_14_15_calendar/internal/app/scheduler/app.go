package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/interfaces"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/rmq"
	"github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
	sql "github.com/dianapovarnitsina/hw-test/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/streadway/amqp"
)

type AppScheduler struct {
	logger  interfaces.Logger
	storage interfaces.EventStorage
}

//nolint:gocognit
func NewSchedulerApp(ctx context.Context, conf *config.SchedulerConfig) (*AppScheduler, error) {
	app := &AppScheduler{}

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

	eventsProdMq, err := rmq.New(
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
		return nil, fmt.Errorf("can't initialize rmq for events-producer: %w", err)
	}

	if err := eventsProdMq.Init(ctx); err != nil {
		logger.Error("RMQ failed: %v", err)
	}

	// Горутина для отправки уведомлений
	go func() {
		// Используем таймер для запуска каждую минуту.
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return // Если контекст отменен, завершаем горутину.
			case <-ticker.C: // Если таймер сработал, то выполняем нужные действия.
				// Выбираем события для отправки уведомлений из БД.
				err := deleteOldEvents(ctx, app)
				if err != nil {
					logger.Error("Failed select events for notifications: %v", err)
					continue
				}

				events, err := selectEventsForNotifications(ctx, app)
				if err != nil {
					logger.Error("Failed select events for notifications: %v", err)
					continue
				}

				for _, event := range events {
					// Создаем уведомление для события.
					notification := createNotification(event)

					// Сериализуем уведомление в JSON.
					notificationJSON, err := serializeNotification(notification)
					if err != nil {
						logger.Error("Failed to serialize notification: %v", err)
						continue
					}

					// Отправляем уведомление в RabbitMQ.
					err = publishNotificationToRMQ(notificationJSON, eventsProdMq)
					if err != nil {
						logger.Error("Failed to publish notification to RMQ:  %v", err)
						continue
					}
					fmt.Println("отправили уведомление", time.Now().Format("2006-01-02 15:04"))
				}
			}
		}
	}()

	return app, nil
}

func selectEventsForNotifications(ctx context.Context, app *AppScheduler) ([]*storage.Event, error) {
	events, err := app.storage.SelectEventsForNotifications(ctx)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func deleteOldEvents(ctx context.Context, app *AppScheduler) error {
	return app.storage.DeleteOldEvents(ctx)

}

func createNotification(event *storage.Event) storage.Notification {
	notification := storage.Notification{
		EventID:  event.ID,
		Title:    event.Title,
		DateTime: event.DateTime,
		UserID:   event.UserID,
	}
	return notification
}

func serializeNotification(notification storage.Notification) ([]byte, error) {
	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return nil, err
	}
	return notificationJSON, nil
}

func publishNotificationToRMQ(notificationJSON []byte, rmq *rmq.Rmq) error {
	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        notificationJSON,
	}

	return rmq.Publish(msg)
}
