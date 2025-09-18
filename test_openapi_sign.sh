#!/bin/bash

# OpenAPI签名验证测试脚本
# 功能: 测试openapi模块的签名验证中间件是否正常工作

BASE_URL="http://127.0.0.1:8001"

echo "=== OpenAPI签名验证中间件测试 ==="
echo ""

# 测试1: 获取公钥接口（无需签名验证）
echo "1. 测试获取公钥接口（无需签名验证）"
echo "请求: GET $BASE_URL/app/openapi/sign/publicKey"
curl -X GET "$BASE_URL/app/openapi/sign/publicKey" \
  -H "Content-Type: application/json" \
  -w "\n状态码: %{http_code}\n" \
  -s
echo ""

# 测试2: 生成签名接口（需要签名验证，但缺少签名头）
echo "2. 测试生成签名接口（缺少签名头，应该返回错误）"
echo "请求: POST $BASE_URL/app/openapi/sign/generate"
curl -X POST "$BASE_URL/app/openapi/sign/generate" \
  -H "Content-Type: application/json" \
  -d '{"data": "test data"}' \
  -w "\n状态码: %{http_code}\n" \
  -s
echo ""

# 测试3: 生成签名接口（提供不完整的签名头）
echo "3. 测试生成签名接口（提供不完整的签名头，应该返回错误）"
echo "请求: POST $BASE_URL/app/openapi/sign/generate"
curl -X POST "$BASE_URL/app/openapi/sign/generate" \
  -H "Content-Type: application/json" \
  -H "X-App-Id: test-app" \
  -H "X-Timestamp: $(date +%s)" \
  -d '{"data": "test data"}' \
  -w "\n状态码: %{http_code}\n" \
  -s
echo ""

# 测试4: 验证签名接口（提供完整但无效的签名头）
echo "4. 测试验证签名接口（提供完整但无效的签名头，应该返回签名验证失败）"
echo "请求: POST $BASE_URL/app/openapi/sign/verify"
curl -X POST "$BASE_URL/app/openapi/sign/verify" \
  -H "Content-Type: application/json" \
  -H "X-App-Id: test-app" \
  -H "X-Timestamp: $(date +%s)" \
  -H "X-Nonce: test-nonce" \
  -H "X-Signature: invalid-signature" \
  -d '{"data": "test data"}' \
  -w "\n状态码: %{http_code}\n" \
  -s
echo ""

echo "=== 测试完成 ==="
echo ""
echo "预期结果："
echo "1. 获取公钥接口应该正常返回（200状态码）"
echo "2. 缺少签名头的请求应该返回400错误"
echo "3. 不完整签名头的请求应该返回400错误"
echo "4. 无效签名的请求应该返回401错误"