package sysconfig

import (
	"errors"
	"k8s.io/klog/v2"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// InitProxy
var InitProxy *httputil.ReverseProxy

// NewProxy . 按照配置生成新的反向代理
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
		klog.Info("request url host: ", resp.Request.URL.Host)

		// FIXME: 主要处理重定向的场景
		r, ok := InitProxyMap[resp.Request.URL.Host]
		if !ok {
			return errors.New("not found InitProxy in InitProxyMap")
		}
		InitProxy = r
		klog.Info("resp code: ", resp.StatusCode)

		return nil
	}

}

func modifyRequest(req *http.Request) {

}

// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler(proxys map[string]*httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	klog.Info("use proxy gateway !!")
	return func(w http.ResponseWriter, req *http.Request) {

		// 处理 map panic 的问题
		defer func() {
			if err := recover(); err != nil {
				log.Println("Recovered from panic:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		// 这里需要分割url
		s := strings.Split(req.URL.Path, "/")
		res1, res2 := handler(s)
		klog.Infof("prefix: [%v], service url: [%v]", res1, res2)

		// FIXME: 这里会有重定向的问题，重定向的请求不能执行
		// FIXME: 目前是把没有在map中找到的请求都直接return，长期会有问题
		res, ok1 := proxys[res1]
		if !ok1 {
			klog.Error("proxy map not found")
			InitProxy.ServeHTTP(w, req)
		}

		r, ok2 := HostMap[res1]
		req.URL.Host = r
		if !ok2 {
			klog.Error("HostMap map not found")
			InitProxy.ServeHTTP(w, req)

		}
		req.URL.Path = res2
		klog.Infof("request URL.Host: [%v], URL.Path: [%v]", req.URL.Host, req.URL.Path)
		if ok1 && ok2 {
			res.ServeHTTP(w, req)
		}

	}
}

// handler 临时方案，处理 proxy 的 url
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
