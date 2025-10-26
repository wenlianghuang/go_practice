# 快速开始指南 ⚡

## 🎯 最简单的方式（推荐）

**无需安装任何额外软件！**

```bash
# 进入项目目录
cd complex_breakpoint_example

# 直接运行
go run main.go
```

就这几步！项目使用**内存数据库**，完全不需要：
- ❌ 不需要 Docker
- ❌ 不需要 PostgreSQL  
- ❌ 不需要配置文件
- ❌ 不需要任何额外安装

## ✅ 使用内存数据库（当前默认）

### 运行应用
```bash
go run main.go
```

### 测试 API
服务器会在 `:9090` 端口启动：

```bash
# 健康检查
curl http://localhost:9090/health

# 获取所有用户
curl http://localhost:9090/api/v1/users

# 获取特定用户（ID 1）
curl http://localhost:9090/api/v1/users/1

# 创建用户
curl -X POST http://localhost:9090/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","balance":1000}'
```

## 📋 完整的 API 端点

```
POST   /api/v1/users              - 创建用户
GET    /api/v1/users              - 获取所有用户
GET    /api/v1/users/{id}         - 获取用户
GET    /api/v1/users/{id}/account - 获取用户账户
GET    /api/v1/users/{id}/transactions - 获取用户交易
GET    /api/v1/users/{id}/loans   - 获取用户贷款申请
POST   /api/v1/users/{id}/deposit - 存款
POST   /api/v1/users/{id}/withdraw - 提款
POST   /api/v1/users/{id}/apply-loan - 申请贷款
POST   /api/v1/transfer           - 转账
POST   /api/v1/test/concurrent   - 并发测试
GET    /health                    - 健康检查
```

## 🧪 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定测试
go test ./handlers -v
```

## 📦 编译

```bash
# 编译为可执行文件
go build -o complex_breakpoint_example .

# 运行编译后的文件
./complex_breakpoint_example
```

## 🔧 如果需要 PostgreSQL（可选）

目前项目已配置好 GORM 支持，但**不需要**立即使用。

如果将来需要使用 PostgreSQL：

1. 查看 `GORM_README.md` 了解安装 PostgreSQL 的方法
2. 安装 PostgreSQL（可以使用 Homebrew、apt、或直接从官网下载）
3. 配置环境变量并修改 `main.go`

## 💡 提示

- 端口可在环境变量中设置：`export PORT=8080`
- 日志会输出到控制台
- 内存数据库在重启后会重置数据

## ❓ 有问题？

- 查看 `README.md` 了解项目概述
- 查看 `DEBUG_GUIDE.md` 了解调试方法
- 查看 `GO_CHI_DEBUG_GUIDE.md` 了解 Go Chi 框架调试

