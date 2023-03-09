package middleware

import (
	"github.com/myoperator/proxyoperator/pkg/middleware/limit"
	"k8s.io/klog/v2"
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

//
func LimiterMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ip := req.RemoteAddr
		klog.Info("ip: ", ip)

		var limiter *limit.Bucket

		if v, ok := limit.IpCache.Data.Load(ip); ok {
			limiter = v.(*limit.Bucket)
		} else {
			limiter = limit.NewBucket(limit.DefaultCap, limit.DefaultRate)
			limit.IpCache.Data.Store(ip, limiter)
		}

		// 如果限流器接受，则走到下一个中间件，不然就报错
		if limiter.IsAccept() {
			next(w, req)
		} else {
			w.Write([]byte("this ip is too many request!!"))
		}
	}
}

