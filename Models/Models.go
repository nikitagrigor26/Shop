package Models

import (
	database "Shop/Database"
	"fmt"
	"time"
)

type Product struct {
	Id          int
	Name        string
	Price       float64
	Description string
	Quantity    int
}
type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

type Order struct {
	Id      int
	user_id int
	Status  string
}
type Order_item struct {
	Id         int
	product_id int
	order_id   int
	quantity   int
}

func StatsProduct() {
	for {
		time.Sleep(1 * time.Minute)

		var count int
		err := database.Db.QueryRow("SELECT COUNT(*) FROM Products").Scan(&count)
		if err != nil {
			fmt.Println("[Статистика] Ошибка подсчета продуктов :", err)
			continue
		}

		fmt.Printf("[Статистика] Прямо сейчас в базе хранится продуктов: %d шт.\n", count)
	}
}
func StatsUsers() {
	for {
		time.Sleep(1 * time.Minute)

		var count int
		err := database.Db.QueryRow("SELECT COUNT(*) FROM Users").Scan(&count)
		if err != nil {
			fmt.Println("[Статистика] Ошибка подсчета пользователей :", err)
			continue
		}

		fmt.Printf("[Статистика] Прямо сейчас в базе хранится пользователей: %d шт.\n", count)
	}
}
