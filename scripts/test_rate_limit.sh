#!/bin/bash

# 限速功能测试脚本
# 用于测试公开接口的令牌桶限速功能

BASE_URL="http://localhost:1500"
SERVICES_URL="$BASE_URL/api/open/services"
SUBMIT_URL="$BASE_URL/api/open/submit"

echo "=== 公开接口限速测试 ==="
echo "测试URL: $SERVICES_URL"
echo "限速规则: 每秒10个请求，桶容量20"
echo ""

# 测试1: 正常请求速率
echo "测试1: 正常请求速率 (10个请求，间隔0.1秒)"
for i in {1..10}; do
    response=$(curl -s -w "%{http_code}" -o /dev/null "$SERVICES_URL")
    echo "请求 $i: HTTP $response"
    sleep 0.1
done
echo ""

# 测试2: 突发请求 (应该能处理，因为桶容量为20)
echo "测试2: 突发请求 (20个并发请求)"
for i in {1..20}; do
    curl -s -w "请求 $i: HTTP %{http_code}\n" -o /dev/null "$SERVICES_URL" &
done
wait
echo ""

# 测试3: 超出限制的请求
echo "测试3: 超出限制的请求 (50个并发请求)"
success_count=0
rate_limited_count=0

for i in {1..50}; do
    response=$(curl -s -w "%{http_code}" -o /dev/null "$SERVICES_URL" 2>/dev/null)
    if [ "$response" = "200" ]; then
        ((success_count++))
        echo "请求 $i: HTTP 200 (成功)"
    else
        ((rate_limited_count++))
        echo "请求 $i: HTTP $response (限速)"
    fi
done

echo ""
echo "结果统计:"
echo "成功请求: $success_count"
echo "被限速请求: $rate_limited_count"
echo ""

# 测试4: 提交接口限速测试
echo "=== 提交接口限速测试 ==="
echo "测试URL: $SUBMIT_URL"
echo "限速规则: 每秒2个请求，桶容量5"
echo ""

echo "测试4: 提交接口并发请求 (10个并发请求)"
submit_success=0
submit_rate_limited=0

for i in {1..10}; do
    response=$(curl -s -w "%{http_code}" -o /dev/null \
        -X POST "$SUBMIT_URL" \
        -H "Content-Type: application/json" \
        -d '{"env_id": 1, "value": "test", "key": "test123"}' 2>/dev/null)
    
    if [ "$response" = "200" ]; then
        ((submit_success++))
        echo "提交 $i: HTTP 200 (成功)"
    else
        ((submit_rate_limited++))
        echo "提交 $i: HTTP $response (限速)"
    fi
done

echo ""
echo "提交接口结果统计:"
echo "成功请求: $submit_success"
echo "被限速请求: $submit_rate_limited"
echo ""

# 测试5: 等待令牌恢复
echo "测试5: 等待令牌恢复 (等待3秒后重试)"
sleep 3

echo "恢复后的请求测试:"
for i in {1..5}; do
    response=$(curl -s -w "%{http_code}" -o /dev/null "$SERVICES_URL")
    echo "恢复请求 $i: HTTP $response"
    sleep 0.5
done

echo ""
echo "=== 测试完成 ==="
echo ""
echo "预期结果:"
echo "1. 正常请求应该全部成功 (HTTP 200)"
echo "2. 突发请求大部分成功 (桶容量内)"
echo "3. 超出限制的请求部分被限速"
echo "4. 提交接口更严格，更多请求被限速"
echo "5. 等待后令牌恢复，请求重新成功"
