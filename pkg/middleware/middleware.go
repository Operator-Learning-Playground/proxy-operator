package middleware

import (
	"log"
	"net/http"
)

type Middleware func(handlerFunc http.HandlerFunc) http.HandlerFunc

func ApplyMiddleware(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

// LoggerMiddleware 日志中间件
func LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next(w, r)
	}
}
