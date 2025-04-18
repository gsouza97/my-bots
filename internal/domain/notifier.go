package domain

type Notifier interface {
	SendMessage(chatID int64, message string) error
}
