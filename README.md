# ChifCalendar
A Discord webhook tool that fetches a calendar and posts a pretty summary every day. 
![](https://i.imgur.com/z9TslZm.png)

## Requirements
1. Discord Webhook with ID and URL
2. Direct URL to the .ics file from a calendar
3. MySQL Database with Connection URL

Create a .env file and fill in these details like the `.env.sample` file.

### Database Setup
Create a `message` Table
```mysql
CREATE TABLE message (
    id INT NOT NULL AUTO_INCREMENT,
    webhookId INT NOT NULL,
    hashValue VARCHAR(256),
    PRIMARY KEY(id),
    FOREIGN KEY(webhookId) REFERENCES webhook(id) ON DELETE CASCADE
);
```

Create a `webhook` Table

```mysql
CREATE TABLE webhook (
    id INT NOT NULL AUTO_INCREMENT,
    discordId BIGINT NOT NULL,
    discordToken VARCHAR(512),
    calendarPrivateUrl VARCHAR(2048),
    calendarPublicUrl VARCHAR(2048),
    nickname VARCHAR(256),
    primaryColor INT,
    PRIMARY KEY(id)
);
```

Register a webhook
```mysql
INSERT INTO webhook(discordId, discordToken, calendarPrivateUrl, calendarPublicUrl, nickname, primaryColor) 
VALUES(123, "abc",
"https://calendar.url/data.ics",
"https://publiccalendar.url",
"Bot Name",
0xf2c94c);
```

## How to use?
Create an event in the calendar with one of the following emojis as the prefix. The bot will categorize the events accordingly with the prefixes stripped of them. 
```go
HomeworkEmoji = "📘"
ExamEmoji     = "📝"
OtherEmoji    = "👻"
```
I would recommend setting a calendar entry to "all day".

![](https://i.imgur.com/2JOx7uR.png)
![](https://i.imgur.com/7Ub5mmk.png)
### Hyperlinks
The webhook creates a hyperlink if the description contains a URL that redirects to a Discord message.

### Duplicate messages
The bot uses the `message` table to store a SHA256 hash value of each message. Every day at 17:00, the bot creates a message, fetches the last hash value, and compares it to the current hash value. If the hash is different, the bot will send the message and store the hash in the database.
