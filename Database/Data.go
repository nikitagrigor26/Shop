package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

var Db *sql.DB

func InitDb() {
	connectionString := "server = localhost; port =1433; database = ShopDB; Trusted=True;"
	var err error
	Db, err = sql.Open("sqlserver", connectionString)
	if err != nil {
		log.Fatal("Ошибка не получилось открыть базу", err)

	}
	err = Db.Ping()
	if err != nil {
		fmt.Println("Ошибка не удалось открыть БД")
		return
	}

	productQuery := `IF NOT EXISTS (SELECT * FROM sysobjects WHERE name = 'products' and xtype='U') 
   BEGIN 
   			CREATE TABLE Products(
					Id int NOT NULL IDENTITY(1,1) PRIMARY KEY,
					Name nvarchar(255) NOT NULL,
                    Price decimal(10,2) NOT NULL,
					Description varchar(255) NOT NULL,
					Quantity int NOT NULL)
   END`

	_, err = Db.Exec(productQuery)
	if err != nil {
		log.Fatal("Ошибка подключения к таблице products", err)

	}

	userQuery := `IF NOT EXISTS (SELECT * FROM sysobjects WHERE name = 'users' and xtype='U') 
	BEGIN 
		create table Users(
				Id int NOT NULL IDENTITY(1,1) PRIMARY KEY,
				Name nvarchar(255) NOT NULL,
				Email nvarchar(255) NOT NULL,
				Password nvarchar(255) NOT NULL)
	END`
	_, err = Db.Exec(userQuery)
	if err != nil {
		log.Fatal("Ошибка с таблицей User", err)
	}
	orderQuery := `IF NOT EXISTS (SELECT * FROM sysobjects WHERE name = 'orders' and xtype='U')
    BEGIN
    	create table Orders(
				Id int NOT NULL IDENTITY(1,1) PRIMARY KEY,
				user_id int FOREIGN KEY REFERENCES Users(Id),
				Status nvarchar(255) NOT NULL)
	END`
	_, err = Db.Exec(orderQuery)
	if err != nil {
		log.Fatal("Ошибка с табллицей Order", err)

	}
	orderItemQuery := `IF NOT EXISTS (SELECT * FROM sysobjects WHERE name = 'order_item' and xtype='U')
	BEGIN
		CREATE TABLE Order_item (
			Id int NOT NULL IDENTITY(1,1) PRIMARY KEY,
			product_id int FOREIGN KEY REFERENCES Products(Id),
			order_id int FOREIGN KEY REFERENCES Orders(Id),
			quantity int NOT NULL)
	END`
	_, err = Db.Exec(orderItemQuery)
	if err != nil {
		log.Fatal("Ошибка с таблицей Order_item", err)

	}
	fmt.Println("Бд подключилась корректно")
}
