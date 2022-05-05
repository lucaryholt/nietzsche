package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"

	"github.com/joho/godotenv"
)

type DBMessage struct {
	ID        int    `json:"id"`
	Message   string `json:"message"`
	Topic     int    `json:"topic"`
	CreatedAt string `json:"createdAt"`
}

type Message struct {
	message string `json:"message"`
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func makeConnectionString() string {
	var user = os.Getenv("DB_USER")
	var password = os.Getenv("DB_PASSWORD")
	var ip = os.Getenv("DB_IP")
	var port = os.Getenv("DB_PORT")
	var name = os.Getenv("DB_NAME")

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		user,
		password,
		ip,
		port,
		name,
	)
}

func authenticate(token string) bool {
	return token == "testToken"
}

func getMessage(c *gin.Context) {
	topic := c.Param("topic")
	offset := c.Param("offset")
	limit := c.Param("limit")
	token := c.Param("token")
	// format := c.Param("format")

	if !authenticate(token) {
		c.String(http.StatusUnauthorized, "Not authorized. Provide valid token.")
		return
	}

	db, connectionError := sql.Open("mysql", makeConnectionString())

	if connectionError != nil {
		panic(connectionError.Error())
	}

	defer db.Close()

	results, queryError := db.Query(
		"SELECT * FROM message " +
			"WHERE topic = " +
			"(SELECT id FROM topic WHERE name = '" + topic + "') " +
			"AND id > " + offset + " LIMIT " + limit,
	)
	if queryError != nil {
		panic(queryError.Error())
	}

	var messages []Message

	for results.Next() {
		var message DBMessage

		var scanningError = results.Scan(&message.ID, &message.Message, &message.Topic, &message.CreatedAt)
		if scanningError != nil {
			panic(scanningError.Error())
		}

		var messageData bytes.Buffer
		jsonFormatError := json.Indent(&messageData, []byte(message.Message), "", "\t")
		if jsonFormatError != nil {
			panic(jsonFormatError.Error())
		}

		var returnMessage Message
		returnMessage.message = messageData.String()

		// fmt.Println(returnMessage)

		messages = append(messages, returnMessage)
	}

	fmt.Println(messages)

	c.IndentedJSON(http.StatusOK, messages)
}

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		panic("Error loading .env file")
	}
}

func main() {
	router := gin.Default()
	router.GET("/ping", ping)
	router.GET("/read-messages/topic/:topic/from/:offset/limit/:limit/user-token/:token/format/:format", getMessage)

	router.Run("localhost:8080")
}
