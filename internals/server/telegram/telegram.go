package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func NewTelegramBot(token string, chatID int64) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %v", err)
	}

	return &TelegramBot{
		bot:    bot,
		chatID: chatID,
	}, nil
}
func (t *TelegramBot) SendContactInfo(name, contactType, address, message string, option string, ip string) error {
	text := fmt.Sprintf(
		"\n\nИмя: %s\nТип контакта: %s\nАдрес: %s\nСообщение: %s\nВыбранный вариант: %s\nIP: %s",
		name,
    contactType,
		address,
		message,
    option,
    ip,
	)

	msg := tgbotapi.NewMessage(t.chatID, text)
	_, err := t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}
