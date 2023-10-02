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


### RoadMap
