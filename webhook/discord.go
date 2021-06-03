package webhook

import (
	"gebes.io/calendar/utils"
	"github.com/nickname32/discordhook"
	"log"
)

var WA *discordhook.WebhookAPI  = nil

func InitDiscord() {
	log.Println("Initializing Discord Webhook")
	keys := utils.GetDiscordWebhookKeys()
	wa, err := discordhook.NewWebhookAPI(keys.WebhookID, keys.WebhookToken, true, nil);

	if err != nil{
		panic(err)
	}
	WA = wa
	log.Println("Initialized Discord Webhook")
}