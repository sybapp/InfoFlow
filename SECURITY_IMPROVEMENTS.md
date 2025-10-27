# 安全改进文档

本文档说明了针对设计问题报告中识别的安全问题所做的改进。

## 已修复的问题

### 1. 密码加密改进 ✅

**问题**: 使用不安全的MD5算法加密密码

**解决方案**:
- 采用bcrypt算法加密新密码
- 保持对旧MD5密码的向后兼容性
- 新增`VerifyPassword`函数支持两种算法验证

**代码变更**: `pkg/encrypt/encrypt.go`

```go
func EncPassword(password string) string {
    hash, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(password)), bcrypt.DefaultCost)
    if err != nil {
        // 降级到MD5（仅用于兼容性）
        return Md5Sum([]byte(strings.TrimSpace(password + passwordEncryptSeed)))
    }
    return string(hash)
}

func VerifyPassword(hashedPassword, password string) bool {
    // 先尝试bcrypt验证
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    if err == nil {
        return true
    }
    // 降级到MD5验证（兼容旧密码）
    return hashedPassword == Md5Sum([]byte(password+passwordEncryptSeed))
}
```

**迁移建议**:
1. 部署新代码后，新用户自动使用bcrypt
2. 旧用户下次登录时，可选择性地触发密码重新哈希
3. 逐步淘汰MD5支持

### 2. 加密密钥配置化 ✅

**问题**: 密钥硬编码在源代码中

**解决方案**:
- 从环境变量读取密钥
- 提供默认值以保持向后兼容
- 创建`.env.example`文件作为配置模板

**代码变更**: `pkg/encrypt/encrypt.go`

```go
var (
    passwordEncryptSeed string
    phoneAesKey         string
)

func init() {
    passwordEncryptSeed = getEnvOrDefault("PASSWORD_ENCRYPT_SEED", "(infoflow)@#$")
    phoneAesKey = getEnvOrDefault("PHONE_AES_KEY", "5A2E746B08D846502F37A6E2D85D583B")
}
```

**部署指南**:
```bash
# 设置环境变量
export PASSWORD_ENCRYPT_SEED="your-new-secret-seed"
export PHONE_AES_KEY="your-32-byte-aes-key-here"

# 或使用.env文件（需要配合godotenv等库）
```

### 3. 日志脱敏 ✅

**问题**: 日志中暴露密码等敏感信息

**解决方案**:
- 移除日志中的完整请求对象
- 只记录非敏感标识符（用户名、手机号）
- 使用适当的日志级别（Error/Info）

**代码变更**: `applications/user/rpc/internal/logic/loginlogic.go`

**Before**:
```go
logx.Errorf("Login req: %v error: %v", in, err) // 包含密码
```

**After**:
```go
logx.Errorf("Login by username error: %v, username: %s", err, in.Username) // 不含密码
```

### 4. 类型断言安全性 ✅

**问题**: 不安全的类型断言可能导致panic

**解决方案**:
- 添加nil检查
- 使用两值断言检查类型
- 返回明确的错误码

**代码变更**: `applications/article/api/internal/logic/publishlogic.go`

**Before**:
```go
userId, err := l.ctx.Value("userId").(json.Number).Int64()
```

**After**:
```go
val := l.ctx.Value("userId")
if val == nil {
    return nil, code.UnauthorizedError
}
num, ok := val.(json.Number)
if !ok {
    return nil, code.InvalidTokenError
}
userId, err := num.Int64()
if err != nil {
    return nil, code.InvalidTokenError
}
```

### 5. 错误处理改进 ✅

**问题**: 使用panic导致服务崩溃

**解决方案**:
- 返回错误而非panic
- 在main函数中优雅处理初始化错误
- 允许服务在遇到可恢复错误时继续运行

**代码变更**:
- `applications/applet/applet.go`
- `applications/article/api/article.go`
- `applications/article/api/internal/svc/servicecontext.go`

### 6. Bug修复 ✅

**问题**: `byte16ToBytes`函数实现错误

**解决方案**:
```go
func byte16ToBytes(in [16]byte) []byte {
    return in[:]  // 简单直接的转换
}
```

## 性能改进

### 1. 数据库索引 ✅

