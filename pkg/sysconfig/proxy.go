package sysconfig

import (
	"fmt"
	"k8s.io/klog/v2"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)


// NewProxy .
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	return proxy, nil
}

// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler(proxys map[string]*httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	fmt.Println("use proxy!!")
	return func(w http.ResponseWriter, req *http.Request) {



		// 这里需要分割url
		s := strings.Split(req.URL.Path, "/")
		res1, res2 := handler(s)
		klog.Info("res1: ", res1, " res2: ", res2)

		// FIXME: 这里会有重定向的问题，重定向的请求不能执行
		// FIXME: 目前是把没有在map中找到的请求都直接return，长期会有问题
		res, ok := proxys[res1]
		if !ok {
			return
		}
		req.URL.Host, ok = HostMap[res1]
		if !ok {
			return
		}
		req.URL.Path = res2
		klog.Info("request: ", req.URL.Host, req.URL.Path)
		res.ServeHTTP(w, req)

	}
}

func handler(url []string) (string, string) {
	res1 := "/" + url[1]
	url = url[2:]
	res2 := ""
	for _, v := range url {
		res2 += "/"
		res2 += v
	}

	return res1, res2
}


