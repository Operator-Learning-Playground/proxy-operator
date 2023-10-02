package controller

import (
	"context"
	proxyv1alpha1 "github.com/myoperator/proxyoperator/pkg/apis/proxy/v1alpha1"
	"github.com/myoperator/proxyoperator/pkg/sysconfig"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ProxyController struct {
	client.Client
}

func NewProxyController() *ProxyController {
	return &ProxyController{}
}

// Reconcile 调协loop
func (r *ProxyController) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {

	proxy := &proxyv1alpha1.Proxy{}
	err := r.Get(ctx, req.NamespacedName, proxy)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			klog.Error("get proxy error: ", err)
			return reconcile.Result{}, err
		}
		// 如果未找到的错误，不再进入调协
		return reconcile.Result{}, nil
	}
	klog.Info(proxy)

	// 修改 proxy 配置
	err = sysconfig.AppConfig(proxy)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// InjectClient 使用controller-runtime 需要注入的client
func (r *ProxyController) InjectClient(c client.Client) error {
	r.Client = c
	return nil
}

// TODO: 删除逻辑并未处理
