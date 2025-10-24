#!/bin/bash

# 複雜 Go Breakpoint 練習範例啟動腳本

echo "🎯 複雜 Go Breakpoint 練習範例"
echo "================================"

# 檢查 Go 是否安裝
if ! command -v go &> /dev/null; then
    echo "❌ Go 未安裝，請先安裝 Go"
    exit 1
fi

# 檢查依賴
echo "📦 檢查依賴..."
if ! go mod tidy; then
    echo "❌ 依賴安裝失敗"
    exit 1
fi

# 編譯程序
echo "🔨 編譯程序..."
if ! go build -o complex_breakpoint_example main.go; then
    echo "❌ 編譯失敗"
    exit 1
fi

echo "✅ 編譯成功！"
echo ""
echo "🚀 啟動服務器..."
echo "服務器將在 http://localhost:8080 運行"
echo "按 Ctrl+C 停止服務器"
echo ""

# 啟動服務器
./complex_breakpoint_example
