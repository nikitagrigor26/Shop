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
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
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
	idtext := r.PathValue("id")
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
	err = database.Db.QueryRow(query, id).Scan(&user.Name, &user.Password, &user.Email)
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
func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	var user Models.User
	json.NewDecoder(r.Body).Decode(&user)
	Query := `INSERT INTO Users (Name,Email,Password) VALUES (@p1,@p2,@p3)`
	_, err := database.Db.Exec(Query, user.Name, user.Email, user.Password)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "Ошибка с БД")
		return
	}
	fmt.Fprintf(w, "Пользовать добавлена")
}
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	idtext := r.PathValue("id")
	var id int
	var err error
	if idtext != "" {
		id, err = strconv.Atoi(idtext)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "Ошибка конвертации Id, проверьте запрос и повторите попытку")
			return
		}
	}
	query := `DELETE FROM Users WHERE Id = @p1`
	res, err := database.Db.Exec(query, id)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "Ошибка БД упала :", err)
		return
	}
	result, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	if result == 0 {
		fmt.Fprintf(w, "Пользователя не существует, удаление невозможно")
		return
	}
	fmt.Fprintf(w, "Пользователь удален")
}
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	idtext := r.PathValue("id")
	var id int
	var err error
	if idtext != "" {
		id, err = strconv.Atoi(idtext)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "Ошибка конвертации Id, проверьте запрос и попробуйте снова")
			return
		}
	}
	var user Models.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println(err)
		return
	}
	Query := `UPDATE Users SET Name = @p2,Email = @p3,Password = @p4 WHERE Id = @p1`
	res, errs := database.Db.Exec(Query, id, user.Name, user.Email, user.Password)
	if errs != nil {
		fmt.Println(errs)
		return
	}
	RowsAffected, errs := res.RowsAffected()
	if errs != nil {
		fmt.Println(errs)
	}
	if RowsAffected == 0 {
		fmt.Fprintf(w, "Такого пользователя не существует")
		return
	}
	fmt.Fprintf(w, "пользователь обновлен")
}
