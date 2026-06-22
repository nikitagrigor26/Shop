package middleware

import (
	"fmt"
	"net/http"
)

func LogingMiddelware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Получен запрос :\n Метод: %s\n Путь: %s\n", r.Method, r.URL)
		next.ServeHTTP(w, r)
		fmt.Println("---Все обработано---")
	})
}
