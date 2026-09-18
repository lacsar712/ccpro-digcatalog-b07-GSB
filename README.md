# 考古发掘出土文物编目系统（DigCatalog）

面向考古工地出土文物登记与编目的全栈演示项目：支持发掘工地、探方/发掘单位、出土文物、材质字典的 CRUD，以及概览统计。

## 技术栈

- **前端**: Vue 3 + Vite + Pinia + Vue Router（Composition API + `<script setup>`）
- **后端**: Go 1.21+ + Gin + GORM
- **数据库**: MySQL 8.0
- **认证**: JWT + bcrypt

## 一键启动

```bash
docker compose up --build
```

启动完成后访问：

| 服务 | 地址 |
|------|------|
| 前端 | http://localhost:3200 |
| 后端 API | http://localhost:8200/api |
| MySQL | localhost:3307（用户 `root` / 密码 `root`，库名 `digcatalog`） |

停止服务：

```bash
docker compose down
```

清除数据卷后重建：

```bash
docker compose down -v
docker compose up --build
```

## 测试账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| `admin` | `123456` | 管理员 |
| `recorder` | `123456` | 记录员 |

## 功能模块

1. **登录认证** — 管理员 / 记录员角色，JWT 鉴权；未登录访问任意页面均跳转登录页
2. **工作台（默认首页）** — 今日新增文物数、各实体总数卡片（可点击跳转 Finds/Sites/Units）、最近 5 条登记文物
3. **发掘工地 Site** — 名称、时代、经纬度、负责人
4. **探方/发掘单位 Unit** — 所属工地、编号、深度区间、地层简述
5. **出土文物 Find** — 所属探方、登记号、器物类型、材质、完整度、出土日期、描述、存放位置
6. **材质分类 Material** — 名称、描述（字典表）
7. **概览页** — 工地数、探方数、文物总数、按器物类型统计（保留在侧栏「概览」）

## 工作台统计口径

工作台数据由 `GET /api/workspace/today` 实时查询数据库返回，**不使用任何静态假数据**：

| 指标 | 口径 |
|------|------|
| 今日新增文物 | `finds.created_at` 落在**东八区（UTC+8）当日 [00:00, 24:00)** 的记录数；后端用 `time.FixedZone("UTC+8", 8*3600)` 计算当日边界，不依赖容器 tzdata，与 `Asia/Shanghai` 等价（中国已无夏令时） |
| 最近登记文物 | `finds` 按 `id` 倒序取前 5 条，预加载所属探方与工地 |
| 各实体总数 | 工地 / 探方 / 文物总数，与 `GET /api/overview` 复用同一统计函数（`countEntities`），口径完全一致 |

接口响应附带 `timezone`（`UTC+8`）与 `today`（东八区当日日期，如 `2026-09-18`）字段，前端按 `Asia/Shanghai` 时区展示登记时间，与后端口径一致。

## API 前缀

所有接口以 `/api` 开头：

- `POST /api/auth/login`
- `GET|POST|PUT|DELETE /api/sites`
- `GET|POST|PUT|DELETE /api/units`
- `GET|POST|PUT|DELETE /api/finds`
- `GET|POST|PUT|DELETE /api/materials`
- `GET /api/overview`
- `GET /api/workspace/today`

前端经 Nginx 将 `/api` 反代至后端容器 `http://backend:8080`。

## 端口映射

| 服务 | 宿主机 | 容器内 |
|------|--------|--------|
| Frontend | 3200 | 80 |
| Backend | 8200 | 8080 |
| MySQL | 3307 | 3306 |

## 目录结构

```
DigCatalog/
├── docker-compose.yml
├── README.md
├── .gitignore
├── backend/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   └── internal/
│       ├── config/
│       ├── models/
│       ├── handlers/
│       ├── middleware/
│       └── seed/
└── frontend/
    ├── Dockerfile
    ├── nginx.conf
    ├── package.json
    ├── vite.config.js
    ├── index.html
    └── src/
```

## 本地开发（可选）

### 后端

```bash
cd backend
go mod tidy
# 确保 MySQL 已启动且环境变量正确
go run .
```

### 前端

```bash
cd frontend
npm install --registry=https://registry.npmmirror.com
npm run dev
```
