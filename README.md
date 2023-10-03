## proxy-operator 简易型网关控制器

### 项目思路与设计
![](https://github.com/googs1025/proxy-operator/blob/main/image/%E6%B5%81%E7%A8%8B%E5%9B%BE.jpg?raw=true)
设计背景：集群的网关通常都是采用 nginx-controller 部署的方式，对自用小集群难免存在部署步骤复杂等问题。
本项目在此问题上，基于 k8s 的扩展功能，实现 proxy 的自定义资源，做出一个有反向代理功能的 controller 应用。调用方可在 cluster 中部署与启动相关配置即可使用。

思路：当应用启动后，会启动一个 **controller** 与 **proxy** 反向代理服务，controller 会监听 crd 资源，并执行相应的业务逻辑。

### 项目功能
1. 支持多个 service 服务的配置，实现集群的网关功能。
2. 实现中间件框架，可以自定义中间件。
3. 提供 ip 限流、query 限流功能。

- cr 资源对象如下所示：
```yaml
apiVersion: api.practice.com/v1alpha1
kind: Proxy
metadata:
  name: myproxy
spec:
  server:                               # 反向代理 server ip:port
    ip: localhost
    port: 10086
  rules:                                # 后端路由列表
    - path:
        backend:                        # 后端
          prefix: /service1             # 自定义前缀
          url: http://localhost:8899    # 真正要访问的服务
    - path:
        backend:
          prefix: /service2
          url: http://localhost:8800
```

### 项目部署
1. 打镜像。
```bash
# 项目根目录执行
[root@VM-0-16-centos proxyoperator]# pwd
/root/proxyoperator
# 可以直接使用 docker 镜像部署
[root@VM-0-16-centos proxyoperator]# docker build -t myproxyoperator:v1 .
Sending build context to Docker daemon  49.16MB
Step 1/19 : FROM golang:1.18.7-alpine3.15 as builder
 ---> 33c97f935029
Step 2/19 : WORKDIR /app
 ---> Using cache ...
```   
2. apply crd 资源
```bash
[root@VM-0-16-centos yaml]# ls
deploy_docker.yaml  deploy.yaml  example.yaml  proxy.yaml  rbac.yaml
[root@VM-0-16-centos yaml]# kubectl apply -f proxy.yaml
customresourcedefinition.apiextensions.k8s.io/proxys.api.practice.com unchanged
```   
3. 启动 controller 服务(需要先执行 rbac.yaml，否则服务会报错)
```bash
[root@VM-0-16-centos yaml]# kubectl apply -f rbac.yaml
serviceaccount/myproxy-sa unchanged
clusterrole.rbac.authorization.k8s.io/myproxy-clusterrole unchanged
clusterrolebinding.rbac.authorization.k8s.io/myproxy-ClusterRoleBinding unchanged
[root@VM-0-16-centos yaml]# kubectl apply -f deploy_docker.yaml
deployment.apps/myproxy-controller unchanged
service/myproxy-svc unchanged
```   
4. 查看 operator 服务

```bash
[root@VM-0-16-centos yaml]# kubectl get pods | grep proxy
myproxy-controller-789c6f7c66-jfdmj                1/1     Running            0          13h
```

### RoadMap
