# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此代码库中工作时提供指导。

## 项目概述

这是一个用 Go 语言编写的分布式蜜罐诱捕系统。它通过创建诱饵网络服务来检测和捕获恶意活动。系统由 5 个主要组件组成，通过 gRPC 和 RabbitMQ 进行通信。

## 系统架构

### 核心组件

1. **honey_server** (`apps/honey_server/`) - 管理服务器
   - 蜜罐基础设施的中央编排器
   - 提供 REST API（端口 8079）和 gRPC 服务器（端口 8081）
   - 管理节点注册、IP 分配、端口转发
   - 使用技术：Gin、gRPC、GORM (MySQL)、Redis、RabbitMQ、JWT

2. **honey_node** (`apps/honey_node/`) - 节点代理
   - 部署在目标主机上，创建虚拟 IP 和诱饵服务
   - 集成 Suricata IDS 进行流量监控
   - 执行网络操作（IP 创建、端口隧道）
   - 使用技术：gRPC、RabbitMQ、SQLite、arping、tail 日志监控

3. **alert_server** (`apps/alert_server/`) - 告警处理
   - 从 RabbitMQ 队列消费告警信息
   - 使用地理位置（ip2region）和服务映射丰富数据
   - 将告警存储到 Elasticsearch（端口 8078）
   - 使用技术：Gin、Elasticsearch (olivere/elastic)、RabbitMQ、Redis

4. **image_server** (`apps/image_server/`) - 容器管理
   - 编排 Docker 容器用于蜜罐服务
   - 管理镜像、虚拟网络、主机模板（端口 8080）
   - 使用技术：Gin、Docker SDK、GORM、Redis

5. **goroutine_test** (`apps/goroutine_test/`) - 测试代码（非生产环境）

### 数据流

```
攻击者 → 蜜罐 IP:端口 → Suricata IDS → honey_node
  → RabbitMQ (alert_topic) → alert_server → Elasticsearch

honey_server → RabbitMQ 交换机 → honey_node
  (createIpExchange, bindPortExchange, deleteIpExchange)

honey_server ←gRPC 流式通信→ honey_node
  (双向命令/状态通信)
```

### 基础设施

- **RabbitMQ** (5672, 15672)：消息代理，支持 SSL，运行在 10.3.0.8
- **Elasticsearch** (9200)：告警存储，运行在 10.3.0.7
- **MySQL**：honey_server、alert_server、image_server 的共享数据库
- **Redis**：所有服务器的缓存层
- **Suricata IDS**：在 honey_node 主机上的 Docker 中运行（输出 eve.json）

## 构建和运行

### 启动基础设施
```bash
cd deploy
docker-compose up -d  # 启动 RabbitMQ + Elasticsearch
```

### 生成 SSL 证书（仅首次运行）
```bash
cd apps/honey_server
bash ca.sh  # 创建 RabbitMQ SSL 证书
# 将生成的证书复制到其他应用的 mq_cert 目录
```

### 启动管理服务器
```bash
cd apps/honey_server
go mod download
go run main.go -db          # 初始化数据库架构
go run main.go -m user -t create -v '{"username":"admin","password":"admin123"}'
go run main.go              # 启动服务器（HTTP :8079, gRPC :8081）
```

### 启动告警服务器
```bash
cd apps/alert_server
go mod download
go run main.go              # 运行在 :8078
```

### 启动镜像服务器
```bash
cd apps/image_server
go mod download
go run main.go              # 运行在 :8080
```

### 部署节点代理
```bash
cd apps/honey_node
go mod download
cd deploy && docker-compose up -d  # 部署 Suricata IDS
cd .. && go run main.go            # 自动注册到 honey_server
```

### 其他有用的命令
```bash
# 查看应用版本
go run main.go -vv

# 重新生成 protobuf 文件（honey_server）
cd apps/honey_server && bash rpc.sh
```

## 配置

