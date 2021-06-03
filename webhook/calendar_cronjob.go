package webhook

import (
	"github.com/Gebes/ChifCalendar/utils"
	"github.com/apognu/gocal"
	"github.com/nickname32/discordhook"
	"github.com/robfig/cron"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	HomeworkEmoji = "📘"
	ExamEmoji     = "📝"
	Other         = "👻"
)

func InitCronJob() {
	log.Println("Initializing CronJob")
	c := cron.New()
	c.AddFunc("0 17 * * *", writeEvents)
	c.Start()
	log.Println("Initialized CronJob")
	select {}
}

func writeEvents() {
	exams, homework, other := "", "", ""
	for _, event := range fetchTomorrowEvents() {
		if strings.Contains(event.Summary, HomeworkEmoji) {
			event.Summary = strings.TrimPrefix(event.Summary, HomeworkEmoji)
			addRow(&homework, event)
		} else if strings.Contains(event.Summary, ExamEmoji) {
			event.Summary = strings.TrimPrefix(event.Summary, ExamEmoji)
			addRow(&exams, event)
		} else {
			addRow(&other, event)
		}
	}

	if len(exams) == 0 && len(homework) == 0 && len(other) == 0 {
		return
	}

	// important := strings.Contains(exams, "Schularbeit")
	important := len(exams) != 0

	embedFields := []*discordhook.EmbedField{}

	if len(exams) != 0 {
		embedFields = append(embedFields, &discordhook.EmbedField{
			Name:   ExamEmoji + " Prüfungen",
			Value:  exams,
			Inline: false,
		})
	}

	if len(homework) != 0 {
		embedFields = append(embedFields, &discordhook.EmbedField{
			Name:   HomeworkEmoji + " Hausübungen",
			Value:  homework,
			Inline: false,
		})
	}
	if len(other) != 0 {
		embedFields = append(embedFields, &discordhook.EmbedField{
			Name:   Other + " Anderes",
			Value:  other,
			Inline: false,
		})
	}

	content := &discordhook.WebhookExecuteParams{
		Embeds: []*discordhook.Embed{
			{
				Title:  "Was steht morgen an?",
				Fields: embedFields,
				URL:    "https://calendar.google.com/calendar/u/0?cid=aGEydTh1Z3Z2bDFrMWVucGhoM3M0cTE4bTBAZ3JvdXAuY2FsZW5kYXIuZ29vZ2xlLmNvbQ",
				Footer: &discordhook.EmbedFooter{
					IconURL: "https://i.imgur.com/meY76jE.png",
					Text: "Calendar",
				},

			},
		},
	}

	if important {
		content.Content = "@everyone"
		content.Embeds[0].Color = 0xff0000
	}

	msg, err := WA.Execute(nil, content, nil, "")
	if err != nil {
		panic(err)
	}
	log.Println(msg.ID)
}

func addRow(message *string, toAppend gocal.Event) {
	*message += "• "+ toAppend.Summary
	if len(toAppend.Description) != 0 {
		*message += " (" + toAppend.Description + ")"
	}
	*message += "\n"
}

func fetchTomorrowEvents() []gocal.Event {
	start, end := time.Now(), time.Now().Add(24*time.Hour)

	resp, err := http.Get(utils.GetCalendarKeys().CalenderURL)
	if err != nil {
		panic(err)
	}

	c := gocal.NewParser(resp.Body)

	c.Start, c.End = &start, &end
	err = c.Parse()
	if err != nil {
		panic(err)
	}

	return c.Events
}
