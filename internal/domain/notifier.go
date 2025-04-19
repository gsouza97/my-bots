package domain

type Notifier interface {
	SendMessage(message string) error
}