每个应用都有 `settings.yaml` 配置文件，包含：
- **db**：MySQL 连接（host、port、db_name、user、password）
- **redis**：Redis 连接（addr、db）
- **mq**：RabbitMQ SSL 连接（host、port、ssl、certificates）
- **system**：监听地址（webAddr、grpcAddr）
- **jwt**：身份验证密钥（仅在 honey_server 中）

示例：
```yaml
db:
  host: "111.111.111.133"
  port: 3306
  db_name: "honey_db"
mq:
  host: "111.111.111.133"
  port: 5671
  ssl: true
  clientCertificate: mq_cert/client_certificate.pem
```

## 代码组织

### honey_server
- `main.go`：入口点，初始化所有服务
- `internal/api/`：REST API 处理器（Gin 路由）
- `internal/service/grpc_service/`：gRPC 服务器实现
- `internal/service/mq_service/`：RabbitMQ 生产者
- `internal/model/`：GORM 数据库模型
- `global/`：共享实例（数据库、Redis、MQ 连接）

### honey_node
- `main.go`：启动命令循环、定时任务、MQ 消费者
- `internal/service/command/`：gRPC 双向流式通信
- `internal/service/suricata_service/`：IDS 日志监控
- `internal/service/ip_service/`：网络接口操作（arping）
- `internal/service/port_service/`：TCP 隧道/转发
- `internal/service/mq_service/`：RabbitMQ 消费者

### alert_server
- `internal/service/mq_service/rev_alert_mq.go`：告警消费
- `internal/es_model/`：Elasticsearch 数据模型
- `internal/api/alert_api/`：查询 API

### image_server
- `internal/service/vs_net_service/`：Docker 虚拟网络管理
- `internal/api/vs_api/`：虚拟服务 API
- `internal/api/mirror_cloud_api/`：容器镜像管理

## 数据库模型

主要的 MySQL 表（通过 GORM）：
- **nodes**：已注册的蜜罐节点
- **honey_ips**：虚拟 IP 地址
- **honey_ports**：端口转发规则
- **services**：服务类型定义
- **users**：管理员用户
- **white_ips**：IP 白名单（仅 alert_server）
- **nets**：网络段
- **hosts**：目标主机
- **node_networks**：节点网络接口

## RabbitMQ 拓扑

**交换机**（honey_server 使用）：
- `createIpExchange` (direct)：IP 创建命令
- `deleteIpExchange` (direct)：IP 删除命令
- `bindPortExchange` (direct)：端口绑定命令

**队列**：
- `alert_topic`：从节点到 alert_server 的告警消息

所有 MQ 连接使用 SSL 和客户端证书。

## API 端点

### honey_server (:8079)
- `/honey_server/login`, `/honey_server/captcha`
- `/honey_server/user/*`, `/honey_server/node/*`
- `/honey_server/honey_ip/*`, `/honey_server/honey_port/*`
- `/honey_server/host/*`, `/honey_server/net/*`

### alert_server (:8078)
- `/alert_server/alert/*` - 查询告警
- `/alert_server/white_ip/*` - 白名单管理

### image_server (:8080)
- `/image_server/mirror_cloud/*` - 容器镜像
- `/image_server/vs/*` - 虚拟服务
- `/image_server/host_template/*`, `/image_server/matrix_template/*`

## 测试

```bash
cd apps/<app_name>
go test ./...
```

`testdata/` 目录中的测试示例演示了 Elasticsearch、JWT、Docker API 的使用。

## 重要说明

- **需要 Linux 系统**：honey_node 使用 arping 和网络接口操作
- **Suricata 集成**：honey_node 监控来自 Suricata Docker 容器的 eve.json 日志文件
- **SSL 证书**：首次运行前必须生成并分发证书
- **网络模式**：Suricata 使用 `network_mode: host` 来监控物理接口
- **Go 版本**：所有应用使用 Go 1.24+
- **端口隧道**：honey_node 创建从蜜罐 IP 到真实服务的 TCP 隧道
