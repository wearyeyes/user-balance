package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
)

type User struct {
	Id      int `json:"id"`
	Balance int `json:"balance"`
}

// Handler that allow to find out the user's balance.
// http://localhost:8090/balance/:id
func getBalance(c *gin.Context) {
	id := c.Param("id")

	var balance int
	err := db.QueryRow("SELECT balance FROM users WHERE id=$1", id).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("User with this ID doesn't exist.")
			return
		} else {
			log.Fatal(err)
		}	
	}

	c.JSON(200, gin.H{
		"id": id,
		"balance": balance,
	})
}

// It allows to top up the user's balance.
func topUp(c *gin.Context) {
	var user *User
	err := c.BindJSON(&user)
	if err != nil {
		log.Println(err)
		return
	}

	var currentBalance int
	err = db.QueryRow("SELECT balance FROM users WHERE id=$1", user.Id).Scan(&currentBalance)
	if err != nil {
		_, err := db.Exec("INSERT INTO users(balance) VALUES($1)", user.Balance)
		if err != nil {
			log.Println(err)
		}

		c.JSON(200, gin.H{
			"id":      user.Id,
			"balance": user.Balance,
		})
	} else {
		currentBalance += user.Balance
		_, err = db.Exec("UPDATE users SET balance=$1 WHERE id=$2", currentBalance, user.Id)
		if err != nil {
			log.Println(err)
		}

		c.JSON(200, gin.H{
			"id":      user.Id,
			"balance": currentBalance,
		})
	}
}

// This serves to withdraw money from the user's account.
func withdraw(c *gin.Context) {
	var user *User
	err := c.BindJSON(&user)
	if err != nil {
		log.Println(err)
	}

	var currentBalance int
	err = db.QueryRow("SELECT balance FROM users WHERE id=$1", user.Id).Scan(&currentBalance)
	if err != nil {
		log.Println(err)
	}

	currentBalance -= user.Balance
	if currentBalance < 0 {
		c.JSON(200, gin.H{
			"message": "unsufficient funds on your account",
		})
		return
	}

	_, err = db.Exec("UPDATE users SET balance=$1 WHERE id=$2", currentBalance, user.Id)
	if err != nil {
		log.Println(err)
	}

	c.JSON(200, gin.H{
		"id":      user.Id,
		"balance": currentBalance,
	})
}

// This function transfers money from one user to another.
func transfer(c *gin.Context) {
	var tInfo = struct {
		SenderId    int `json:"senderId"`
		RecipientId int `json:"recipientId"`
		Money       int `json:"money"`
	}{}

	err := c.BindJSON(&tInfo)
	if err != nil {
		log.Println(err)
	}

	var senderCurrentBalance int
	err = db.QueryRow("SELECT balance FROM users WHERE id=$1", tInfo.SenderId).Scan(&senderCurrentBalance)
	if err != nil {
		log.Println(err)
	}

	senderCurrentBalance -= tInfo.Money
	if senderCurrentBalance < 0 {
		c.JSON(200, gin.H{
			"message": "unsufficient funds on your account",
		})
		return
	}

	_, err = db.Exec("UPDATE users SET balance=$1 WHERE id=$2", senderCurrentBalance, tInfo.SenderId)
	if err != nil {
		log.Println(err)
	}

	var recipientCurrentBalance int
	err = db.QueryRow("SELECT balance FROM users WHERE id=$1", tInfo.RecipientId).Scan(&recipientCurrentBalance)
	if err != nil {
		log.Println(err)
	}

	recipientCurrentBalance += tInfo.Money
	_, err = db.Exec("UPDATE users SET balance=$1 WHERE id=$2", recipientCurrentBalance, tInfo.RecipientId)
	if err != nil {
		log.Println(err)
	}

	c.JSON(200, gin.H{
		"id":      tInfo.SenderId,
		"balance": senderCurrentBalance,
	})
}