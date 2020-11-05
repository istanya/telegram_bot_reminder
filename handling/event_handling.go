package handling

import(
	"log"
	"context"
	"time"
	db "telegram_bot_reminder/db/sqlc"
	tgbotapi "gopkg.in/go-telegram-bot-api/telegram-bot-api.v5"
)
const timeSleep = time.Second*10

func EventsUpdate(store *db.Store, events []db.Event, state string){
	for _, event := range events{
		store.UpdateEventState(context.Background(), db.UpdateEventStateParams{
			ID: event.ID,
			State: state,
		})
	}
}

func EventHandling(server *Server) {
	for {
		events,err:= server.store.ListEventsByTimeAndState(context.Background(), db.ListEventsByTimeAndStateParams{
			DtReminder: time.Now().Add(timeSleep+3),
			State:      string(PENDING),
		})
		if err != nil {
			log.Printf("cannot get events: %s", err)
		}
		EventsUpdate(server.store, events, PROCESSED)

		for _, event := range events{
			go sendReminder(server, event)
		}

		time.Sleep(timeSleep)
	}
}

func sendReminder(server *Server, event db.Event){
	log.Printf("sendReminder[%s]: time now %s  ----  time reminder %s", event.UserName, time.Now(), event.DtReminder)
	time.Sleep(event.DtReminder.Sub(time.Now()))

	mgConfig := tgbotapi.NewMessage(event.UserID, event.Message)
	server.bot.Send(mgConfig)

	server.store.UpdateEventState(context.Background(), db.UpdateEventStateParams{
		ID: event.ID,
		State: SUCCESS,
	})
}
