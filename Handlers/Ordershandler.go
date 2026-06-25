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

func ViewAllOrders(w http.ResponseWriter, r *http.Request) {
	idfilter := r.URL.Query().Get("id")

	// 1. ИСПРАВЛЕНО: Запрашиваем только те колонки, которые будем сканировать
	Query := `SELECT Id, Status FROM Orders WHERE 1=1`

	var id int
	var err error
	var args []any

	if len(idfilter) > 0 {
		id, err = strconv.Atoi(idfilter)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Ошибка конвертации Id, проверьте запрос", http.StatusBadRequest)
			return
		}
		args = append(args, id)
		Query += fmt.Sprintf(" AND Id = @p%d", len(args))
	}

	rows, err := database.Db.Query(Query, args...)
	if err != nil {
		fmt.Println("Ошибка Query:", err)
		http.Error(w, "Ошибка подключения к БД", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []Models.Order
	for rows.Next() {
		var order Models.Order
		err = rows.Scan(&order.OrderId, &order.Status)
		if err != nil {
			fmt.Println("Ошибка Scan при чтении заказа:", err)
			continue
		}
		orders = append(orders, order)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(orders)
}
func AddOrder(w http.ResponseWriter, r *http.Request) {
	var order Models.OrderRequest
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Ошибка чтения запроса, проверьте формат JSON", http.StatusBadRequest)
		return
	}
	tx, err := database.Db.Begin()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	var userIdCheck int
	checkQueryUser := `SELECT Id FROM Users WHERE Id = @p1`
	err = tx.QueryRow(checkQueryUser, order.UserID).Scan(&userIdCheck)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Пользователя с таким Id не существует", http.StatusNotFound)
		} else {
			fmt.Println("Ошибка проверки пользователя:", err)
			http.Error(w, "Ошибка работы с БД", http.StatusInternalServerError)
		}
		return
	}

	var newOrderID int
	createOrderQuery := `INSERT INTO Orders (user_id, status) OUTPUT INSERTED.Id VALUES (@p1, 'created')`
	err = tx.QueryRow(createOrderQuery, order.UserID).Scan(&newOrderID)
	if err != nil {
		fmt.Println("Ошибка создания заказа:", err)
		http.Error(w, "Ошибка оформления", http.StatusInternalServerError)
		return
	}

	for _, item := range order.Items {

		var currentQuantity int
		checkQueryProduct := `SELECT Quantity FROM Products WHERE Id = @p1`
		err = tx.QueryRow(checkQueryProduct, item.ProductID).Scan(&currentQuantity)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, fmt.Sprintf("Товар с Id %d не найден", item.ProductID), http.StatusNotFound)
			} else {
				fmt.Println("Ошибка поиска товара:", err)
				http.Error(w, "Ошибка работы с БД", http.StatusInternalServerError)
			}
			return
		}
		if currentQuantity < item.Quantity {
			http.Error(w, fmt.Sprintf("Недостаточно товара с Id %d на складе. В наличии: %d", item.ProductID, currentQuantity), http.StatusBadRequest)
			return
		}
		updateProductQuery := `UPDATE Products SET Quantity = Quantity - @p1 WHERE Id = @p2`
		_, err = tx.Exec(updateProductQuery, item.Quantity, item.ProductID)
		if err != nil {
			fmt.Println("Ошибка списания товара:", err)
			http.Error(w, "Ошибка обновления склада", http.StatusInternalServerError)
			return
		}
		insertItemQuery := `INSERT INTO Order_item (order_id, product_id, quantity) VALUES (@p1, @p2, @p3)`
		_, err = tx.Exec(insertItemQuery, newOrderID, item.ProductID, item.Quantity)
		if err != nil {
			fmt.Println("Ошибка добавления товара в чек:", err)
			http.Error(w, "Ошибка формирования корзины", http.StatusInternalServerError)
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		fmt.Println("Ошибка Commit:", err)
		http.Error(w, "Ошибка сохранения заказа", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Заказ успешно оформлен! Номер вашего заказа: %d", newOrderID)
}
func UpdateStatus(w http.ResponseWriter, r *http.Request) {
	newstatus := r.PathValue("status")
	idtext := r.PathValue("id")

	if newstatus == "" {
		http.Error(w, "Статус не может быть пустым", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idtext)
	if err != nil {
		fmt.Println("Ошибка конвертации:", err)
		http.Error(w, "Неверный формат Id", http.StatusBadRequest)
		return
	}
	Query := `UPDATE Orders SET Status = @p1 WHERE Id = @p2`
	result, err := database.Db.Exec(Query, newstatus, id)
	if err != nil {
		fmt.Println("Ошибка обновления в БД:", err)
		http.Error(w, "Ошибка на сервере", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Заказ с таким Id не найден", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Статус заказа успешно обновлен")
}
