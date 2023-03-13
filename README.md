## proxy-operator 简易型网关控制器

### 项目思路与设计
![](https://github.com/googs1025/proxy-operator/blob/main/image/%E6%B5%81%E7%A8%8B%E5%9B%BE.jpg?raw=true)
设计背景：集群的网关通常都是采用nginx-controller部署的方式，对自用小集群难免存在部署步骤复杂等问题。
本项目在此问题上，基于k8s的扩展功能，实现proxy的自定义资源，做出一个有反向代理功能的controller应用。调用方可在cluster中部署与启动相关配置即可使用。
思路：当应用启动后，会启动一个controller与proxy反向代理服务，controller会监听crd资源，并执行相应的业务逻辑。

### 项目功能
1. 支持多个service服务的配置，实现集群的网关功能。
2. 实现中间件框架，可以自定义中间件。
3. 提供ip限流、query限流功能。
### 本地调适


### 项目部署


### RoadMap
