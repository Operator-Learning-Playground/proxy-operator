package controller

import (
	"context"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	//"github.com/myoperator/proxyoperator/pkg/sysconfig"
)

const (
	ProxyControllerAnnotation = "myproxy"
	ingressAnnotationKey = "kubernetes.io/ingress.class"
)

type ProxyController struct {
	client.Client
}

func NewProxyController() *ProxyController {
	return &ProxyController{}
}

// Reconcile 调协loop
func (r *ProxyController) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {

	klog.Info(req.NamespacedName)

	return reconcile.Result{}, nil
}

// 使用controller-runtime 需要注入的client
func(r *ProxyController) InjectClient(c client.Client) error {
	r.Client = c
	return nil
}



