# 部署指南

本文档提供了部署修复后系统的完整步骤。

## 前置条件

- Go 1.21.0+
- MySQL 8.0.23+
- Redis 6.2.3+
- Etcd 3.4.13+

## 快速开始

### 1. 克隆并检出分支

```bash
git checkout investigate-design-issues
```

### 2. 安装依赖

```bash
go mod download
go mod tidy
```

### 3. 验证编译

```bash
go build ./...
```

### 4. 运行测试

```bash
# 运行所有测试
go test ./...

# 只运行加密相关测试
go test ./pkg/encrypt/ -v

# 运行集成测试
go test ./test/ -v
```

## 环境配置

### 1. 复制环境变量模板

```bash
cp .env.example .env
```

### 2. 配置密钥（重要！）

编辑 `.env` 文件：

```bash
# 生产环境必须修改这些值！
PASSWORD_ENCRYPT_SEED="your-random-seed-here"
PHONE_AES_KEY="your-32-byte-key-here-12345678"

# 数据库配置
MYSQL_DSN="user:password@tcp(localhost:3306)/infoflow?charset=utf8mb4&parseTime=True"

# Redis配置
REDIS_HOST="localhost:6379"
REDIS_PASS=""

# Kafka配置
KAFKA_BROKERS="localhost:9092"

# OSS配置
OSS_ENDPOINT="oss-cn-hangzhou.aliyuncs.com"
OSS_ACCESS_KEY_ID="your-key-id"
OSS_ACCESS_KEY_SECRET="your-key-secret"
OSS_BUCKET_NAME="your-bucket"
```

### 3. 生成安全密钥

```bash
# 生成随机密码种子
openssl rand -base64 32

# 生成32字节AES密钥
openssl rand -hex 16
```

## 数据库迁移

### 1. 备份现有数据库

```bash
mysqldump -u root -p infoflow_user > backup_user_$(date +%Y%m%d).sql
```

### 2. 运行迁移脚本

```bash
mysql -u root -p infoflow_user < db/migrations/001_update_user_table.sql
```

### 3. 验证迁移

```bash
mysql -u root -p infoflow_user -e "SHOW INDEX FROM user;"
```

应该看到：
- `ix_mtime`
- `ix_username`
- `ix_phone`

```bash
mysql -u root -p infoflow_user -e "DESCRIBE user;"
```

密码字段应该是 `varchar(256)`

## 服务启动

### User RPC Service

```bash
cd applications/user/rpc
go run user.go -f etc/user.yaml
```

### Article API Service

```bash
cd applications/article/api
go run article.go -f etc/article-api.yaml
```

### Applet API Service

```bash
cd applications/applet
go run applet.go -f etc/applet-api.yaml
```

### Like RPC Service

```bash
cd applications/like/rpc
go run like.go -f etc/like.yaml
```

## 验证部署

### 1. 健康检查

```bash
# 检查各服务是否正常运行
curl http://localhost:8888/health  # Applet API
curl http://localhost:8889/health  # Article API
```

### 2. 测试用户注册

```bash
curl -X POST http://localhost:8888/user/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test123!",
    "phone": "13800138000"
  }'
```

### 3. 测试用户登录

```bash
curl -X POST http://localhost:8888/user/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test123!"
  }'
```

应该返回JWT token。

### 4. 验证密码加密

登录数据库检查：

```sql
SELECT username, password FROM user WHERE username = 'testuser';
```

新用户的密码应该以 `$2` 开头（bcrypt），长度约60字符。

## 监控

### 关键指标

1. **登录成功率**: 应 >99%
2. **登录响应时间**: P95应 <500ms
3. **密码验证错误**: 监控失败次数

### 日志监控

```bash
# 监控错误日志
tail -f /var/log/infoflow/error.log | grep -E "Login|Password|Verification"

# 监控性能
tail -f /var/log/infoflow/access.log | grep -E "POST.*login"
```

### 告警配置

建议配置以下告警：
- 登录失败率 >5%
- 登录响应时间 P95 >1s
- 服务panic事件
- 数据库连接失败

## 回滚计划

如果出现问题，可以回滚：

### 1. 回滚代码

```bash
git checkout main
go build ./...
# 重启服务
```

### 2. 数据兼容性

- ✅ 旧代码仍能读取新数据库（索引不影响）
- ✅ 旧代码仍能验证MD5密码
- ⚠️ Bcrypt密码在旧代码中无法验证

### 3. 数据库回滚（如果必要）

```bash
# 恢复备份
mysql -u root -p infoflow_user < backup_user_20240101.sql
```

## 性能优化

### 1. Bcrypt工作因子调整

如果bcrypt太慢，可以在 `pkg/encrypt/encrypt.go` 中调整：

