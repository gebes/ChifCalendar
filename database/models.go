package database

type Message struct {
	id        int
	webhookId int
	hashValue string
}

type Webhook struct {
	Id int
	DiscordId    uint64
	DiscordToken string
	CalendarPrivateUrl string
	CalendarPublicUrl string
	Nickname    string
	PrimaryColor int
}

