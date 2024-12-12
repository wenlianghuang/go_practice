## Microservice

在這個有三個範例
- user_service.go => 專門處理用戶資訊微服務
- order_service.go => 專門處理訂單的微服務
- api_gateway.go => 整個微服務架構的入口，負責將用務的請求分發到對應的微服務

#### 定義proto文件
proto文件是用來描述gRPC服務接口和消息結構的。
- user.proto
- order.proto
使用protoc編譯生成Go程式碼
- protobuf(Protocol Buffer)
- google.golang.org/protobuf/cmd/protoc-gen-go@latest
- google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
生成Go程式碼,執行如下
- protoc --go_out=. --go_grpc_out=. user.proto
- protoc --go_out=. --go_grpc_out=. out.proto

運行微服務有三個啟動地點
- go run user_service.go
- go run order_service.go
- go run api_gateway