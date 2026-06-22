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

func ViewUsersHandler(w http.ResponseWriter, r *http.Request) {
	namefilter := r.URL.Query().Get("name")
	query := `SELECT * FROM Users WHERE 1=1`
	var args []any
	if namefilter != "" {
		args = []any{namefilter}
		query += fmt.Sprintf(" AND Name = %s", namefilter)
	}
	rows, err := database.Db.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "Ошибка с запросом БД к пользователю")
		return
	}
	defer rows.Close()
	var users []Models.User
	for rows.Next() {
		var user Models.User
		err = rows.Scan(&user.Id, &user.Name, user.Email, user.Password)
		if err != nil {
			fmt.Println(err)
			continue
		}
		users = append(users, user)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(users)
}
func ViewUserhandler(w http.ResponseWriter, r *http.Request) {
	idtext := r.URL.Query().Get("id")
	var id int
	var err error
	if idtext != "" {
		id, err = strconv.Atoi(idtext)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "Ошибка в запросе, не удалось обработать Id")
			return
		}
	}
	var user Models.User
	query := `SELECT Name,Email,Password FROM Users WHERE Id = @p1`
	err = database.Db.QueryRow(query, id).Scan(user.Name, user.Password, user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprintf(w, "Ошибка, пользователя с таким Id не существует")
			return
		}
		fmt.Println(err)
		fmt.Fprintf(w, "Такого пользователя нету в баще")
		return
	}
}
