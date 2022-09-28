package notifications

import (
	"context"

	"github.com/ArtemVoronov/indefinite-studies-notifications-service/internal/services"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/services/notifications"
	"google.golang.org/grpc"
)

type NotificationsServiceServer struct {
	notifications.UnimplementedNotificationsServiceServer
}

func RegisterServiceServer(s *grpc.Server) {
	notifications.RegisterNotificationsServiceServer(s, &NotificationsServiceServer{})
}

func (s *NotificationsServiceServer) SendEmail(ctx context.Context, in *notifications.SendEmailRequest) (*notifications.SendEmailReply, error) {
	err := services.Instance().Mail().SendEmail(in.GetSender(), in.GetRecepient(), in.GetBody())
	if err != nil {
		return nil, err
	}
	return &notifications.SendEmailReply{}, nil
}
