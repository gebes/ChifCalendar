package webhook

import (
	"fmt"
	"gebes.io/calendar/database"
	"gebes.io/calendar/utils"
	"github.com/andersfylling/snowflake"
	"github.com/apognu/gocal"
	"github.com/mvdan/xurls"
	"github.com/nickname32/discordhook"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

type Category struct {
	title string
	emoji string
}

const (
	HomeworkEmoji = "📘"
	ExamEmoji     = "📝"
	OtherEmoji    = "👻"

	MessageTime = 17
)

var (
	Categories = []Category{
		{title: "Hausübungen", emoji: HomeworkEmoji},
		{title: "Prüfungen", emoji: ExamEmoji},
		{title: "Anderes", emoji: OtherEmoji},
	}
	ImportantKeywords = []string{"SA", "Schularbeit", "Test", "Revision"}
	shortDayNames     = []string{
		"So",
		"Mo",
		"Di",
		"Mi",
		"Do",
		"Fr",
		"Sa",
	}
)

type Calendar struct {
	database.Webhook
	discordWebhook *discordhook.WebhookAPI
}

func NewCalendar(webhook database.Webhook) *Calendar {
	return &Calendar{Webhook: webhook}
}

func (calendar *Calendar) BuildWebhook() (*discordhook.WebhookAPI, error) {

	if calendar.discordWebhook != nil {
		return calendar.discordWebhook, nil
	}

	wa, err := discordhook.NewWebhookAPI(snowflake.Snowflake(calendar.DiscordId), calendar.DiscordToken, true, nil)

	if err != nil {
		return nil, err
	}
	calendar.discordWebhook = wa
	return wa, nil
}

func (calendar *Calendar) SendEventsMessage() {

	tomorrowEmbed, tomorrowHash := calendar.getEmbedFor("Plan für morgen", calendar.fetchTomorrowEvents(), false)
	nextWeekEmbed, nextWeekHash := calendar.getEmbedFor("Demnächst", calendar.fetchNextWeekEvents(), true)

	hash := tomorrowHash + nextWeekHash
	lastHash := database.LastMessageHash(calendar.Id)

	if lastHash != nil && hash == *lastHash {
		log.Println("Last hash matches with current hash for", calendar.Nickname)
		return
	}

	log.Println("Got new hash for", calendar.Nickname, hash)

	_, err := calendar.BuildWebhook()

	if err != nil {
		log.Println("Could not build a webhook connection for", calendar.Nickname, err)
		return
	}

	err = calendar.SendEmbed(tomorrowEmbed)
	if err != nil {
		log.Println("Could not send tomorrowEmbed to", calendar.Nickname, err, tomorrowEmbed)
	} else if tomorrowEmbed != nil{
		log.Println("Sent tomorrowEmbed to", calendar.Nickname)
	}

	err = calendar.SendEmbed(nextWeekEmbed)
	if err != nil {
		log.Println("Could not send nextWeekEmbed to", calendar.Nickname, err, nextWeekEmbed)
	} else if nextWeekEmbed != nil {
		log.Println("Sent nextWeekEmbed to", calendar.Nickname)
	}

	log.Println("Saved hash for", calendar.Nickname)
	database.SaveHash(calendar.Id, hash)

}

func (calendar *Calendar) SendEmbed(embed *discordhook.Embed) error {
	if embed == nil {
		return nil
	}
	_, err := calendar.discordWebhook.Execute(nil, &discordhook.WebhookExecuteParams{
		Embeds: []*discordhook.Embed{
			embed,
		},
	}, nil, "")

	return err
}
func (calendar *Calendar) getEmbedFor(title string, events *[]gocal.Event, showDates bool) (*discordhook.Embed, string) {
	messages := make([]string, len(Categories))
	important := false

	var embedFields []*discordhook.EmbedField
	for i, category := range Categories {

		for _, event := range *events {
			if strings.Index(event.Summary, category.emoji) != 0 {
				continue
			}

			addRow(&category, &messages[i], showDates, &event)
		}

		if !important && containsImportantKeyword(&messages[i]) {
			important = true
		}

		if len(messages[i]) == 0 {
			continue
		}

		embedFields = append(embedFields, &discordhook.EmbedField{
			Name:   category.emoji + " " + category.title,
			Value:  messages[i],
			Inline: false,
		})

	}

	if len(embedFields) == 0 {
		return nil, utils.Hash("")
	}

	embed := &discordhook.Embed{
		Title:  title,
		Fields: embedFields,
		URL:    calendar.CalendarPublicUrl,
	}

	if important {
		embed.Color = calendar.PrimaryColor
	}

	return embed, utils.HashList(messages)

}

func addRow(category *Category, message *string, showDate bool, eventToAppend *gocal.Event) {
	*message += "• "
	if showDate {
		year, month, day := eventToAppend.Start.Date()
		dayOfWeek := eventToAppend.Start.Weekday()
		*message += fmt.Sprintf("**%s. %02d.%02d.%04d** ", shortDayNames[dayOfWeek], day, month, year)
	}
	summary := strings.TrimPrefix(eventToAppend.Summary, category.emoji)
	if len(eventToAppend.Description) != 0 && strings.Contains(eventToAppend.Description, "https://discord.com/channels") {
		url := xurls.Relaxed.FindString(eventToAppend.Description)
		*message += "[" + summary + "](" + url + ")"
	} else {
		*message += summary
	}

	*message += "\n"
}

func (calendar *Calendar) fetchTomorrowEvents() *[]gocal.Event {
	start, end := tomorrow().Add(8*time.Hour), tomorrow().Add(8*time.Hour)
	return calendar.fetchEventsInPeriod(&start, &end)
}
func (calendar *Calendar) fetchNextWeekEvents() *[]gocal.Event {
	start, end := tomorrow().Add(48*time.Hour), tomorrow().Add(8*24*time.Hour)
	return calendar.fetchEventsInPeriod(&start, &end)
}

func (calendar *Calendar) fetchEventsInPeriod(start *time.Time, end *time.Time) *[]gocal.Event {
	resp, err := http.Get(calendar.CalendarPrivateUrl)
	if err != nil {
		panic(err)
	}

	c := gocal.NewParser(resp.Body)

	c.Start, c.End = start, end
	err = c.Parse()
	if err != nil {
		panic(err)
	}

	events := c.Events
	sort.Slice(events[:], func(i, j int) bool {
		return events[i].Start.Unix() < events[j].Start.Unix()
	})
	return &events
}

func tomorrow() time.Time {
	loc, _ := time.LoadLocation("Europe/Vienna")
	year, month, day := time.Now().AddDate(0, 0, 1).Date()
	return time.Date(year, month, day, 0, 0, 0, 0, loc)
}

func containsImportantKeyword(toCheck *string) bool {
	for _, keyword := range ImportantKeywords {
		if strings.Contains(*toCheck, keyword) {
			return true
		}
	}
	return false
}
