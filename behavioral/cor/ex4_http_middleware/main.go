package main

import "net/http"

func RequireAPIKey(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		apiSecret := r.Header.Get("X-API-Key")
		if apiSecret != "secret123"{
			http.Error(w,"Unauthorized",http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w,r)
	})
}
