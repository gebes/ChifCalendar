package webhook

import (
	"fmt"
	"gebes.io/calendar/database"
	"gebes.io/calendar/utils"
	"github.com/apognu/gocal"
	"github.com/nickname32/discordhook"
	"github.com/robfig/cron"
	"log"
	"mvdan.cc/xurls"
	"net/http"
	"sort"
	"strconv"
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

func InitCronJob() {
	log.Println("Initializing CronJob")
	location, err := time.LoadLocation("Europe/Vienna")
	if err != nil {
		panic(err);
	}
	c := cron.NewWithLocation(location)
	c.AddFunc("0 0 "+strconv.Itoa(MessageTime)+" * * *", sendEventsMessage)
	c.Start()
	select {}
}

func sendEventsMessage() {

	tomorrowEmbed, tomorrowHash := getEmbedFor("Was steht morgen an?", fetchTomorrowEvents(), false)
	nextWeekEmbed, nextWeekHash := getEmbedFor("Was steht in den nächsten Tagen an?", fetchNextWeekEvents(), true)

	hash := tomorrowHash + nextWeekHash
	lastHash := database.LastMessageHash()

	if lastHash != nil && hash == *lastHash {
		return
	}

	database.SaveHash(hash)

	sendEmbed(tomorrowEmbed)
	sendEmbed(nextWeekEmbed)

}


func sendEmbed(embed *discordhook.Embed) {
	if embed == nil {
		return
	}
	_, err := WA.Execute(nil, &discordhook.WebhookExecuteParams{
		Embeds: []*discordhook.Embed{
			embed,
		},
	}, nil, "")
	if err != nil {
		panic(err)
	}
}

func getEmbedFor(title string, events *[]gocal.Event, showDates bool) (*discordhook.Embed, string) {
	messages := make([]string, len(Categories))
	important := false

	embedFields := []*discordhook.EmbedField{}
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
		URL:    "https://calendar.google.com/calendar/u/0?cid=aGEydTh1Z3Z2bDFrMWVucGhoM3M0cTE4bTBAZ3JvdXAuY2FsZW5kYXIuZ29vZ2xlLmNvbQ",
	}

	if important {
		embed.Color = 0xf2c94c
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
	if len(eventToAppend.Description) != 0 && strings.Contains(eventToAppend.Description, "https://discord.com/channels"){
		url := xurls.Relaxed.FindString(eventToAppend.Description)
		*message += "["+summary + "]("+ url +")"
	} else {
		*message += summary
	}

	*message += "\n"
}

func fetchTomorrowEvents() *[]gocal.Event {
	start, end := tomorrow(), tomorrow().Add(24*time.Hour)
	return fetchEventsInPeriod(&start, &end)
}
func fetchNextWeekEvents() *[]gocal.Event {
	start, end := tomorrow().Add(48*time.Hour), tomorrow().Add(8*24*time.Hour)
	return fetchEventsInPeriod(&start, &end)
}

func fetchEventsInPeriod(start *time.Time, end *time.Time) *[]gocal.Event {
	resp, err := http.Get(utils.GetCalendarKeys().CalenderURL)
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
