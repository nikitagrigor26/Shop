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
	mux.HandleFunc("GET /", Handlers.HelloHandler)
	mux.HandleFunc("GET /products/", Handlers.Productshandler)
	mux.HandleFunc("GET /products/{id}", Handlers.ProductHandler)
	mux.HandleFunc("POST /products/add/", Handlers.AddProductHandler)
	mux.HandleFunc("GET /users/", Handlers.ViewUsersHandler)
	mux.HandleFunc("GET /user/{id}", Handlers.ViewUserhandler)
	fmt.Println("Серваер запущен на порте http://localhost:8000")
	http.ListenAndServe(":8000", loggedmux)
}
