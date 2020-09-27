package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

type User struct {
	Id      int     `json:"id"`
	Balance float64 `json:"balance"`
}

// Handler that allow to find out the user's balance.
// http://localhost:8090/balance/:id?currency=RUB
func getBalance(c *gin.Context) {
	id := c.Param("id")
	currency := c.DefaultQuery("currency", "RUB")

	var balance float64
	err := db.QueryRow("SELECT balance FROM users WHERE id=$1", id).Scan(&balance)
	if err != nil {
		log.Println(err)
		return
	}

	switch {
	case currency == "RUB":
		c.JSON(200, gin.H{
			"id":      id,
			"balance": balance,
		})
	case currency != "RUB":
		balance, err := currencyConv(balance, currency)
		if err != nil {
			c.JSON(400, gin.H{
				"message": fmt.Sprint(err),
			})
			return
		}
		c.JSON(200, gin.H{
			"id":      id,
			"balance": balance,
		})
	}
}

// It allows to top up the user's balance.
func topUp(c *gin.Context) {
	var user *User
	err := c.BindJSON(&user)
	if err != nil {
		log.Println(err)
		return
	}

	var currentBalance float64
	err = db.QueryRow("SELECT balance FROM users WHERE id=$1", user.Id).Scan(&currentBalance)
	if err != nil {
		// If database have no such user, create new row.
		_, err := db.Exec("INSERT INTO users(balance) VALUES($1)", user.Balance)
		if err != nil {
			log.Println(err)
			return
		}

		c.JSON(200, gin.H{
			"id":      user.Id,
			"balance": user.Balance,
		})
	} else {
		// Add money to current user's balance.
		currentBalance += user.Balance
		_, err := db.Exec("UPDATE users SET balance=$1 WHERE id=$2", currentBalance, user.Id)
		if err != nil {
			log.Println(err)
			return
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
		return
	}

	var currentBalance float64
	err = db.QueryRow("SELECT balance FROM users WHERE id=$1", user.Id).Scan(&currentBalance)
	if err != nil {
		log.Println(err)
		return
	}

	// Subtracting money from the current balance and update balance in database.
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
		return
	}

	c.JSON(200, gin.H{
		"id":      user.Id,
		"balance": currentBalance,
	})
}

// This function transfers money from one user to another.
// Accept sender's ID, recipient's ID and money which sender want to send.
// Then check the current balance of the sender, and if he has enough money
// in the account - func sends it.
func transfer(c *gin.Context) {
	var tInfo = struct {
		SenderId    int     `json:"senderId"`
		RecipientId int     `json:"recipientId"`
		Money       float64 `json:"money"`
	}{}

	err := c.BindJSON(&tInfo)
	if err != nil {
		log.Println(err)
		return
	}

	var senderCurrentBalance float64
	err = db.QueryRow("SELECT balance FROM users WHERE id=$1", tInfo.SenderId).Scan(&senderCurrentBalance)
	if err != nil {
		log.Println(err)
		return
	}

	var recipientCurrentBalance float64
	err = db.QueryRow("SELECT balance FROM users WHERE id=$1", tInfo.RecipientId).Scan(&recipientCurrentBalance)
	if err != nil {
		log.Println(err)
		return
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
		return
	}

	recipientCurrentBalance += tInfo.Money
	_, err = db.Exec("UPDATE users SET balance=$1 WHERE id=$2", recipientCurrentBalance, tInfo.RecipientId)
	if err != nil {
		log.Println(err)
		return
	}

	c.JSON(200, gin.H{
		"id":      tInfo.SenderId,
		"balance": senderCurrentBalance,
	})
}
