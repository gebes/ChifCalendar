package webhook

import (
	"gebes.io/calendar/database"
	"github.com/robfig/cron"
	"log"
	"strconv"
	"time"
)



func InitCronJob() {

	log.Println("Initializing CronJob")
	broadcastToAllWebhooks()
	location, err := time.LoadLocation("Europe/Vienna")
	if err != nil {
		panic(err)
	}
	c := cron.NewWithLocation(location)
	c.AddFunc("0 0 "+strconv.Itoa(MessageTime)+" * * *", broadcastToAllWebhooks)
	c.Start()

	select {}
}

func broadcastToAllWebhooks(){

	webhooks, err := database.GetWebhooks()

	if err != nil {
		log.Println("Could not get webhooks from database", err)
		return
	}

	for _, webhook := range webhooks {
		calendar := NewCalendar(webhook)
		calendar.SendEventsMessage()
	}

}

