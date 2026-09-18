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

1. **登录认证** — 管理员 / 记录员角色，JWT 鉴权
2. **工作台（默认首页）** — 今日新增文物数、最近 5 条登记文物、各实体总数，卡片可点击进入对应管理页
3. **发掘工地 Site** — 名称、时代、经纬度、负责人
4. **探方/发掘单位 Unit** — 所属工地、编号、深度区间、地层简述
5. **出土文物 Find** — 所属探方、登记号、器物类型、材质、完整度、出土日期、描述、存放位置
6. **材质分类 Material** — 名称、描述（字典表）
7. **概览页** — 工地数、探方数、文物总数、按器物类型统计

### 工作台口径说明

- **今日新增文物**：按**东八区（UTC+8）自然日**（当日 00:00:00 至次日 00:00:00）统计 `finds.created_at`（登记时间），与出土日期 `find_date` 无关；不含软删除记录。
- **最近登记文物**：按登记先后取最近 5 条（`id` 倒序），含登记号、工地/探方、器物类型、材质、完整度、出土日期、登记时间。
- **各实体总数**：工地数、探方数、文物总数，与概览页同一统计口径（复用同一后端查询，均不含软删除）。
- 工作台数据全部来自 `GET /api/workspace/today` 实时查询，无静态假数据。

## API 前缀

所有接口以 `/api` 开头：

- `POST /api/auth/login`
- `GET|POST|PUT|DELETE /api/sites`
- `GET|POST|PUT|DELETE /api/units`
- `GET|POST|PUT|DELETE /api/finds`
- `GET|POST|PUT|DELETE /api/materials`
- `GET /api/overview`
- `GET /api/workspace/today` — 工作台今日摘要（今日新增文物数、最近 5 条文物、各实体总数）

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
