package handling

import (
	"gopkg.in/go-telegram-bot-api/telegram-bot-api.v5"
	db "telegram_bot_reminder/db/sqlc"
)

type Server struct {
	store *db.Store
	bot   *tgbotapi.BotAPI
}

func NewServer(store *db.Store, bot *tgbotapi.BotAPI) *Server {
	server := &Server{store: store, bot: bot}
	return server
}
