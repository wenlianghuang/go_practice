#!/bin/bash

# 複雜 Go Breakpoint 練習測試腳本
# 使用方法: ./test_scripts.sh

echo "🎯 複雜 Go Breakpoint 練習測試腳本"
echo "=================================="

# 檢查服務器是否運行
check_server() {
    if curl -s http://localhost:8080/health > /dev/null; then
        echo "✅ 服務器正在運行"
        return 0
    else
        echo "❌ 服務器未運行，請先啟動服務器"
        return 1
    fi
}

# 測試案例 1: 用戶創建
test_create_user() {
    echo ""
    echo "📝 測試案例 1: 用戶創建"
    echo "----------------------"
    
    echo "創建新用戶..."
    curl -X POST http://localhost:8080/api/v1/users \
        -H "Content-Type: application/json" \
        -d '{
            "username": "testuser1",
            "email": "test1@example.com"
        }' | jq .
    
    echo ""
    echo "嘗試創建重複用戶名..."
    curl -X POST http://localhost:8080/api/v1/users \
        -H "Content-Type: application/json" \
        -d '{
            "username": "testuser1",
            "email": "test2@example.com"
        }' | jq .
}

# 測試案例 2: 存款處理
test_deposit() {
    echo ""
    echo "💰 測試案例 2: 存款處理"
    echo "----------------------"
    
    echo "正常存款..."
    curl -X POST http://localhost:8080/api/v1/users/1/deposit \
        -H "Content-Type: application/json" \
        -d '{
            "amount": 100.0,
            "description": "Test deposit"
        }' | jq .
    
    echo ""
    echo "無效金額存款..."
    curl -X POST http://localhost:8080/api/v1/users/1/deposit \
        -H "Content-Type: application/json" \
        -d '{
            "amount": -50.0,
            "description": "Invalid deposit"
        }' | jq .
    
    echo ""
    echo "查看用戶帳戶..."
    curl http://localhost:8080/api/v1/users/1/account | jq .
}

# 測試案例 3: 提款處理
test_withdraw() {
    echo ""
    echo "💸 測試案例 3: 提款處理"
    echo "----------------------"
    
    echo "正常提款..."
    curl -X POST http://localhost:8080/api/v1/users/1/withdraw \
        -H "Content-Type: application/json" \
        -d '{
            "amount": 50.0,
            "description": "Test withdrawal"
        }' | jq .
    
    echo ""
    echo "餘額不足提款..."
    curl -X POST http://localhost:8080/api/v1/users/1/withdraw \
        -H "Content-Type: application/json" \
        -d '{
            "amount": 10000.0,
            "description": "Insufficient funds test"
        }' | jq .
}

# 測試案例 4: 轉帳處理
test_transfer() {
    echo ""
    echo "🔄 測試案例 4: 轉帳處理"
    echo "----------------------"
    
    echo "正常轉帳..."
    curl -X POST http://localhost:8080/api/v1/transfer \
        -H "Content-Type: application/json" \
        -d '{
            "from_user_id": 1,
            "to_user_id": 2,
            "amount": 25.0,
            "description": "Transfer test"
        }' | jq .
    
    echo ""
    echo "餘額不足轉帳..."
    curl -X POST http://localhost:8080/api/v1/transfer \
        -H "Content-Type: application/json" \
        -d '{
            "from_user_id": 2,
            "to_user_id": 1,
            "amount": 1000.0,
            "description": "Insufficient funds transfer"
        }' | jq .
    
    echo ""
    echo "查看兩個用戶的帳戶..."
    echo "用戶 1:"
    curl http://localhost:8080/api/v1/users/1/account | jq .
    echo "用戶 2:"
    curl http://localhost:8080/api/v1/users/2/account | jq .
}

