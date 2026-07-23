# 石小易 AI 迎新助手 · 工单系统

> 当知识库检索不到答案时，智能体自动提交工单 → 人工补答 → 用户凭工单号查结果。

## 📦 项目结构

```
ticket-system/
├── backend/              # Go 后端（Gin + GORM + SQLite）
│   ├── main.go           # 入口
│   ├── config/           # 配置（JWT Secret、默认 admin 等）
│   ├── models/           # 数据模型（工单 + 用户）
│   ├── handlers/         # API 处理器（工单 / 认证 / 用户管理）
│   ├── middleware/       # 双轨鉴权（API Key + JWT）
│   ├── go.mod
│   └── Dockerfile
├── frontend/             # Vue 3 + Element Plus 管理后台
│   ├── src/
│   │   ├── views/        # 登录 / 工单列表 / 详情 / 用户管理 / 设置
│   │   ├── api/          # axios 封装
│   │   ├── router/       # 路由守卫（强制改密码 / 权限控制）
│   │   ├── stores/       # Pinia 状态管理
│   │   └── App.vue
│   ├── vite.config.js
│   ├── nginx.conf        # 生产环境反向代理
│   └── Dockerfile
├── deploy/
│   ├── docker-compose.yml
│   └── .env.example
├── test_e2e.py           # 端到端测试脚本（53 用例）
└── README.md
```

## 🚀 快速启动

### 方式一：本地开发（推荐先跑这个）

```powershell
# 1. 启动后端
cd backend
$env:API_KEY="sk-test-key-12345"
$env:PORT="8000"
$env:JWT_SECRET="change-me-in-production"
go run .

# 2. 另开终端：启动前端
cd frontend
npm install
npm run dev
# 浏览器访问 http://localhost:5173
```

### 方式二：Docker 一键部署（生产环境）

```bash
# 1. 配置环境变量
cd deploy
cp .env.example .env
# 修改 .env 里的 API_KEY、JWT_SECRET、ADMIN_PASSWORD

# 2. 构建 + 启动
docker-compose up -d --build

# 3. 访问
# 后端 API:  http://<server-ip>:8000
# 管理后台:  http://<server-ip>:80
# 健康检查:  http://<server-ip>:80/health
```

## 🔌 API 端点速查

### 工单 API（双轨鉴权：API Key 或 JWT）

| 方法 | 路径 | 功能 |
|------|------|------|
| `POST` | `/api/v1/tickets` | 智能体调用提交工单，返回 ticket_id |
| `GET`  | `/api/v1/tickets` | 查询工单列表（分页 + 筛选） |
| `GET`  | `/api/v1/tickets/{ticket_id}` | 用户凭工单号查答复结果 |
| `POST` | `/api/v1/tickets/{ticket_id}/answer` | 人工答复 |
| `POST` | `/api/v1/tickets/{ticket_id}/close` | 关闭工单 |

### 认证 API（仅 JWT）

| 方法 | 路径 | 功能 | 权限 |
|------|------|------|------|
| `POST` | `/api/v1/auth/login` | 用户名密码登录 | 公开（无需鉴权） |
| `GET`  | `/api/v1/auth/me` | 获取当前用户信息 | 登录用户 |
| `POST` | `/api/v1/auth/change-password` | 修改密码 | 登录用户 |

### 用户管理 API（仅 admin + JWT）

| 方法 | 路径 | 功能 | 权限 |
|------|------|------|------|
| `GET`  | `/api/v1/users` | 用户列表 | admin |
| `POST` | `/api/v1/users` | 创建用户 | admin |
| `PATCH` | `/api/v1/users/{id}` | 更新用户信息 | admin |

### 健康检查

| 方法 | 路径 | 功能 |
|------|------|------|
| `GET` | `/health` | 健康检查（无需鉴权） |

完整规范见 `ticket-plugin.openapi.json`。

## 🔐 鉴权说明

本系统采用**双轨鉴权**：

### 1. API Key（智能体使用）

固定密钥通过请求头 `X-API-Key` 传递，由环境变量 `API_KEY` 配置。
适用于 HiAgent 智能体调用工单提交/查询接口。

```
X-API-Key: sk-your-api-key
```

### 2. Bearer JWT（管理后台使用）

先通过 `POST /auth/login` 获取令牌，然后在后续请求中携带：

```
Authorization: Bearer <jwt-token>
```

JWT 令牌包含用户 ID、用户名、角色（admin / staff），由环境变量 `JWT_SECRET` 签名。

### 角色权限