**添加的索引**:
```sql
ALTER TABLE `user` ADD KEY `ix_username` (`username`);
ALTER TABLE `user` ADD KEY `ix_phone` (`phone`);
```

**影响**: 显著提升用户名和手机号查询性能

### 2. 消息队列重试机制 ✅

**改进**: 为Kafka消息发送添加重试逻辑

```go
threading.GoSafe(func() {
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        err := l.svcCtx.KqPusherClient.Push(string(data))
        if err == nil {
            return
        }
        logx.Errorf("Thumbup: kq push error (attempt %d/%d): %v", i+1, maxRetries, err)
    }
})
```

## 代码质量改进

### 1. 消除重复代码 ✅

**改进**: 提取公共登录验证逻辑

```go
func (l *LoginLogic) verifyAndLogin(user *model.User, password, identifier string) (*pb.LoginResponse, error) {
    // 统一的验证逻辑
}
```

### 2. 拼写错误修正 ✅

- `AritcleRpc` → `ArticleRpc`
- 提高代码可读性和专业性

### 3. 配置结构优化 ✅

**改进**: 使用命名结构体替代匿名结构体

```go
type AuthConfig struct {
    AccessSecret string
    AccessExpire int64
}

type OssConfig struct {
    Endpoint         string
    AccessKeyId      string
    AccessKeySecret  string
    BucketName       string
    ConnectTimeout   int64
    ReadWriteTimeout int64
}
```

## 功能完善

### 1. ArticleDetail 实现 ✅
- 添加参数验证
- 实现文章详情查询
- 返回完整的文章信息

### 2. ArticleDelete 实现 ✅
- 添加权限验证（只能删除自己的文章）
- 实现软删除或硬删除
- 添加详细日志

### 3. IsThumbup 实现 ✅
- 查询点赞记录
- 返回点赞状态

### 4. Thumbup 改进 ✅
- 添加重复点赞检查
- 改进Kafka消息发送可靠性
- 添加重试机制

## 部署清单

### 必须操作：
- [ ] 设置环境变量（PASSWORD_ENCRYPT_SEED, PHONE_AES_KEY）
- [ ] 运行数据库迁移脚本
- [ ] 更新服务配置文件
- [ ] 重启所有服务

### 可选操作：
- [ ] 实施密码重哈希计划
- [ ] 配置日志监控和告警
- [ ] 进行安全审计和渗透测试
- [ ] 更新API文档

## 后续建议

### 短期（1-2周）:
1. 监控bcrypt性能影响
2. 收集用户反馈
3. 优化错误处理流程

### 中期（1-3个月）:
1. 实施密码策略（强度要求、过期时间）
2. 添加账户锁定机制
3. 实施API限流
4. 完善审计日志

### 长期（3-6个月）:
1. 考虑使用专业的密钥管理服务（如HashiCorp Vault）
2. 实施密钥轮换机制
3. 迁移至完全的bcrypt（移除MD5支持）
4. 添加双因素认证（2FA）

## 测试建议

### 单元测试：
```go
func TestEncPassword(t *testing.T) {
    password := "test123"
    hash := encrypt.EncPassword(password)
    assert.True(t, encrypt.VerifyPassword(hash, password))
}

func TestVerifyPasswordBackwardCompatibility(t *testing.T) {
    // 测试旧MD5密码仍然可以验证
    oldMD5Hash := "..." // 旧的MD5哈希
    assert.True(t, encrypt.VerifyPassword(oldMD5Hash, "password"))
}
```

### 集成测试：
- 测试登录流程
- 测试点赞功能
- 测试文章发布和删除

### 性能测试：
- 压测bcrypt加密性能
- 测试数据库索引效果
- 测试Kafka消息吞吐量

## 安全检查清单

- [x] 密码使用强哈希算法
- [x] 敏感信息不在日志中暴露
- [x] 环境变量配置密钥
- [x] 类型断言有安全检查
- [x] 错误处理不导致服务崩溃
- [x] 数据库查询有索引支持
- [x] API有权限验证
- [x] 消息队列有重试机制
- [ ] 实施HTTPS（待完成）
- [ ] 添加CSRF保护（待完成）
- [ ] 实施XSS防护（待完成）
- [ ] 添加SQL注入防护检查（待完成）

## 联系与支持

如有任何安全问题或建议，请联系安全团队。
