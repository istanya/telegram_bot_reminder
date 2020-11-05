package handling

import (
	"fmt"
	tgbotapi "gopkg.in/go-telegram-bot-api/telegram-bot-api.v5"
	"log"
	"strconv"
	"sync"
)

var yearsKeyboard = bildKeyboardByMap(Years)
var monthsKeyboard = bildKeyboardByMap(commonMonth)

func bildKeyboardByMap(keyboards map[string]int) tgbotapi.ReplyKeyboardMarkup {
	columns := 4
	keyboardButton := [][]tgbotapi.KeyboardButton{}
	keyboardButton = append(keyboardButton, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Выберите дату напоминания"),
	))

	i := 0
	keysRow := []tgbotapi.KeyboardButton{}
	for key, _ := range keyboards {
		i++
		keysRow = append(keysRow, tgbotapi.NewKeyboardButton(key))
		if ((i % columns) == 0) || i == len(keyboards) {
			keyboardButton = append(keyboardButton, tgbotapi.NewKeyboardButtonRow(keysRow...))
			keysRow = keysRow[:0]
		}
	}
	return tgbotapi.NewReplyKeyboard(keyboardButton...)
}

func bildKeyboardByDgit(num int) tgbotapi.ReplyKeyboardMarkup {
	columns := 7
	keyboardButton := [][]tgbotapi.KeyboardButton{}
	keyboardButton = append(keyboardButton, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Выберите дату напоминания"),
	))
	keysRow := []tgbotapi.KeyboardButton{}
	for i := 1; i <= num; i++ {
		keysRow = append(keysRow, tgbotapi.NewKeyboardButton(strconv.Itoa(i)))
		if ((i % columns) == 0) || i == num {
			keyboardButton = append(keyboardButton, tgbotapi.NewKeyboardButtonRow(keysRow...))
			keysRow = keysRow[:0]
		}
	}
	return tgbotapi.NewReplyKeyboard(keyboardButton...)
}

func DeleteMessage(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
	deleteMessageConfig := tgbotapi.NewDeleteMessage(chatID, messageID)
	bot.DeleteMessage(deleteMessageConfig)
}

func UpdateCurentBotMessage(bot *tgbotapi.BotAPI, chatID int64, messageID int, newMessage string, keyboard tgbotapi.ReplyKeyboardMarkup) tgbotapi.Message {
	mgConfig := tgbotapi.NewMessage(chatID, newMessage)
	mgConfig.ReplyMarkup = keyboard
	messageTo, _ := bot.Send(mgConfig)

	deleteMessageConfig := tgbotapi.NewDeleteMessage(chatID, messageID)
	bot.DeleteMessage(deleteMessageConfig)
	return messageTo
}

func getDaysMonth(year string, month string) int {
	var num int
	if Years[year] > 365 {
		num = leapMonth[month]
	} else {
		num = commonMonth[month]
	}
	return num
}

