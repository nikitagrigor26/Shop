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
	Query := `SELECT Id, Name, Price, Description,Quantity FROM Products WHERE 1=1`
	if idfilter != "" {
		id, err = strconv.Atoi(idfilter)
		if err != nil {
			fmt.Println("Ошибка конвертации Id", err)
			fmt.Fprintf(w, "Не удалось обработать Id")
			return
		}
		args = append(args, id)
		Query += fmt.Sprintf(" AND Id = @p%d", len(args))
	}
	if nameFilter != "" {
		args = append(args, nameFilter)
		Query += fmt.Sprintf(" AND Name = @p%d", len(args))
	}
	if priceFilter != "" {
		price, err = strconv.Atoi(priceFilter)
		if err != nil {
			fmt.Println("Ошибка конвертации Price", err)
			fmt.Fprintf(w, "Не удалось обработать Price")
			return
		}
		args = append(args, price)
		Query += fmt.Sprintf(" AND Price = @p%d", price)
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
		err = rows.Scan(&prod.Id, &prod.Name, &prod.Price, &prod.Description, &prod.Quantity)
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
	idfilter := r.PathValue("id")
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
	Query := `SELECT Id, Name, Quantity, Price, Description FROM Products WHERE Id=@p1`
	err = database.Db.QueryRow(Query, id).Scan(&prod.Id, &prod.Name, &prod.Quantity, &prod.Price, &prod.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Товар не найден", http.StatusNotFound)
			return
		}
		fmt.Println(err)
		http.Error(w, "Ошибка с БД, запрос пока не доступен", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(prod)
}
func AddProductHandler(w http.ResponseWriter, r *http.Request) {
	var newprod Models.Product
	err := json.NewDecoder(r.Body).Decode(&newprod)
	if err != nil {
		fmt.Println("возникла ошибка с обработкой product :", err)
		return
	}
	Query := `INSERT INTO Products (Name, Price, Description, Quantity) VALUES (@p1, @p2, @p3, @p4)`
	_, err = database.Db.Exec(Query, newprod.Name, newprod.Price, newprod.Description, newprod.Quantity)
	if err != nil {
		fmt.Println("Ошибка добавление записи в БД : ", err)
		fmt.Fprintf(w, "Ошибка добавление в БД")
		return
	}
	fmt.Fprintf(w, "Книга добавлена")
}
func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	textId := r.PathValue("id")
	var id int
	var err error
	if textId != "" {
		id, err = strconv.Atoi(textId)
		if err != nil {
			fmt.Println("Ошибка чтения Id : ", err)
			fmt.Fprintf(w, "Ошибка конвертации Id, проверьте запрос")
			return
		}
	}
	Query := `DELETE FROM Products WHERE Id=@p1`
	result, err := database.Db.Exec(Query, id)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "Ошибка БД: %v", err)
		return
	}
	rows, error := result.RowsAffected()
	if error != nil {
		fmt.Fprintf(w, "Ошибка удаления")
		return
	}
	if rows == 0 {
		fmt.Fprintf(w, "Ошибка удаления, продукт не найден")
		return
	}
	fmt.Fprintf(w, "Товар удалена")
}
func UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	textId := r.PathValue("id")
	var id int
	var err error
	if textId != "" {
		id, err = strconv.Atoi(textId)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "Ошибка конвертации Id, проверьте запрос")
			return
		}
	}
	var updateProd Models.Product
	err = json.NewDecoder(r.Body).Decode(&updateProd)
	if err != nil {
		fmt.Println(err)
		return
	}
	Query := `UPDATE Products SET Name =@p2, Price = @p3, Description=@p4, Quantity =@p5 WHERE Id=@p1`
	res, errs := database.Db.Exec(Query, id, updateProd.Name, updateProd.Price, updateProd.Description, updateProd.Quantity)
	if errs != nil {
		fmt.Println(err)
		fmt.Println(w, "Ошибка изменений")
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	if rowsAffected == 0 {
		fmt.Fprintf(w, "Ошибка, товара не найден")
	}
	fmt.Fprintf(w, "Товар обновлен")
}
