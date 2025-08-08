# 介绍
## consumer
使用 polaris-sidecar 的 dnsagent 模式进行服务发现，作为服务主调方，服务本身未注册到北极星
服务发现时使用格式为<service>.<namespace>的域名进行访问。

## provider
使用polaris-go sdk 注册到北极星，作为服务被调方

# 部署流程
## provider部署
- 在 provider 目录下执行
### 编译
```shell
make build
```
### 执行
```shell
# 设置北极星地址
export POLARIS_SERVER=127.0.0.1:8091
# 启动 provider
./provider
```

## consumer部署
### 前提条件
按照安装指引安装 polaris-sidecar组件，[安装指引链接](https://github.com/polarismesh/polaris-sidecar/blob/main/README-zh.md)
- 以下在 consumer 目录下执行
### 编译
```shell
make build
```
### 执行
```shell
./consumer
```

# 验证
## consumer 调用 provider
```shell
curl localhost:20000/echo
```