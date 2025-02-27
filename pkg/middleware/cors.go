package middleware

import (
	"fmt"
	"net/http"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		fmt.Println(origin)
		if origin != "http://localhost:3000" {
			next.ServeHTTP(w, r) // Если origin не совпадает, передаем дальше без CORS
			return
		}

		header := w.Header()
		header.Set("Access-Control-Allow-Origin", origin) // Разрешаем только этот Origin
		header.Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			header.Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, HEAD, PATCH")
			header.Set("Access-Control-Allow-Headers", "authorization, content-type, content-length")
			header.Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusNoContent) // Прерываем обработку OPTIONS
			return
		}

		next.ServeHTTP(w, r)
	})
}