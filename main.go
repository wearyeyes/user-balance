package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
)

var db *sql.DB

func main() {
	// Open a database called 'user_balance' using postgreSQL.
	var err error
	db, err = sql.Open("postgres", "user=postgres password=password1 dbname=user_balance sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	// Close a database in the end of program's work.
	defer db.Close()

	// Create new router with default widdleware,
	// using gin framework.
	r := gin.Default()

	// Handlers to different paths.
	r.GET("/balance/:id", getBalance)
	r.PUT("/topup", topUp)
	r.PUT("/withdraw", withdraw)
	r.PUT("/transfer", transfer)

	// Run the server on 8090 port.
	r.Run(":8090")
}
