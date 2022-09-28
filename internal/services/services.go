package services

import (
	"log"
	"sync"

	"github.com/ArtemVoronov/indefinite-studies-notifications-service/internal/services/notifications/mail"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/app"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/services/auth"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/utils"
)

type Services struct {
	auth *auth.AuthGRPCService
	mail *mail.EmailNotificationsService
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
	return &Services{
		auth: auth.CreateAuthGRPCService(utils.EnvVar("AUTH_SERVICE_GRPC_HOST")+":"+utils.EnvVar("AUTH_SERVICE_GRPC_PORT"), &authcreds),
		mail: mail.CreateEmailNotificationsService(utils.EnvVar("SMTP_SERVER_HOST") + ":" + utils.EnvVar("SMTP_SERVER_PORT")),
	}
}

func (s *Services) Shutdown() {
	s.auth.Shutdown()
	s.mail.Shutdown()
}

func (s *Services) Auth() *auth.AuthGRPCService {
	return s.auth
}

func (s *Services) Mail() *mail.EmailNotificationsService {
	return s.mail
}
