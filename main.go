package main

import (
	"agentoff/internals/server/database"
	"agentoff/internals/server/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	database.InitDB()
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/osago", handlers.OsagoHandler)
	http.HandleFunc("/kasko", handlers.KaskoHandler)
	http.HandleFunc("/house", handlers.HouseHandler)
	http.HandleFunc("/dom", handlers.DomHandler)
	http.HandleFunc("/contact", handlers.ContactHandler)
	http.HandleFunc("/admin", handlers.AdminHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Server starting at :8088")
	log.Fatal(http.ListenAndServe(":8088", nil))
}
