package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getDatabaseConnection(c *gin.Context) *sql.DB {
	db, connectionError := sql.Open("mysql", makeConnectionString())

	if connectionError != nil {
		var errorMessage = "Could not connect to database"
		c.JSON(http.StatusInternalServerError, gin.H{"message": errorMessage})
		log.Fatal(errorMessage, connectionError.Error())
	}

	return db
}

func authenticateToken(inputToken string, c *gin.Context) bool {
	db := getDatabaseConnection(c)
	defer db.Close()

	statement := "SELECT * FROM token " +
		"WHERE token = ?;"

	var id int
	var email string
	var phone int
	var token string

	queryError := db.QueryRow(statement, inputToken).Scan(&id, &email, &phone, &token)
	if queryError != nil {
		return false
	}

	return true
}

func getMessages(topic string, offset string, limit string, c *gin.Context) []Message {
	db := getDatabaseConnection(c)
	defer db.Close()

	statement, statementCreationError := db.Prepare(
		"SELECT * FROM message " +
			"WHERE topic = ? " +
			"AND id > ? LIMIT ?",
	)
	if statementCreationError != nil {
		var errorMessage = "Error preparing database query"
		c.JSON(http.StatusInternalServerError, gin.H{"message": errorMessage})
		log.Fatal(errorMessage, statementCreationError.Error())
	}

	results, queryError := statement.Query(topic, offset, limit)
	if queryError != nil {
		var errorMessage = "Error querying database"
		c.JSON(http.StatusInternalServerError, gin.H{"message": errorMessage})
		log.Fatal(errorMessage, queryError.Error())
	}

	var messages []Message

	for results.Next() {
		var message Message

		var scanningError = results.Scan(&message.ID, &message.Content, &message.Topic, &message.CreatedAt, &message.Creator)
		if scanningError != nil {
			var errorMessage = "Error scanning results from database"
			c.JSON(http.StatusInternalServerError, gin.H{"message": errorMessage})
			log.Fatal(errorMessage, scanningError.Error())
		}

		messages = append(messages, message)
	}

	return messages
}

func insertMessage(message Message, token string, c *gin.Context) {
	db := getDatabaseConnection(c)
	defer db.Close()

	statement, statementCreationError := db.Prepare(
		"INSERT INTO message (content, topic, creator) VALUES (?, ?, (SELECT id FROM token WHERE token = ?))",
	)
	if statementCreationError != nil {
		var errorMessage = "Error preparing database query"
		c.JSON(http.StatusInternalServerError, gin.H{"message": errorMessage})
		log.Fatal(errorMessage, statementCreationError.Error())
	}

	_, insertError := statement.Exec(message.Content, message.Topic, token)
	if insertError != nil {
		var errorMessage = "Error inserting into database"
		c.JSON(http.StatusInternalServerError, gin.H{"message": errorMessage})
		log.Fatal(errorMessage, insertError.Error())
	}
}
