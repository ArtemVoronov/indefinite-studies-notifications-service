package services

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/ArtemVoronov/indefinite-studies-notifications-service/internal/services/notifications/mail"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/app"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/services/auth"
	kafkaService "github.com/ArtemVoronov/indefinite-studies-utils/pkg/services/kafka"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/services/watcher"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/utils"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Services struct {
	auth    *auth.AuthGRPCService
	mail    *mail.EmailNotificationsService
	watcher *watcher.WatcherService
}

var once sync.Once
var instance *Services

func Instance() *Services {
	once.Do(func() {
		if instance == nil {
			instance = createServices()
		}
	})
	return instance
}

func createServices() *Services {
	authcreds, err := app.LoadTLSCredentialsForClient(utils.EnvVar("AUTH_SERVICE_CLIENT_TLS_CERT_PATH"))
	if err != nil {
		log.Fatalf("unable to load TLS credentials")
	}

	mailService := mail.CreateEmailNotificationsService(
		utils.EnvVar("SMTP_SERVER_HOST")+":"+utils.EnvVar("SMTP_SERVER_PORT"),
		utils.EnvVarDurationDefault("SMPT_SERVER_CONNECT_TIMEOUT_IN_SECONDS", time.Second, 10*time.Second),
	)

	// TODO: use configured logrus
	watcherService := watcher.CreateWatcherService(
		utils.EnvVar("KAFKA_HOST")+":"+utils.EnvVar("KAFKA_PORT"),
		utils.EnvVar("KAFKA_GROUP_ID"),
		kafkaService.EVENT_TYPE_SEND_EMAIL,
		30_000,
		func(e *kafka.Message) {
			var dto kafkaService.SendEmailEvent
			err := json.Unmarshal(e.Value, &dto)
			if err != nil {
				log.Printf("Error during parsing SEND_MAIL event message: %s\n", err)
				return
			}

			// TODO: clean logging
			log.Printf("Message on %s: %s\n", e.TopicPartition, dto)

			err = mailService.SendEmail(dto.Sender, dto.Recepient, dto.Subject, dto.Body)
			if err != nil {
				log.Printf("Error during sending email: %s\n", err)
			}
		},

		func(e error) {
			log.Printf("Error during watching events: %s\n", e)
		},
	)

	return &Services{
		auth:    auth.CreateAuthGRPCService(utils.EnvVar("AUTH_SERVICE_GRPC_HOST")+":"+utils.EnvVar("AUTH_SERVICE_GRPC_PORT"), &authcreds),
		mail:    mailService,
		watcher: watcherService,
	}
}

func (s *Services) Shutdown() {
	s.auth.Shutdown()
	s.mail.Shutdown()
	s.watcher.Shutdown()
}

func (s *Services) Auth() *auth.AuthGRPCService {
	return s.auth
}

func (s *Services) Mail() *mail.EmailNotificationsService {
	return s.mail
}
