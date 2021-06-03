package utils

import (
	"github.com/andersfylling/snowflake"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func InitEnvironment(){
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}


type DiscordWebhookKeys struct {
	WebhookID snowflake.Snowflake
	WebhookToken string
}

func GetDiscordWebhookKeys() *DiscordWebhookKeys {
	number, err := strconv.ParseUint(os.Getenv("DiscordWebhookID"), 10, 64)
	if err != nil {
		panic(err)
	}
	keys := DiscordWebhookKeys{WebhookID: snowflake.NewSnowflake(number), WebhookToken: os.Getenv("DiscordWebhookToken")}
	return &keys
}

type CalendarKeys struct {
	CalenderURL string
}

func GetCalendarKeys() *CalendarKeys {
	notion := CalendarKeys{CalenderURL: os.Getenv("CalendarUrl")}
	return &notion
}