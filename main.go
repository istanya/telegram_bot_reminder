package main

import (
	"database/sql"
	"log"
	"net/http"
	db "telegram_bot_reminder/db/sqlc"
	"telegram_bot_reminder/handling"
	"telegram_bot_reminder/util"

	_ "github.com/lib/pq"
	tgbotapi "gopkg.in/go-telegram-bot-api/telegram-bot-api.v5"
)


func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)

	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Fatal("cannot connect to bot:", err)
	}

	// bot.Debug = true
	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(config.WebhookURL))
	if err != nil {
		panic(err)
	}

	bot.SetMyCommands([]tgbotapi.BotCommand{
		{Command: "/set_reminder", Description: "set reminder"},
		{Command: "/close", Description: "close"},
	})

	updates := bot.ListenForWebhook("/")

	go http.ListenAndServe(":8080", nil)
	log.Println("start listen :8080")

	server := handling.NewServer(store, bot)

	go handling.EventHandling(server)

	// получаем все обновления из канала updates
	for update := range updates {
		go handling.MessageHandling(server, update)
	}
}