| 角色 | 工单操作 | 用户管理 | 修改密码 | 说明 |
|------|---------|---------|---------|------|
| `admin` | ✅ 全部 | ✅ 全部 | ✅ | 初始账号，可管理所有资源 |
| `staff` | ✅ 全部 | ❌ 无权限 | ✅ | 客服人员，仅操作工单 |

> 默认 admin 账号 `admin / admin123`（通过 `ADMIN_PASSWORD` 环境变量覆盖）。
> 首次登录强制修改密码，密码要求至少 8 位、含字母和数字。

## 🧪 端到端测试

后端启动后：

```bash
python test_e2e.py
```

**测试覆盖 53 个用例**：

- ✅ 健康检查
- ✅ 鉴权测试（无 key / 错 key / API Key / JWT）
- ✅ admin 登录、token 获取
- ✅ 提交工单、5 分钟内去重
- ✅ 查询工单、人工答复、关闭工单
- ✅ 状态流转（pending → answered → closed）
- ✅ 来源渠道（hiagent / wechat / feishu / yiban）
- ✅ 参数校验、404
- ✅ 错误密码验证
- ✅ /auth/me 用户信息
- ✅ 弱密码校验（太短 / 无数字）
- ✅ 修改密码、旧密码失效
- ✅ 创建 staff 用户、重复用户名拒绝
- ✅ 首次登录强制改密码
- ✅ 角色权限（staff 不能管理用户）
- ✅ admin 用户列表、编辑、禁用/启用
- ✅ 自禁用防护、最后 admin 降级防护

实测结果：**53 / 53 通过**。

> ⚠ 测试前需清理旧数据库：`mavis-trash backend/tickets.db` 然后重启后端。
> 后端启动时自动将 admin 密码重置为 `ADMIN_PASSWORD` 环境变量值。

## ⚙️ 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `API_KEY` | `sk-change-me-in-production` | API 鉴权密钥（**生产必改**） |
| `JWT_SECRET` | `change-me-in-production` | JWT 签名密钥（**生产必改**） |
| `ADMIN_PASSWORD` | `admin123` | 默认 admin 密码（首次启动创建 / 后续启动自动重置到此值） |
| `DATABASE_URL` | `tickets.db`（开发）/ `/app/data/tickets.db`（Docker） | SQLite 数据库路径 |
| `PORT` | `8000` | 后端监听端口 |
| `HOST` | `0.0.0.0` | 后端监听地址 |
| `TICKET_PREFIX` | `T` | 工单号前缀 |
| `FRONTEND_PORT` | `80` | docker-compose 暴露的前端端口 |

## 👤 用户管理操作

### 默认 admin 账号

首次启动自动创建：`admin / admin123`（可通过 `ADMIN_PASSWORD` 环境变量覆盖）

### 创建 staff 用户（通过 API）

```bash
# admin 登录获取 token
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 创建 staff 用户
curl -X POST http://localhost:8000/api/v1/users \
  -H "Authorization: Bearer <admin-token>" \
  -H "Content-Type: application/json" \
  -d '{"username":"xiaowang","password":"Wang1234","role":"staff","display_name":"小王","email":"xiaowang@shzu.edu.cn"}'
```

### 安全约束

- 不可禁用自己
- 不可将最后一名 admin 降级为 staff
- 密码至少 8 位，含字母和数字
- 首次登录（包括默认 admin）必须修改密码

## 🔄 工作流

```
新生提问 ──▶ 智能体 RAG 检索
                │
                ├─ 命中 ──▶ L1/L2 直接回答 ✅
                │
                └─ 完全无结果 ──▶ L3 兜底
                                  │
                                  ▼
                            POST /api/v1/tickets
                            返回 ticket_id + 友好话术
                                  │
                                  ▼
                            用户看到工单号
                                  │
                  ┌───────────────┼───────────────┐
                  ▼               ▼               ▼
            立即再问         几小时后再问     等不到回复
            (5min 内)        (查答复结果)
            返回原 ticket    GET /tickets/{id}
                  │               │
                  └───────┬───────┘
                          ▼
                  管理后台看到 pending
                          │
                          ▼
                  POST /tickets/{id}/answer
                  (人工补标准答案)
                          │
                          ▼
                  status = answered
                  用户再次查询即可看到答案
                          
                  可选: sync_to_kb=true
                          │
                          ▼
                  TODO: 同步到 HiAgent 知识库
```

## 📋 后续扩展（暂未实现）

- **Phase 2**: 知识库自动同步（`sync_to_kb=true` 时的实际推送）
- **Phase 3**: 数据分析（TOP 20 高频问题、周报、看板）
- **Phase 4**: 多渠道接入（飞书/易班后台入口）
- **PostgreSQL**: 当前用 SQLite 适合中小规模；日均 > 1000 单时建议切到 PostgreSQL
