# Gate Copy Trading - UniApp 前端

基于 UniApp 的 Gate.io 跟单交易前端应用，支持微信小程序、H5、App 等多端运行。

## 功能特性

- ✅ 用户认证（登录/注册）
- ✅ 交易员列表展示（筛选、排序）
- ✅ 交易员详情页
- ✅ 收益曲线图
- ✅ 个人中心
- ✅ API Key 配置

## 技术栈

- **框架**: UniApp (Vue 3 + TypeScript)
- **UI 组件**: uView UI
- **状态管理**: Pinia
- **图表**: uCharts

## 快速开始

### 1. 安装依赖

```bash
npm install
```

### 2. 配置 API 地址

修改 `api/request.ts` 中的 `BASE_URL`：

```typescript
const BASE_URL = '/api'  // H5 开发时使用代理
```

### 3. 运行项目

**微信小程序**:
```bash
npm run dev:mp-weixin
```

**H5**:
```bash
npm run dev:h5
```

**App**:
```bash
npm run dev:app
```

### 4. 启动后端 API

确保后端服务运行在 `http://localhost:8080`

```bash
cd ..
go run ./cmd/api/
```

## 项目结构

```
uniapp-frontend/
├── pages/                    # 页面
│   ├── index/               # 首页（交易员列表）
│   ├── trader/              # 交易员详情
│   ├── stats/               # 历史统计
│   ├── auth/                # 认证
│   ├── mine/                # 个人中心
│   └── settings/            # 设置
├── components/              # 组件
├── store/                   # Pinia 状态管理
│   └── modules/
│       ├── auth.ts
│       └── trader.ts
├── api/                     # API 封装
│   ├── request.ts
│   ├── auth.ts
│   ├── trader.ts
│   └── user.ts
├── types/                   # TypeScript 类型
├── utils/                   # 工具函数
└── static/                  # 静态资源
```

## API 接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/auth/login` | POST | 用户登录 |
| `/api/auth/register` | POST | 用户注册 |
| `/api/user` | GET | 获取用户信息 |
| `/api/user/api-key` | POST | 绑定 API Key |
| `/api/copytrading/traders` | GET | 交易员列表 |
| `/api/copytrading/traders/:id` | GET | 交易员详情 |
| `/api/copytrading/traders/:id/stats` | GET | 历史统计 |

## 构建发布

**微信小程序**:
```bash
npm run build:mp-weixin
```

**H5**:
```bash
npm run build:h5
```

## 注意事项

1. 首次使用需要先注册账号
2. 配置 API Key 后才能使用完整功能
3. Cookie 过期需要重新绑定

## License

MIT
