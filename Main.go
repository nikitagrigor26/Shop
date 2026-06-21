package main

import (
	"Shop/Database"
	"Shop/Handlers"
	"fmt"
	"net/http"
)

func main() {
	database.InitDb()
	mux := http.NewServeMux()
	mux.HandleFunc("/", Handlers.HelloHandler)

	http.ListenAndServe(":8000", mux)
	fmt.Println("Серваер запущен на порте http:/localhost:8000")
}
