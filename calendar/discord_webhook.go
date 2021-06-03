package calendar

import (
	"fmt"
	"github.com/nickname32/discordhook"
	"log"

)




func initWebhook() {
	log.Println("Initializing Discord Webhook")
	keys := getWebhookKeys()
	wa, err := discordhook.NewWebhookAPI(keys.WebhookID, keys.WebhookToken, true, nil);

	if err != nil{
		panic(err)
	}


	msg, err := wa.Execute(nil, &discordhook.WebhookExecuteParams{
		Content: "Example text",
		Embeds: []*discordhook.Embed{
			{
				Title:       "Hi there",
				Description: "This is description",
			},
		},
	}, nil, "")
	if err != nil {
		panic(err)
	}

	fmt.Println(msg.ID)
}