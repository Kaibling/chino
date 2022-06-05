package services

type notify interface {
	Send(message string) error
}

type NotificationService struct {
	n notify
}

func NewNotificationService(n notify) *NotificationService {
	return &NotificationService{n}
}

func (s *NotificationService) Send(message string) error {
	return s.n.Send(message)
}
