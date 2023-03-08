package sysconfig

import (
	"errors"
	"k8s.io/klog/v2"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// InitProxy
var InitProxy *httputil.ReverseProxy

// NewProxy .
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		modifyRequest(req)
	}

	proxy.ModifyResponse = modifyResponse()

	return proxy, nil
}

func modifyResponse() func(response *http.Response) error {
	return func(resp *http.Response) error {
		klog.Info("")
		klog.Info(resp.Request.URL.Host)

		r, ok := InitProxyMap[resp.Request.URL.Host]
		if !ok {
			return errors.New("not found InitProxy in InitProxyMap")
		}
		InitProxy = r
		klog.Info(resp.StatusCode)

		return nil
	}

}

func modifyRequest(req *http.Request)  {

}

// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler(proxys map[string]*httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	klog.Info("use proxy!!")
	return func(w http.ResponseWriter, req *http.Request) {



		// 这里需要分割url
		s := strings.Split(req.URL.Path, "/")
		res1, res2 := handler(s)
		klog.Info("res1: ", res1, " res2: ", res2)

		// FIXME: 这里会有重定向的问题，重定向的请求不能执行
		// FIXME: 目前是把没有在map中找到的请求都直接return，长期会有问题
		res, ok1 := proxys[res1]
		if !ok1 {
			klog.Error("proxy map 中没有找到")
			InitProxy.ServeHTTP(w, req)

		}

		r, ok2 := HostMap[res1]
		req.URL.Host = r
		if !ok2 {
			klog.Error("HostMap map 中没有找到")
			InitProxy.ServeHTTP(w, req)

		}
		req.URL.Path = res2
		klog.Info("request: ", req.URL.Host, req.URL.Path)
		if ok1 && ok2 {
			res.ServeHTTP(w, req)
		}

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