```go
// 降低成本（不推荐）
bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

// 默认成本
bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// 提高安全性（更慢）
bcrypt.GenerateFromPassword([]byte(password), bcrypt.MaxCost)
```

### 2. 连接池优化

在各服务的配置文件中：

```yaml
DataSource: "user:pass@tcp(localhost:3306)/db?maxIdleConns=10&maxOpenConns=100"
```

### 3. Redis缓存

确保Redis配置正确：

```yaml
CacheRedis:
  - Host: localhost:6379
    Pass: ""
    Type: node
```

## 密码迁移策略

### 渐进式迁移（推荐）

用户登录时自动升级密码：

1. 用户使用旧密码登录
2. 系统验证MD5哈希成功
3. 后台生成新bcrypt哈希
4. 更新数据库中的密码字段
5. 下次登录使用bcrypt验证

实现示例（可选）：

```go
// 在loginlogic.go中
func (l *LoginLogic) verifyAndLogin(user *model.User, password, identifier string) (*pb.LoginResponse, error) {
    if user == nil {
        return nil, xcode.New(xcode.ParameterError, "用户不存在")
    }

    if !encrypt.VerifyPassword(user.Password, password) {
        return nil, xcode.New(xcode.ParameterError, "密码错误")
    }

    // 如果是MD5密码，触发重新哈希
    if len(user.Password) == 32 {
        go func() {
            newHash := encrypt.EncPassword(password)
            // 异步更新数据库
            l.svcCtx.UserModel.UpdatePassword(context.Background(), user.Id, newHash)
        }()
    }

    return &pb.LoginResponse{UserId: user.Id}, nil
}
```

### 批量迁移（可选）

如果需要一次性迁移所有密码（不推荐，因为需要明文密码）：

```sql
-- 仅用于演示，实际场景无法获取明文密码
-- 只能通过渐进式迁移
```

## 安全最佳实践

### 1. 环境变量

```bash
# 使用systemd环境文件
# /etc/systemd/system/infoflow-user.service
[Service]
EnvironmentFile=/etc/infoflow/.env
```

### 2. 密钥轮换

定期更换密钥：
1. 生成新密钥
2. 同时支持新旧密钥验证（过渡期）
3. 逐步迁移数据
4. 移除旧密钥

### 3. 日志审计

启用审计日志：
- 所有登录尝试
- 密码修改操作
- 敏感数据访问

## 故障排查

### 问题1: 登录失败

**症状**: 所有用户无法登录

**检查**:
```bash
# 检查环境变量
echo $PASSWORD_ENCRYPT_SEED

# 检查数据库连接
mysql -u root -p -e "SELECT 1"

# 检查服务日志
tail -f /var/log/infoflow/user-rpc.log
```

### 问题2: 性能下降

**症状**: 登录响应时间增加

**检查**:
```bash
# 检查bcrypt性能
go test -bench=BenchmarkVerifyPassword ./pkg/encrypt/

# 检查数据库连接
mysql -u root -p -e "SHOW PROCESSLIST"

# 检查Redis
redis-cli info
```

### 问题3: 旧密码无法验证

**症状**: 老用户无法登录

**检查**:
```bash
# 确认向后兼容代码存在
grep "Md5Sum" pkg/encrypt/encrypt.go

# 测试MD5验证
go test -v ./pkg/encrypt/ -run TestVerifyPasswordBackwardCompatibility
```

## 支持与联系

如有问题，请查看：
- `DESIGN_ISSUES.md` - 设计问题详解
- `SECURITY_IMPROVEMENTS.md` - 安全改进说明
- `TEST_VERIFICATION_REPORT.md` - 测试报告

## 附录

### A. 配置文件示例

#### user.yaml
```yaml
Name: user.rpc
ListenOn: 0.0.0.0:8080

Etcd:
  Hosts:
  - localhost:2379
  Key: user.rpc

DataSource: user:pass@tcp(localhost:3306)/infoflow_user?charset=utf8mb4&parseTime=True

CacheRedis:
  - Host: localhost:6379
    Pass: ""
    Type: node

BizRedis:
  Host: localhost:6379
  Pass: ""
  Type: node
```

### B. 系统要求

| 组件 | 最低要求 | 推荐配置 |
|------|---------|---------|
| CPU | 2核 | 4核+ |
| 内存 | 4GB | 8GB+ |
| 磁盘 | 50GB | 100GB+ SSD |
| 网络 | 100Mbps | 1Gbps |

### C. 容量规划

| 并发用户 | CPU | 内存 | 数据库 |
|---------|-----|------|--------|
| 100 | 2核 | 4GB | 1实例 |
| 1000 | 4核 | 8GB | 主从 |
| 10000 | 8核+ | 16GB+ | 集群 |

---

**文档版本**: 1.0
**最后更新**: 2024
**维护者**: DevOps Team
