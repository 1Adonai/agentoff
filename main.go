package main

import (
	"agentoff/internals/server/database"
	"agentoff/internals/server/handlers"
	"agentoff/internals/server/logger"
	"fmt"
	"log"
	"net/http"
)

func main() {
	logger.InitLogger()
	defer logger.CloseLogger()
	database.InitDB()
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/osago", handlers.OsagoHandler)
	http.HandleFunc("/kasko", handlers.KaskoHandler)
	http.HandleFunc("/house", handlers.HouseHandler)
	http.HandleFunc("/dom", handlers.DomHandler)
	http.HandleFunc("/contact", handlers.ContactHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Server starting at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
