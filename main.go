package main

import (
	"github.com/Gebes/ChifCalendar/utils"
	"github.com/Gebes/ChifCalendar/webhook"
	"log"
)

func main() {

	utils.InitEnvironment()

	log.Println("Starting ChifCalendar")
	webhook.InitDiscord()
	webhook.InitCronJob()
}
