# Go-Zero + Polaris 服务发现示例

本项目展示了如何使用 go-zero 微服务框架结合 Polaris 服务发现来构建分布式服务。项目包含一个 gRPC 服务端和一个客户端，演示了服务注册、发现和调用的完整流程。

## 项目结构

```
.
├── server/                 # gRPC 服务端
│   ├── server.go          # 服务端主程序
│   ├── go.mod             # Go 模块依赖
│   ├── go.sum             # 依赖锁定文件
│   ├── polaris.yaml       # Polaris 配置文件
│   ├── etc/
│   │   └── polaris.yaml   # 服务端配置
│   └── internal/
│       └── config/
│           └── config.go  # 配置结构定义
├── client/                 # gRPC 客户端
│   ├── client.go          # 客户端主程序
│   ├── go.mod             # Go 模块依赖
│   ├── go.sum             # 依赖锁定文件
│   ├── polaris.yaml       # Polaris 配置文件
│   └── etc/
│       └── polaris.yaml   # 客户端配置
└── README.md              # 项目说明文档
```

## 环境要求

- **Go 版本**: 1.18 或更高版本
- **Polaris 服务端**: 需要运行 Polaris 服务发现服务
- **操作系统**: 支持 Linux、macOS、Windows

## 依赖组件

### 主要依赖
- [go-zero](https://github.com/zeromicro/go-zero): 微服务框架
- [zero-contrib](https://github.com/zeromicro/zero-contrib): go-zero 扩展组件
- [Polaris](https://polarismesh.cn/): 服务发现和治理平台

### 核心功能
- **服务注册**: 服务端自动注册到 Polaris
- **服务发现**: 客户端通过 Polaris 发现服务
- **负载均衡**: 自动实现客户端负载均衡
- **健康检查**: 服务健康状态监控

## 快速开始

### 1. 启动 Polaris 服务

首先需要启动 Polaris 服务发现服务。可以通过以下方式之一：

#### 使用 Docker (推荐)
```bash
# 拉取并启动 Polaris 服务
docker run -d --name polaris -p 8091:8091 -p 8090:8090 polarismesh/polaris-server:latest
```

#### 使用 Docker Compose
```bash
# 创建 docker-compose.yml 文件
cat > docker-compose.yml << EOF
version: '3.8'
services:
  polaris:
    image: polarismesh/polaris-server:latest
    ports:
      - "8091:8091"
      - "8090:8090"
    environment:
      - POLARIS_LOG_LEVEL=info
EOF

# 启动服务
docker-compose up -d
```

### 2. 构建和运行服务端

```bash
# 进入服务端目录
cd server

# 构建服务端
go build -mod=vendor -o server server.go

# 运行服务端
./server -f ./etc/polaris.yaml
```

服务端启动后会：
- 监听 `0.0.0.0:3453` 端口
- 自动注册到 Polaris 服务发现中心
- 服务名为 `EchoServerZero`

### 3. 构建和运行客户端

```bash
# 新开一个终端，进入客户端目录
cd client

# 构建客户端
go build -mod=vendor -o client client.go

# 运行客户端
./client -f ./etc/polaris.yaml
```

客户端启动后会：
- 通过 Polaris 发现 `EchoServerZero` 服务
- 每秒向服务端发送一次请求
- 打印服务端响应结果

## 配置说明

### 服务端配置 (`server/etc/polaris.yaml`)
```yaml
Name: server              # 服务名称
ListenOn: 0.0.0.0:3453   # 监听地址和端口
```

### 客户端配置 (`client/etc/polaris.yaml`)
```yaml
Name: client                                                    # 客户端名称
Target: polaris://127.0.0.1:38091/EchoServerZero?timeout=5000s # 目标服务地址
```

### Polaris 连接配置 (`polaris.yaml`)
```yaml
global:
  serverConnector:
    addresses:
      - 127.0.0.1:8091    # Polaris 服务地址
```

## 测试指南

### 1. 基本功能测试

启动服务端和客户端后，客户端应该能够成功调用服务端：

```bash
# 客户端输出示例
15:04:05 => hello from hostname-xxx
15:04:06 => hello from hostname-xxx
15:04:07 => hello from hostname-xxx
```

### 2. 多实例测试

启动多个服务端实例来测试负载均衡：

```bash
# 终端1: 启动第一个服务端实例
cd server
./server -f ./etc/polaris.yaml

# 终端2: 修改端口启动第二个实例
cd server
# 临时修改配置文件中的端口为 3454
./server -f ./etc/polaris.yaml

# 终端3: 启动客户端
cd client
./client -f ./etc/polaris.yaml
```

观察客户端输出，应该能看到请求在不同的服务端实例之间轮询。

### 3. 服务发现测试

测试服务动态上下线：

1. 启动客户端和服务端
2. 停止服务端进程
3. 观察客户端是否能检测到服务不可用
4. 重新启动服务端
5. 观察客户端是否能自动重连

### 4. 性能测试

可以使用以下命令进行简单的性能测试：

```bash
# 使用 go test 进行基准测试
cd server
go test -bench=. -benchmem

# 或者使用 hey 工具进行压力测试
# 首先安装 hey: go install github.com/rakyll/hey@latest
hey -n 1000 -c 10 -m POST -H "Content-Type: application/grpc" http://localhost:3453
```

## 故障排除

### 常见问题

1. **连接 Polaris 失败**
   ```
   错误信息: connection refused
   解决方案: 确保 Polaris 服务正在运行，检查端口 8091 是否可访问
   ```

2. **服务注册失败**
   ```
   错误信息: register service failed
   解决方案: 检查 polaris.yaml 配置文件中的服务地址是否正确
   ```

3. **客户端发现服务失败**
   ```
   错误信息: service not found
   解决方案: 确保服务端已启动并成功注册，检查服务名称是否匹配
   ```

4. **端口冲突**
   ```
   错误信息: bind: address already in use
   解决方案: 修改配置文件中的端口号或停止占用端口的进程
   ```

### 调试技巧

1. **查看详细日志**
   ```bash
   # 启动时添加详细日志
   ./server -f ./etc/polaris.yaml -v
   ```

2. **检查 Polaris 控制台**
   - 访问 `http://localhost:8090` 查看 Polaris 管理界面
   - 检查服务注册状态和健康检查结果

3. **网络连接测试**
   ```bash
   # 测试 Polaris 连接
   telnet 127.0.0.1 8091
   
   # 测试服务端连接
   telnet 127.0.0.1 3453
   ```

## 开发指南

### 添加新的 gRPC 方法

1. 修改 proto 文件定义
2. 重新生成 Go 代码
3. 在服务端实现新方法
4. 在客户端添加调用代码

### 自定义配置

可以通过修改配置文件来自定义：
- 服务监听端口
- Polaris 服务地址
- 超时时间
- 重试策略

### 扩展功能

项目支持以下扩展：
- 添加中间件 (认证、限流、监控)
- 集成链路追踪
- 添加指标监控
- 实现熔断降级

## 相关链接

- [go-zero 官方文档](https://go-zero.dev/)
- [Polaris 官方文档](https://polarismesh.cn/docs/)
- [gRPC Go 教程](https://grpc.io/docs/languages/go/)
- [项目源码](https://github.com/polarismesh/examples)

## 许可证

本项目采用 MIT 许可证，详情请参见 LICENSE 文件。 