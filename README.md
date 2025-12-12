# 蜜罐诱捕系统

一个基于 Go 语言开发的分布式蜜罐诱捕系统，通过创建诱饵网络服务来检测和捕获恶意活动。

![蜜罐诱捕系统架构图](https://raw.githubusercontent.com/cxposition/picture_storage/main/%E5%A8%81%E8%83%81%E8%AF%B1%E6%8D%95%E7%9F%A9%E9%98%B5/%E7%B3%BB%E7%BB%9F%E6%9E%B6%E6%9E%84%E5%9B%BE/1.png)

## 目录

- [系统简介](#系统简介)
- [核心功能](#核心功能)
- [系统架构](#系统架构)
- [技术栈](#技术栈)
- [快速开始](#快速开始)
- [配置说明](#配置说明)
- [API 文档](#api-文档)
- [开发指南](#开发指南)

## 系统简介

本系统是一个企业级分布式蜜罐平台，由 5 个微服务组件组成，通过 gRPC 和 RabbitMQ 实现服务间通信。系统能够：

- 在目标主机上动态创建虚拟 IP 和诱饵服务
- 集成 Suricata IDS 进行实时流量监控
- 自动捕获并分析恶意活动
- 将告警信息存储到 Elasticsearch 进行检索和分析
- 通过 Docker 容器编排多种蜜罐服务

## 核心功能

- **动态 IP 管理**：在节点上创建/删除虚拟 IP 地址
- **端口转发**：建立从蜜罐 IP 到真实服务的 TCP 隧道
- **入侵检测**：集成 Suricata IDS 监控所有流量
- **告警处理**：实时告警处理、地理位置丰富化、存储到 ES
- **容器编排**：管理 Docker 容器化的蜜罐服务
- **集中管理**：统一的 Web API 管理所有节点和服务

## 系统架构

### 核心组件

#### 1. honey_server - 管理服务器

蜜罐基础设施的中央编排器，提供管理 API 和节点通信。

- **端口**：8079 (HTTP)、8081 (gRPC)
- **功能**：
  - 节点注册与管理
  - IP 地址分配
  - 端口转发配置
  - 用户认证（JWT）
- **技术栈**：Gin、gRPC、GORM (MySQL)、Redis、RabbitMQ

#### 2. honey_node - 节点代理

部署在目标主机上的代理程序，执行具体的网络操作。

- **功能**：
  - 创建虚拟 IP 地址
  - 配置端口隧道
  - 监控 Suricata 日志
  - 上报告警信息
- **技术栈**：gRPC、RabbitMQ、SQLite、arping

#### 3. alert_server - 告警处理

从消息队列消费告警并进行处理分析。

- **端口**：8078
- **功能**：
  - 消费 RabbitMQ 告警消息
  - 地理位置查询（ip2region）
  - 告警数据存储到 Elasticsearch
  - 白名单管理
- **技术栈**：Gin、Elasticsearch、RabbitMQ、Redis

#### 4. image_server - 容器管理

负责 Docker 容器的镜像和服务编排。

- **端口**：8080
- **功能**：
  - Docker 镜像管理
  - 虚拟网络创建
  - 容器服务编排
  - 主机模板管理
- **技术栈**：Gin、Docker SDK、GORM、Redis

#### 5. goroutine_test - 测试代码

用于开发和测试的辅助工具（非生产环境）。

### 数据流

```
攻击者流量 → 蜜罐 IP:端口 → Suricata IDS → honey_node
           ↓
    RabbitMQ (alert_topic)
           ↓
    alert_server → Elasticsearch

honey_server → RabbitMQ 交换机 → honey_node
    (createIpExchange, bindPortExchange, deleteIpExchange)

honey_server ←gRPC 双向流→ honey_node
    (命令/状态通信)
```

### 基础设施依赖

- **RabbitMQ** (5672, 15672)：消息代理，支持 SSL
- **Elasticsearch** (9200)：告警存储和检索
- **MySQL**：关系型数据库
- **Redis**：缓存层
- **Suricata IDS**：入侵检测系统（Docker 部署）

## 技术栈

- **语言**：Go 1.24+
- **Web 框架**：Gin
- **RPC**：gRPC + Protocol Buffers
- **消息队列**：RabbitMQ (SSL)
- **数据库**：MySQL (GORM ORM)、SQLite
- **搜索引擎**：Elasticsearch
- **缓存**：Redis
- **容器化**：Docker、Docker Compose
- **IDS**：Suricata
- **认证**：JWT

## 快速开始

### 环境要求

- Linux 操作系统（需要网络接口操作权限）
- Go 1.24+
- Docker & Docker Compose
- MySQL 数据库
- Redis 服务

### 1. 启动基础设施

```bash
cd deploy
docker-compose up -d
```

这将启动 RabbitMQ 和 Elasticsearch 服务。

### 2. 生成 SSL 证书（首次运行）

```bash
cd apps/honey_server
bash ca.sh
```

将生成的证书复制到其他应用的 `mq_cert/` 目录。

### 3. 启动管理服务器

```bash
cd apps/honey_server
go mod download

# 初始化数据库
go run main.go -db

# 创建管理员用户
go run main.go -m user -t create -v '{"username":"admin","password":"admin123"}'

# 启动服务
go run main.go
```

服务将在以下端口启动：
- HTTP API：http://localhost:8079
- gRPC：localhost:8081

### 4. 启动告警服务器

```bash
cd apps/alert_server
go mod download
go run main.go
```

服务运行在 http://localhost:8078

### 5. 启动镜像服务器

```bash
cd apps/image_server
go mod download
go run main.go
```

服务运行在 http://localhost:8080

### 6. 部署节点代理

```bash
cd apps/honey_node

# 部署 Suricata IDS
cd deploy && docker-compose up -d

# 启动节点代理
cd .. && go run main.go
```

节点将自动注册到管理服务器。

### 7. 验证部署

```bash
# 查看应用版本
go run main.go -vv

# 检查服务状态
curl http://localhost:8079/honey_server/health
curl http://localhost:8078/alert_server/health
curl http://localhost:8080/image_server/health
```

## 配置说明

每个应用的 `settings.yaml` 文件包含以下配置：

### 数据库配置

```yaml
db:
  host: "111.111.111.133"
  port: 3306
  db_name: "honey_db"
  user: "root"
  password: "your_password"
```

### Redis 配置

```yaml
redis:
  addr: "111.111.111.133:6379"
  db: 0
  password: ""
```

### RabbitMQ 配置

```yaml
mq:
  host: "111.111.111.133"
  port: 5671
  ssl: true
  clientCertificate: mq_cert/client_certificate.pem
  clientKey: mq_cert/client_key.pem
  rootCertificate: mq_cert/rootCA.pem
```

### 系统配置

```yaml
system:
  webAddr: ":8079"      # HTTP 监听地址
  grpcAddr: ":8081"     # gRPC 监听地址（仅 honey_server）
```

### JWT 配置（仅 honey_server）

```yaml
jwt:
  key: "your_secret_key"
  expire: 86400         # Token 过期时间（秒）
```

## API 文档

### honey_server (:8079)

**认证相关**
- `POST /honey_server/login` - 用户登录
- `GET /honey_server/captcha` - 获取验证码

**节点管理**
- `GET /honey_server/node/*` - 节点列表
- `POST /honey_server/node/*` - 创建节点
- `PUT /honey_server/node/*` - 更新节点
- `DELETE /honey_server/node/*` - 删除节点

**IP 管理**
- `GET /honey_server/honey_ip/*` - IP 列表
- `POST /honey_server/honey_ip/*` - 创建虚拟 IP
- `DELETE /honey_server/honey_ip/*` - 删除虚拟 IP

**端口管理**
- `GET /honey_server/honey_port/*` - 端口列表
- `POST /honey_server/honey_port/*` - 绑定端口
- `DELETE /honey_server/honey_port/*` - 解绑端口

**主机管理**
- `GET /honey_server/host/*` - 主机列表
- `POST /honey_server/host/*` - 创建主机
- `PUT /honey_server/host/*` - 更新主机
- `DELETE /honey_server/host/*` - 删除主机

**网络管理**
- `GET /honey_server/net/*` - 网络列表
- `POST /honey_server/net/*` - 创建网络段
- `PUT /honey_server/net/*` - 更新网络段
- `DELETE /honey_server/net/*` - 删除网络段

### alert_server (:8078)

**告警查询**
- `GET /alert_server/alert/*` - 查询告警
- `GET /alert_server/alert/stats` - 告警统计

**白名单管理**
- `GET /alert_server/white_ip/*` - 白名单列表
- `POST /alert_server/white_ip/*` - 添加白名单
- `DELETE /alert_server/white_ip/*` - 删除白名单

### image_server (:8080)

**镜像管理**
- `GET /image_server/mirror_cloud/*` - 镜像列表
- `POST /image_server/mirror_cloud/*` - 导入镜像
- `DELETE /image_server/mirror_cloud/*` - 删除镜像

**虚拟服务**
- `GET /image_server/vs/*` - 服务列表
- `POST /image_server/vs/*` - 创建服务
- `PUT /image_server/vs/*` - 更新服务
- `DELETE /image_server/vs/*` - 删除服务

**模板管理**
- `GET /image_server/host_template/*` - 主机模板
- `GET /image_server/matrix_template/*` - 矩阵模板

## 开发指南

### 项目结构

```
honey_app/
├── apps/
│   ├── honey_server/      # 管理服务器
│   ├── honey_node/        # 节点代理
│   ├── alert_server/      # 告警服务器
│   ├── image_server/      # 容器管理
│   └── goroutine_test/    # 测试代码
├── deploy/                # Docker Compose 配置
├── CLAUDE.md              # 开发指南
└── README.md              # 本文件
```

### 各应用的代码组织

**honey_server**
```
honey_server/
├── main.go                     # 入口点
├── internal/
│   ├── api/                    # REST API 处理器
│   ├── service/
│   │   ├── grpc_service/       # gRPC 服务器实现
│   │   └── mq_service/         # RabbitMQ 生产者
│   └── model/                  # GORM 数据库模型
└── global/                     # 共享实例
```

**honey_node**
```
honey_node/
├── main.go                     # 启动入口
├── internal/
│   └── service/
│       ├── command/            # gRPC 双向流通信
│       ├── suricata_service/   # IDS 日志监控
│       ├── ip_service/         # 网络接口操作
│       ├── port_service/       # TCP 隧道/转发
│       └── mq_service/         # RabbitMQ 消费者
└── deploy/                     # Suricata 部署配置
```

**alert_server**
```
alert_server/
├── main.go
├── internal/
│   ├── service/
│   │   └── mq_service/         # 告警消费
│   ├── es_model/               # Elasticsearch 模型
│   └── api/alert_api/          # 查询 API
└── global/
```

**image_server**
```
image_server/
├── main.go
├── internal/
│   ├── service/
│   │   └── vs_net_service/     # Docker 网络管理
│   └── api/
│       ├── vs_api/             # 虚拟服务 API
│       └── mirror_cloud_api/   # 镜像管理 API
└── global/
```

### 数据库模型

主要的 MySQL 表（通过 GORM）：

- `nodes` - 已注册的蜜罐节点
- `honey_ips` - 虚拟 IP 地址
- `honey_ports` - 端口转发规则
- `services` - 服务类型定义
- `users` - 管理员用户
- `white_ips` - IP 白名单
- `nets` - 网络段
- `hosts` - 目标主机
- `node_networks` - 节点网络接口

### RabbitMQ 拓扑

**交换机**（honey_server 发布）：
- `createIpExchange` (direct) - IP 创建命令
- `deleteIpExchange` (direct) - IP 删除命令
- `bindPortExchange` (direct) - 端口绑定命令

**队列**：
- `alert_topic` - 告警消息队列（node → alert_server）

所有 MQ 连接使用 SSL 和客户端证书认证。

### 重新生成 protobuf 文件

```bash
cd apps/honey_server
bash rpc.sh
```

### 运行测试

```bash
cd apps/<app_name>
go test ./...
```

测试示例位于 `testdata/` 目录。

## 重要说明

- **需要 Linux 系统**：honey_node 需要执行 arping 和网络接口操作
- **Suricata 集成**：必须使用 `network_mode: host` 监控物理网络接口
- **SSL 证书**：首次运行前必须生成并分发 RabbitMQ SSL 证书
- **端口隧道**：honey_node 创建的 TCP 隧道需要 root 权限或 CAP_NET_ADMIN 能力
- **防火墙配置**：确保各组件间的通信端口已开放

## 故障排查

### 节点无法注册

1. 检查 honey_server 的 gRPC 端口是否可访问
2. 验证 SSL 证书配置是否正确
3. 查看节点日志中的连接错误

### 告警未存储到 Elasticsearch

1. 确认 RabbitMQ 连接正常
2. 检查 Elasticsearch 服务状态
3. 验证 alert_server 日志中的错误信息

### Suricata 未检测流量

1. 确认 Suricata 容器使用 `host` 网络模式
2. 检查网络接口配置
3. 验证 eve.json 日志文件权限

### 端口转发不工作

1. 检查虚拟 IP 是否成功创建（`ip addr show`）
2. 验证目标服务是否可达
3. 查看 honey_node 日志中的端口绑定错误

## 许可证

本项目采用私有许可证，未经授权不得用于商业用途。

## 联系方式

如有问题或建议，请联系开发团队。

---

⚠️ **安全提醒**：本系统用于合法的安全研究和防御目的，请勿用于非法活动。
