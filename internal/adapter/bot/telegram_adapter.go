package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type TelegramAdapter struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramAdapter(token string) (*TelegramAdapter, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &TelegramAdapter{
		bot: bot,
	}, nil
}

func (ta *TelegramAdapter) HandleUpdates() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := ta.bot.GetUpdatesChan(u)
	if err != nil {
		return nil, err
	}

	return updates, nil
}

func (ta *TelegramAdapter) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := ta.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
