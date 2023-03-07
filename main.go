package main

import (
	"fmt"
	"github.com/myoperator/proxyoperator/pkg/controller"
	"github.com/myoperator/proxyoperator/pkg/k8sconfig"
	"github.com/myoperator/proxyoperator/pkg/sysconfig"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/klog/v2"
	"log"
	"net/http"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

/*
	manager 主要用来管理Controller Admission Webhook 包括：
	访问资源对象的client cache scheme 并提供依赖注入机制 优雅关闭机制

	operator = crd + controller + webhook
*/



func main() {

	// 1. 管理器初始化
	mgr, err := manager.New(k8sconfig.K8sRestConfig(), manager.Options{})
	if err  != nil {
		klog.Error(err, "unable to set up manager")
		os.Exit(1)
	}

	// 2. 控制器相关
	//proxyCtl := controller.NewProxyController()
	// 传入资源&v1.Ingress{}，也可以用crd
	err = builder.ControllerManagedBy(mgr).
		For(&networkingv1.Ingress{}).
		//Watches(&source.Kind{ // 加入监听。
		//	Type: &networkingv1.Ingress{},
		//}, handler.Funcs{
		//	DeleteFunc: proxyCtl.IngressDeleteHandler,
		//}).
		Complete(controller.NewProxyController())

	// 3. ++ 注册进入序列化表
	err = k8sconfig.SchemeBuilder.AddToScheme(mgr.GetScheme())
	if err != nil {
		klog.Error(err, "unable add schema")
		os.Exit(1)
	}

	// 4. 载入业务配置
	if err = sysconfig.InitConfig(); err != nil {
		klog.Error(err, "unable to load sysconfig")
		os.Exit(1)
	}
	errC := make(chan error)

	// 3. 启动controller管理器
	go func() {
		klog.Info("controller start!! ")
		if err = mgr.Start(signals.SetupSignalHandler()); err != nil {
			errC <-err
		}
	}()

	// 4. 启动网关
	go func() {
		klog.Info("proxy start!! ")
		http.HandleFunc("/", sysconfig.ProxyRequestHandler(sysconfig.ProxyMap))
		if err = http.ListenAndServe(fmt.Sprintf(":%d", sysconfig.SysConfig1.Server.Port), nil); err != nil {
			errC <-err
		}
	}()

	// 这里会阻塞，两种常驻进程可以使用这个方法
	getError := <-errC
	log.Println(getError.Error())

}