# 測試案例 5: 貸款申請
test_loan() {
    echo ""
    echo "🏦 測試案例 5: 貸款申請"
    echo "----------------------"
    
    echo "申請貸款..."
    curl -X POST http://localhost:8080/api/v1/users/1/apply-loan \
        -H "Content-Type: application/json" \
        -d '{
            "amount": 5000.0,
            "term": 24
        }' | jq .
    
    echo ""
    echo "無效期限貸款申請..."
    curl -X POST http://localhost:8080/api/v1/users/1/apply-loan \
        -H "Content-Type: application/json" \
        -d '{
            "amount": 1000.0,
            "term": 100
        }' | jq .
    
    echo ""
    echo "查看貸款申請狀態..."
    sleep 2  # 等待審核完成
    curl http://localhost:8080/api/v1/users/1/loans | jq .
}

# 測試案例 6: 並發操作
test_concurrent() {
    echo ""
    echo "⚡ 測試案例 6: 並發操作"
    echo "----------------------"
    
    echo "並發存款操作..."
    curl -X POST http://localhost:8080/api/v1/test/concurrent \
        -H "Content-Type: application/json" \
        -d '{
            "user_id": 1,
            "amount": 10.0,
            "operation": "deposit",
            "count": 5
        }' | jq .
    
    echo ""
    echo "並發提款操作..."
    curl -X POST http://localhost:8080/api/v1/test/concurrent \
        -H "Content-Type: application/json" \
        -d '{
            "user_id": 1,
            "amount": 5.0,
            "operation": "withdraw",
            "count": 3
        }' | jq .
}

# 測試案例 7: 錯誤處理
test_error_handling() {
    echo ""
    echo "❌ 測試案例 7: 錯誤處理"
    echo "----------------------"
    
    echo "不存在的用戶..."
    curl http://localhost:8080/api/v1/users/999 | jq .
    
    echo ""
    echo "無效的用戶 ID..."
    curl http://localhost:8080/api/v1/users/invalid | jq .
    
    echo ""
    echo "無效的 JSON..."
    curl -X POST http://localhost:8080/api/v1/users \
        -H "Content-Type: application/json" \
        -d '{"username": "test", "email":}' | jq .
    
    echo ""
    echo "非活躍用戶操作..."
    curl -X POST http://localhost:8080/api/v1/users/3/deposit \
        -H "Content-Type: application/json" \
        -d '{
            "amount": 100.0,
            "description": "Inactive user test"
        }' | jq .
}

# 查看所有數據
view_all_data() {
    echo ""
    echo "📊 查看所有數據"
    echo "----------------"
    
    echo "所有用戶:"
    curl http://localhost:8080/api/v1/users/1 | jq .
    curl http://localhost:8080/api/v1/users/2 | jq .
    curl http://localhost:8080/api/v1/users/3 | jq .
    
    echo ""
    echo "用戶 1 的交易記錄:"
    curl http://localhost:8080/api/v1/users/1/transactions | jq .
    
    echo ""
    echo "用戶 1 的貸款申請:"
    curl http://localhost:8080/api/v1/users/1/loans | jq .
}

# 主菜單
show_menu() {
    echo ""
    echo "🎯 選擇測試案例:"
    echo "1. 用戶創建測試"
    echo "2. 存款處理測試"
    echo "3. 提款處理測試"
    echo "4. 轉帳處理測試"
    echo "5. 貸款申請測試"
    echo "6. 並發操作測試"
    echo "7. 錯誤處理測試"
    echo "8. 查看所有數據"
    echo "9. 運行所有測試"
    echo "0. 退出"
    echo ""
    read -p "請選擇 (0-9): " choice
}

# 主程序
main() {
    if ! check_server; then
        exit 1
    fi
    
    while true; do
        show_menu
        case $choice in
            1) test_create_user ;;
            2) test_deposit ;;
            3) test_withdraw ;;
            4) test_transfer ;;
            5) test_loan ;;
            6) test_concurrent ;;
            7) test_error_handling ;;
            8) view_all_data ;;
            9)
                test_create_user
                test_deposit
                test_withdraw
                test_transfer
                test_loan
                test_concurrent
                test_error_handling
                view_all_data
                ;;
            0) echo "👋 再見！"; exit 0 ;;
            *) echo "❌ 無效選擇，請重新輸入" ;;
        esac
    done
}

# 檢查 jq 是否安裝
if ! command -v jq &> /dev/null; then
    echo "❌ 請先安裝 jq: brew install jq"
    exit 1
fi

# 運行主程序
main
