# ChifCalendar
Simple webhook which fetches a calendar daily and posts a summary to a channel.
![](https://i.imgur.com/7Ub5mmk.png)

## Requirements
1. Discord Webhook with ID and URl
2. Direct URL to the .ics file from the calendar
3. MySQL Database Connection URL

Create a .env file and fill in these details like the `.env.sample` file.

### Database
Simply create a "message" Table
```mysql
CREATE TABLE message (
    id INT NOT NULL AUTO_INCREMENT,
    hashValue VARCHAR(256),
    PRIMARY KEY(id)
);
```

## How to use?
Create an event in your calendar. In the message an event will only show up if an emoji is provided.
```go
HomeworkEmoji = "📘"
ExamEmoji     = "📝"
OtherEmoji    = "👻"
```
It will then create a message based on your calendar.
![](https://i.imgur.com/2JOx7uR.png)
![](https://i.imgur.com/7Ub5mmk.png)  
### Hyperlinks
If the description contains a URL which redirects to a Discord Message, then the webhook creates a hyperlink.