func MessageHandling(server *Server, update tgbotapi.Update) {
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	command := update.Message.Command()

	ChatID := update.Message.Chat.ID
	if command != "start" {
		DeleteMessage(server.bot, ChatID, update.Message.MessageID)
	} else {
		return
	}

	mu := &sync.Mutex{}
	mu.Lock()
	user, okUser := Users[ChatID]
	mu.Unlock()
	if !okUser {
		command = "set_reminder"
	}

	
	// комманда - сообщение, начинающееся с "/"
	switch command {
	case "set_reminder":
		messageText := "Пока не выбрана дата напоминания"
		messageTo := UpdateCurentBotMessage(server.bot, ChatID, 0, messageText, yearsKeyboard)
		mu.Lock()
		Users[ChatID] = User{
			ID:           update.Message.From.ID,
			UserName:     update.Message.From.UserName,
			BotMessageId: messageTo.MessageID,
		}
		mu.Unlock()
		return
	case "close":
		key := tgbotapi.NewMessage(ChatID, update.Message.Command())
		key.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		server.bot.Send(key)
		return
	}

	//Анализируем Меню
	text := update.Message.Text
	if _, ok := Years[text]; ok && user.Year == "" {

		messageText := fmt.Sprintf("Дата напоминания: %sг", text)
		messageTo := UpdateCurentBotMessage(server.bot, ChatID, user.BotMessageId, messageText, monthsKeyboard)

		user.BotMessageId = messageTo.MessageID
		user.Year = text
		mu.Lock()
		Users[ChatID] = user
		mu.Unlock()
	} else if _, ok := commonMonth[text]; ok && user.Year != "" && user.Month == "" {
		numDaysMonth := getDaysMonth(user.Year, text)
		daysMonthKeyboard := bildKeyboardByDgit(numDaysMonth)

		messageText := fmt.Sprintf("Дата напоминания: %sг %s", user.Year, text)
		messageTo := UpdateCurentBotMessage(server.bot, ChatID, user.BotMessageId, messageText, daysMonthKeyboard)

		user.BotMessageId = messageTo.MessageID
		user.Month = text
		mu.Lock()
		Users[ChatID] = user
		mu.Unlock()
	} else if user.Day == "" && user.Year != "" && user.Month != "" {
		numDaysMonth := getDaysMonth(user.Year, user.Month)
		if !validDigitDiapason(text, 1, numDaysMonth) {
			return
		}
		daysMonthKeyboard := bildKeyboardByDgit(24)

		messageText := fmt.Sprintf("Дата напоминания: %sг %s %s", user.Year, user.Month, text)
		messageTo := UpdateCurentBotMessage(server.bot, ChatID, user.BotMessageId, messageText, daysMonthKeyboard)

		user.BotMessageId = messageTo.MessageID
		user.Day = text
		mu.Lock()
		Users[ChatID] = user
		mu.Unlock()
	} else if user.Hour == "" && user.Day != "" && user.Year != "" && user.Month != "" {
		if !validDigitDiapason(text, 1, 24) {
			return
		}
		daysMonthKeyboard := bildKeyboardByDgit(60)

		messageText := fmt.Sprintf("Дата напоминания: %sг %s %s \nВремя: %s", user.Year, user.Month, user.Day, text)
		messageTo := UpdateCurentBotMessage(server.bot, ChatID, user.BotMessageId, messageText, daysMonthKeyboard)

		user.BotMessageId = messageTo.MessageID
		user.Hour = text
		mu.Lock()
		Users[ChatID] = user
		mu.Unlock()
	} else if user.Minute == "" && user.Hour != "" && user.Day != "" && user.Year != "" && user.Month != "" {
		if !validDigitDiapason(text, 1, 60) {
			return
		}
		daysMonthKeyboard := tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Введите сообщение напоминание"),
		))

		messageText := fmt.Sprintf("Дата напоминания: %sг %s %s \nВремя: %s:%s", user.Year, user.Month, user.Day, user.Hour, text)
		messageTo := UpdateCurentBotMessage(server.bot, ChatID, user.BotMessageId, messageText, daysMonthKeyboard)

		user.BotMessageId = messageTo.MessageID
		user.Minute = text
		mu.Lock()
		Users[ChatID] = user
		mu.Unlock()
	} else if user.MessageReminder == "" && user.Minute != "" && user.Hour != "" && user.Day != "" && user.Year != "" && user.Month != "" {
		daysMonthKeyboard := tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/close"),
		))

		messageText := fmt.Sprintf("Дата напоминания: %sг %s %s \nВремя: %s:%s\nСообщение: %s", user.Year, user.Month, user.Day, user.Hour, user.Minute, text)
		messageTo := UpdateCurentBotMessage(server.bot, ChatID, user.BotMessageId, messageText, daysMonthKeyboard)

		user.BotMessageId = messageTo.MessageID
		user.MessageReminder = text
		mu.Lock()
		Users[ChatID] = user
		mu.Unlock()
		err := user.SaveBD(server.store)
		if err != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, err)
		}
	}
}
