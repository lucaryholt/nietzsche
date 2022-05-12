package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"

	"github.com/joho/godotenv"
)

type Message struct {
	ID        int    `form:"id" json:"id" xml:"id" yaml:"id"`
	Content   string `form:"content" json:"content" xml:"content" yaml:"content"`
	Topic     string `form:"topic" json:"topic" xml:"topic" yaml:"topic"`
	CreatedAt string `form:"createdAt" json:"createdAt" xml:"createdAt" yaml:"createdAt"`
	Creator   int    `form:"creator" json:"creator" xml:"creator" yaml:"creator"`
}

type TokenInformation struct {
	ID    int
	Email string
	Phone int
	Token string
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func isSupportedFormat(format string) bool {
	formats := [4]string{"JSON", "XML", "YAML", "TSV"}

	for _, supportedFormat := range formats {
		if supportedFormat == format {
			return true
		}
	}
	return false
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

func authenticate(inputToken string, c *gin.Context) bool {
	return authenticateToken(inputToken, c)
}

func getMessage(c *gin.Context) {
	topic := c.Param("topic")
	offset := c.Param("offset")
	limit := c.Param("limit")
	token := c.Param("token")
	format := c.Param("format")

	if !authenticate(token, c) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized. Provide valid token."})
		return
	}

	if !isSupportedFormat(format) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Not a supported format"})
		return
	}

	messages := getMessages(topic, offset, limit, c)

	if len(messages) == 0 {
		c.JSON(http.StatusNoContent, "")
		return
	}

	transformMessages(format, c, messages)
}

func transformMessages(format string, c *gin.Context, messages []Message) {
	switch format {
	case "JSON":
		c.JSON(http.StatusOK, messages)
	case "XML":
		c.XML(http.StatusOK, messages)
	case "YAML":
		c.YAML(http.StatusOK, messages)
	case "TSV":
		var returnMessages bytes.Buffer

		returnMessages.WriteString("id,content,topic,createdAt,creator\n")

		for _, message := range messages {
			returnMessages.WriteString(strconv.Itoa(message.ID) + "," + message.Content + "," + message.Topic + "," + message.CreatedAt + "," + strconv.Itoa(message.Creator) + "\n")
		}
		c.String(http.StatusOK, returnMessages.String())
	}
}

func createMessage(c *gin.Context) {
	token := c.Param("token")
	format := c.Param("format")

	if !authenticate(token, c) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized. Provide valid token."})
		return
	}

	if !isSupportedFormat(format) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Not a supported format"})
		return
	}

	message := Message{}

	if format == "TSV" {
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Fatal(err.Error())
		}

		lines := strings.Split(string(data), "\n")
		messageData := strings.Split(lines[1], ",")

		message.Content = messageData[0]
		message.Topic = messageData[1]
	} else {
		c.Bind(&message)

		if message.Content == "" || message.Topic == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Provide both content and topic."})
			return
		}
	}

	insertMessage(message, token, c)

	c.JSON(http.StatusOK, gin.H{"message": "Message stored."})
}

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	router := gin.Default()
	router.GET("/ping", ping)
	router.GET("/read-messages/topic/:topic/from/:offset/limit/:limit/user-token/:token/format/:format", getMessage)
	router.POST("/create-message/user-token/:token/format/:format", createMessage)

	router.Run(":" + os.Getenv("PORT"))
}
