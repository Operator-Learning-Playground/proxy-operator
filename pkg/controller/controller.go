package controller

import (
	"context"
	proxyv1alpha1 "github.com/myoperator/proxyoperator/pkg/apis/proxy/v1alpha1"
	"github.com/myoperator/proxyoperator/pkg/sysconfig"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

type ProxyController struct {
	client.Client
}

func NewProxyController() *ProxyController {
	return &ProxyController{}
}

// Reconcile 调协loop
func (r *ProxyController) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {

	// 1. 获取资源对象
	proxy := &proxyv1alpha1.Proxy{}
	err := r.Get(ctx, req.NamespacedName, proxy)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			klog.Error("get proxy error: ", err)
			return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 60}, err
		}
		// 如果未找到的错误，不再进入调协
		return reconcile.Result{}, nil
	}
	klog.Info(proxy)

	// 2. 是否是删除流程
	if !proxy.DeletionTimestamp.IsZero() {
		klog.Info("clean proxy config")
		err := sysconfig.CleanConfig()
		if err != nil {
			klog.Error("clean proxy config error: ", err)
			return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 60}, err
		}

		// 清理完成后，从 Finalizers 中移除 Finalizer
		controllerutil.RemoveFinalizer(proxy, finalizerName)
		err = r.Update(ctx, proxy)
		if err != nil {
			klog.Error("clean proxy finalizer err: ", err)
			return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 60}, err
		}

		klog.Info("successful delete reconcile")

		return reconcile.Result{}, nil
	}

	// 3. 检查是否已添加 Finalizer
	if !containsFinalizer(proxy) {
		// 添加 Finalizer
		controllerutil.AddFinalizer(proxy, finalizerName)
		err = r.Update(ctx, proxy)
		if err != nil {
			klog.Error("update proxy finalizer err: ", err)
			return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 60}, err
		}
	}

	// 4. 修改 proxy 配置
	err = sysconfig.AppConfig(proxy)
	if err != nil {
		klog.Error("apply proxy config error: ", err)
		return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 60}, err
	}
	klog.Info("successful reconcile")
	return reconcile.Result{}, nil
}

const (
	finalizerName = "api.practice.com/finalizer"
)

// InjectClient 使用controller-runtime 需要注入的client
func (r *ProxyController) InjectClient(c client.Client) error {
	r.Client = c
	return nil
}

func containsFinalizer(proxy *proxyv1alpha1.Proxy) bool {
	for _, finalizer := range proxy.Finalizers {
		if finalizer == finalizerName {
			return true
		}
	}
	return false
}

