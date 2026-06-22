package Handlers

import (
	database "Shop/Database"
	"Shop/Models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func Productshandler(w http.ResponseWriter, r *http.Request) {
	idfilter := r.URL.Query().Get("id")
	nameFilter := r.URL.Query().Get("name")
	priceFilter := r.URL.Query().Get("price")
	var price int
	var id int
	var err error
	var args []any
	Query := `SELECT Id, Name, Email, Password FROM Products WHERE 1=1`
	if idfilter != "" {
		id, err = strconv.Atoi(idfilter)
		if err != nil {
			fmt.Println("Ошибка конвертации Id", err)
			fmt.Fprintf(w, "Не удалось обработать Id")
			return
		}
		args = append(args, id)
		Query += fmt.Sprintf(" AND Id = %d", len(args))
	}
	if nameFilter != "" {
		args = append(args, nameFilter)
		Query += fmt.Sprintf(" AND Name = %s", len(args))
	}
	if priceFilter != "" {
		price, err = strconv.Atoi(priceFilter)
		if err != nil {
			fmt.Println("Ошибка конвертации Price", err)
			fmt.Fprintf(w, "Не удалось обработать Price")
			return
		}
		args = append(args, price)
		Query += fmt.Sprintf(" AND Price = %d", price)
	}
	rows, errs := database.Db.Query(Query, args...)
	if errs != nil {
		fmt.Println("Какая-то ошибка с БД", err)
		return
	}
	defer rows.Close()
	var Products []Models.Product
	for rows.Next() {
		var prod Models.Product
		err = rows.Scan(&prod.Id, &prod.Name, prod.Price, prod.Description, prod.Quantity)
		if err != nil {
			fmt.Println(err)
			continue
		}
		Products = append(Products, prod)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(Products)
}
func ProductHandler(w http.ResponseWriter, r *http.Request) {
	idfilter := r.URL.Query().Get("id")
	var id int
	var err error
	if idfilter != "" {
		id, err = strconv.Atoi(idfilter)
		if err != nil {
			fmt.Println("Ошибка конвертации Id", err)
			fmt.Fprintf(w, "Ошибка с запросом, проверьте запрос")
			return
		}
	}
	var prod Models.Product
	Query := `SELECT Name,Price,Description,Quntity FROM Products WHERE Id=@p1`
	err = database.Db.QueryRow(Query, id).Scan(prod.Id, prod.Name, prod.Quantity, prod.Price, prod.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Книга не найдена", http.StatusNotFound)
			return
		}
		fmt.Println(err)
		http.Error(w, "Ошибка с БД, запрос пока не доступен", http.StatusServiceUnavailable)
		return
	}

}
func AddProductHandler(w http.ResponseWriter, r *http.Request) {
	var newprod Models.Product
	err := json.NewDecoder(r.Body).Decode(&newprod)
	if err != nil {
		fmt.Println("возникла ошибка с обработкой product :", err)
		return
	}
	Query := `INSERT INTO Products Id,Name,Price,Description,Quntity values (@p1,@p2,@p3,@p4,@p5)`
	_, err = database.Db.Exec(Query, newprod.Id, newprod.Name, newprod.Price, newprod.Description, newprod.Quantity)
	if err != nil {
		fmt.Println("Ошибка добавление записи в БД : ", err)
		fmt.Fprintf(w, "Ошибка добавление в БД")
		return
	}
	fmt.Fprintf(w, "Книга добавлена")
}
