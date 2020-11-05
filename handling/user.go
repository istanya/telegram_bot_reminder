package handling

import (
	"context"
	"strconv"
	db "telegram_bot_reminder/db/sqlc"
	"time"
)

type StateEvent string

const (
	PENDING   StateEvent = "PENDING"
	PROCESSED            = "PROCESSED"
	SUCCESS              = "SUCCESS"
	FAILURE              = "FAILURE"
)

type User struct {
	ID              int    `json:"id"`
	UserName        string `json:"username"` // optional
	Year            string
	Month           string
	Day             string
	Hour            string
	Minute          string
	MessageReminder string
	BotMessageId    int
}

func (user User) SaveBD(store *db.Store) error {

	year, _ := strconv.Atoi(user.Year)
	day, _ := strconv.Atoi(user.Day)
	hour, _ := strconv.Atoi(user.Hour)
	minute, _ := strconv.Atoi(user.Minute)
	loc, _:= time.LoadLocation("Europe/Moscow")
	userBD := db.CreateEventParams{
		UserID:       int64(user.ID),
		UserName:     user.UserName,
		DtReminder:   time.Date(year, time.Month(mapMonthNum[user.Month]), day, hour, minute, 0, 0, loc),
		BotMessageID: int32(user.BotMessageId),
		Message:      user.MessageReminder,
		State:        string(PENDING),
		DtCreated:    time.Now(),
	}
	_, err := store.CreateEvent(context.Background(), userBD)
	if err != nil {
		return err
	}
	return nil
}

var Users = make(map[int64]User, 0)
