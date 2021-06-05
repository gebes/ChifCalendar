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

func LastMessageHash() *string {

	results, err := Database.Query("SELECT * FROM message ORDER BY message.id DESC LIMIT 1;")

	if err != nil {
		panic(err)
	}

	var lastMessage Message

	results.Next()
	err = results.Scan(&lastMessage.id, &lastMessage.hashValue)
	if err != nil {
		return nil
	}

	return &lastMessage.hashValue
}

func SaveHash(toSave string) {

	_, err := Database.Query("INSERT INTO message(message.hashValue) VALUES(\"" + toSave + "\");")

	if err != nil {
		panic(err)
	}

}
