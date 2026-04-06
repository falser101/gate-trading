# Gate Trading Bot

基于 Go + Flutter 的 Gate.io 自动交易系统

## 技术栈

- **后端**: Go 1.21+, Gin, GORM, PostgreSQL, Redis
- **前端**: Flutter 3.x (Web/iOS/Android)
- **交易所**: Gate.io API v4

## 快速开始

### 1. 环境准备

```bash
# 安装依赖
go mod download

# 复制环境配置
cp .env.example .env

# 编辑 .env 配置你的数据库和 Gate API Key
```

### 2. 启动服务

```bash
# 使用 Docker Compose (推荐)
docker-compose up -d

# 或者直接运行
go run cmd/api/main.go
```

### 3. 测试 API

```bash
curl http://localhost:8080/api/health
```

## 项目结构

```
.
├── cmd/
│   ├── api/           # HTTP API 服务
│   ├── worker/        # 策略执行后台任务
│   └── ws/            # WebSocket 推送服务
├── internal/
│   ├── config/        # 配置管理
│   ├── models/        # 数据模型
│   ├── repository/    # 数据访问层
│   ├── service/       # 业务逻辑层
│   │   ├── auth/      # 认证服务
│   │   ├── strategy/  # 策略引擎
│   │   ├── trading/   # 交易服务
│   │   └── market/    # 行情服务
│   ├── gateway/       # Gate.io API 封装
│   └── middleware/    # 中间件
├── pkg/
│   └── gateapi/       # Gate API SDK
└── deploy/
    └── docker/        # Docker 配置
```

## 功能特性

- [x] 网格交易策略 (Grid)
- [x] 定投策略 (DCA)
- [ ] 组合策略 (Combo)
- [ ] 回测系统
- [ ] 多交易所支持

## 开发文档

- [Gate.io API v4 文档](https://www.gate.com/docs/developers/apiv4/)
- [Gate API GitHub](https://github.com/gateio/rest-v4)

## License

MIT
