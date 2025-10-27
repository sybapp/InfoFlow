#!/bin/bash

# 验证修复脚本
# 用于快速验证所有修复是否正常工作

set -e

echo "=================================="
echo "  InfoFlow 修复验证脚本"
echo "=================================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 计数器
PASS=0
FAIL=0

# 测试函数
test_step() {
    local name=$1
    local command=$2
    
    echo -n "测试: $name ... "
    if eval "$command" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ 通过${NC}"
        ((PASS++))
        return 0
    else
        echo -e "${RED}✗ 失败${NC}"
        ((FAIL++))
        return 1
    fi
}

echo "1. 代码编译检查"
echo "--------------------------------"
test_step "编译所有包" "go build ./..."
echo ""

echo "2. 单元测试"
echo "--------------------------------"
test_step "加密包测试" "go test ./pkg/encrypt/ -v"
test_step "集成测试" "go test ./test/ -v"
echo ""

echo "3. 代码格式检查"
echo "--------------------------------"
test_step "Go格式检查" "test -z \$(gofmt -l .)"
echo ""

echo "4. 关键文件检查"
echo "--------------------------------"
test_step "环境变量示例文件" "test -f .env.example"
test_step "迁移SQL文件" "test -f db/migrations/001_update_user_table.sql"
test_step "设计问题文档" "test -f DESIGN_ISSUES.md"
test_step "安全改进文档" "test -f SECURITY_IMPROVEMENTS.md"
test_step "测试验证报告" "test -f TEST_VERIFICATION_REPORT.md"
echo ""

echo "5. 代码质量检查"
echo "--------------------------------"

# 检查是否还有panic（除了测试文件）
if grep -r "panic(" --include="*.go" --exclude="*_test.go" applications/ pkg/ 2>/dev/null | grep -v "// panic"; then
    echo -e "${RED}✗ 发现不当使用panic${NC}"
    ((FAIL++))
else
    echo -e "${GREEN}✓ 无不当panic使用${NC}"
    ((PASS++))
fi

# 检查拼写错误
if grep -r "Aritcle" --include="*.go" applications/ 2>/dev/null; then
    echo -e "${RED}✗ 发现拼写错误 (Aritcle)${NC}"
    ((FAIL++))
else
    echo -e "${GREEN}✓ 无拼写错误${NC}"
    ((PASS++))
fi

# 检查TODO
TODO_COUNT=$(grep -r "todo:" --include="*.go" applications/ | wc -l)
if [ "$TODO_COUNT" -gt 5 ]; then
    echo -e "${YELLOW}⚠ 发现 $TODO_COUNT 个TODO标记${NC}"
else
    echo -e "${GREEN}✓ TODO标记数量合理 ($TODO_COUNT)${NC}"
    ((PASS++))
fi

echo ""
echo "6. 安全检查"
echo "--------------------------------"

# 检查硬编码密钥
if grep -r "const.*passwordEncryptSeed" pkg/encrypt/*.go 2>/dev/null; then
    echo -e "${RED}✗ 发现硬编码密钥${NC}"
    ((FAIL++))
else
    echo -e "${GREEN}✓ 密钥已配置化${NC}"
    ((PASS++))
fi

# 检查bcrypt使用
if grep -r "bcrypt.GenerateFromPassword" pkg/encrypt/*.go 2>/dev/null > /dev/null; then
    echo -e "${GREEN}✓ 已使用bcrypt加密${NC}"
    ((PASS++))
else
    echo -e "${RED}✗ 未使用bcrypt${NC}"
    ((FAIL++))
fi

# 检查类型断言
if grep -r "\.([a-zA-Z.]*)\." applications/article/api/internal/logic/publishlogic.go 2>/dev/null | grep -v "ok :="; then
    echo -e "${YELLOW}⚠ 可能存在不安全的类型断言${NC}"
else
    echo -e "${GREEN}✓ 类型断言已修复${NC}"
    ((PASS++))
fi

echo ""
echo "7. 数据库检查"
echo "--------------------------------"

# 检查用户表定义
if grep -q "KEY \`ix_username\`" db/user.sql; then
    echo -e "${GREEN}✓ 用户名索引已添加${NC}"
    ((PASS++))
else
    echo -e "${RED}✗ 缺少用户名索引${NC}"
    ((FAIL++))
fi

if grep -q "KEY \`ix_phone\`" db/user.sql; then
    echo -e "${GREEN}✓ 手机号索引已添加${NC}"
    ((PASS++))
else
    echo -e "${RED}✗ 缺少手机号索引${NC}"
    ((FAIL++))
fi

if grep -q "varchar(256)" db/user.sql | grep -q "password"; then
    echo -e "${GREEN}✓ 密码字段长度已扩展${NC}"
    ((PASS++))
else
    echo -e "${YELLOW}⚠ 密码字段可能需要扩展${NC}"
fi

echo ""
echo "=================================="
echo "  测试结果汇总"
echo "=================================="
echo -e "通过: ${GREEN}$PASS${NC}"
echo -e "失败: ${RED}$FAIL${NC}"
echo ""

if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}✓✓✓ 所有检查通过！系统已准备好部署。${NC}"
    exit 0
else
    echo -e "${RED}✗✗✗ 有 $FAIL 项检查失败，请修复后重试。${NC}"
    exit 1
fi
