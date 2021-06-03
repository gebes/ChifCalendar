package main

import (
	"gebes.io/calendar/utils"
	"gebes.io/calendar/webhook"
	"log"
)

func main() {
	utils.InitEnvironment()

	log.Println("Starting ChifCalendar")

	webhook.InitDiscord()
	webhook.InitCronJob()

}
