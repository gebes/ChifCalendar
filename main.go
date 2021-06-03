package main

import (
	"./utils"
	"./webhook"
	"log"
)

func main() {

	utils.InitEnvironment()

	log.Println("Starting ChifCalendar")
	webhook.InitDiscord()
	webhook.InitCronJob()
}
