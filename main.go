package main

import (
	"fmt"
	proxyv1alpha1 "github.com/myoperator/proxyoperator/pkg/apis/proxy/v1alpha1"
	"github.com/myoperator/proxyoperator/pkg/controller"
	"github.com/myoperator/proxyoperator/pkg/k8sconfig"
	"github.com/myoperator/proxyoperator/pkg/middleware"
	"github.com/myoperator/proxyoperator/pkg/sysconfig"
	_ "k8s.io/code-generator"
	"k8s.io/klog/v2"
	"log"

	"net/http"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

/*
	manager 主要用来管理Controller Admission Webhook 包括：
	访问资源对象的client cache scheme 并提供依赖注入机制 优雅关闭机制

	operator = crd + controller + webhook
*/

func main() {

	logf.SetLogger(zap.New())
	// 1. 管理器初始化
	mgr, err := manager.New(k8sconfig.K8sRestConfig(), manager.Options{
		Logger: logf.Log.WithName("proxy-operator"),
	})
	if err != nil {
		mgr.GetLogger().Error(err, "unable to set up manager")
		os.Exit(1)
	}

	// 2. ++ 注册进入序列化表
	err = proxyv1alpha1.SchemeBuilder.AddToScheme(mgr.GetScheme())
	if err != nil {
		klog.Error(err, "unable add schema")
		os.Exit(1)
	}

	// 3. 控制器相关
	proxyCtl := controller.NewProxyController()

	err = builder.ControllerManagedBy(mgr).
		For(&proxyv1alpha1.Proxy{}).
		Complete(proxyCtl)

	// 4. 载入业务配置
	if err = sysconfig.InitConfig(); err != nil {
		klog.Error(err, "unable to load sysconfig")
		os.Exit(1)
	}
	errC := make(chan error)

	// 5. 启动controller管理器
	go func() {
		klog.Info("controller start!! ")
		if err = mgr.Start(signals.SetupSignalHandler()); err != nil {
			errC <- err
		}
	}()

	// 6. 启动网关
	go func() {
		klog.Info("proxy start!! ")
		// 中间件
		http.HandleFunc("/", middleware.ApplyMiddleware(sysconfig.ProxyRequestHandler(sysconfig.ProxyMap), middleware.LoggerMiddleware,
			middleware.IpLimiterMiddleware, middleware.ParamLimiterMiddleware, middleware.PanicMiddleware))
		if err = http.ListenAndServe(fmt.Sprintf(":%d", sysconfig.SysConfig1.Server.Port), nil); err != nil {
			errC <- err
		}
	}()

	// 会阻塞，两种常驻进程可以使用这个方法
	getError := <-errC
	log.Println(getError.Error())

}
