package main

import (
	"gebes.io/calendar/database"
	"gebes.io/calendar/utils"
	"gebes.io/calendar/webhook"
	"log"
)

func main() {
	utils.InitEnvironment()

	log.Println("Starting ChifCalendar")

	database.Connect()

	webhook.InitCronJob()

}
