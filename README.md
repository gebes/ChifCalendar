# ChifCalendar
This bot fetches a predefined calendar every day and posts a summary to a Discord server.
![](https://i.imgur.com/7Ub5mmk.png)

## Requirements
1. Discord Webhook with ID and URL
2. Direct URL to the .ics file from the calendar
3. MySQL Database Connection URL

Create a .env file and fill in these details like the `.env.sample` file.

### Database
Create a "message" Table
```mysql
CREATE TABLE message (
    id INT NOT NULL AUTO_INCREMENT,
    hashValue VARCHAR(256),
    PRIMARY KEY(id)
);
```

## How to use?
Create an event in the calendar with one of the following emojis as the prefix. The bot will categorize the events accordingly with the prefixes stripped of them.
```go
HomeworkEmoji = "📘"
ExamEmoji     = "📝"
OtherEmoji    = "👻"
```
![](https://i.imgur.com/2JOx7uR.png)
![](https://i.imgur.com/7Ub5mmk.png)
### Hyperlinks
If the description contains a URL that redirects to a Discord Message, then the webhook creates a hyperlink.

### Duplicate messages
The bot requires a database to store each message's SHA256 hash. Saving the hash values makes it possible to prevent duplicate messages. In other words, if the message would look the same as the last message sent, then the bot will not send anything.
