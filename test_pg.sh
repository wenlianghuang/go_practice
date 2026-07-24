#!/bin/bash
cd /Users/matthuang/Desktop/go_practice/complex_breakpoint_example

echo "=== 停止旧应用 ==="
pkill -f "go run main.go" 2>/dev/null
sleep 1

echo ""
echo "=== 清除旧数据 ==="
psql breakpoint_db -c "TRUNCATE users CASCADE;" 2>/dev/null

echo ""
echo "=== 启动应用 ==="
go run main.go > /tmp/goapp.log 2>&1 &
APP_PID=$!
sleep 3

echo ""
echo "=== 检查初始用户 ==="
psql breakpoint_db -c "SELECT COUNT(*) as count FROM users;"

echo ""
echo "=== 创建用户 ==="
curl -s -X POST http://localhost:9090/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username":"pgtest","email":"pgtest@test.com","balance":777}'

echo ""
echo ""
echo "=== 等待数据写入 ==="
sleep 2

echo ""
echo "=== 检查数据库 ==="
psql breakpoint_db -c "SELECT id, username, email, balance FROM users ORDER BY id;"

echo ""
echo "=== 停止应用 ==="
kill $APP_PID 2>/dev/null

