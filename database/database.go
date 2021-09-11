package database

import (
	"database/sql"
	"gebes.io/calendar/utils"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var (
	Database *sql.DB
)

func Connect() {
	log.Println("Initializing database connection")
	db, err := sql.Open("mysql", utils.GetDatabaseUrl())
	if err != nil {
		panic(err)
	}
	Database = db
}

func LastMessageHash(webhookId int) *string {

	results, err := Database.Query("SELECT * FROM message WHERE webhookId = ? ORDER BY message.id DESC LIMIT 1;", webhookId)
	defer results.Close()

	if err != nil {
		panic(err)
	}

	var lastMessage Message

	results.Next()
	err = results.Scan(&lastMessage.id, &lastMessage.webhookId, &lastMessage.hashValue)
	if err != nil {
		return nil
	}

	return &lastMessage.hashValue
}

func SaveHash(webhookId int, hash string) {

	_, err := Database.Query("INSERT INTO message(message.webhookId, message.hashValue) VALUES(?, ?);", webhookId, hash)

	if err != nil {
		panic(err)
	}

}

func GetWebhooks() ([]Webhook, error) {
	result, err := Database.Query("SELECT * FROM webhook")
	defer result.Close()
	if err != nil {
		return nil, err
	}

	var webhooks []Webhook
	var webhook Webhook

	for result.Next() {
		err = result.Scan(&webhook.Id, &webhook.DiscordId, &webhook.DiscordToken, &webhook.CalendarPrivateUrl, &webhook.CalendarPublicUrl, &webhook.Nickname, &webhook.PrimaryColor)
		if err != nil {
			return nil, err
		}
		webhooks = append(webhooks, webhook)
	}

	return webhooks, nil
}
