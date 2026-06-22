package main

import (
	"Shop/Database"
	"Shop/Handlers"
	middleware "Shop/Middleware"
	"Shop/Models"
	"fmt"
	"net/http"
)

func main() {
	database.InitDb()
	mux := http.NewServeMux()
	loggedmux := middleware.LogingMiddelware(mux)
	go Models.StatsUsers()
	go Models.StatsProduct()
	mux.HandleFunc("/", Handlers.HelloHandler)
	mux.HandleFunc("/products/", Handlers.Productshandler)
	mux.HandleFunc("/products/{id}", Handlers.ProductHandler)
	mux.HandleFunc("/products/add/", Handlers.AddProductHandler)
	fmt.Println("Серваер запущен на порте http://localhost:8000")
	http.ListenAndServe(":8000", loggedmux)
}
