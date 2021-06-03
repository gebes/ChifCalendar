package calendar

import (
	"github.com/andersfylling/snowflake"
	"os"
	"strconv"
)

type WebhookKeys struct {
	WebhookID snowflake.Snowflake
	WebhookToken string
}

func getWebhookKeys() *WebhookKeys {
	number, err := strconv.ParseUint(os.Getenv("WebhookID"), 10, 64)
	if err != nil {
		panic(err)
	}
	keys := WebhookKeys{WebhookID: snowflake.NewSnowflake(number), WebhookToken: os.Getenv("WebhookToken")}
	return &keys
}