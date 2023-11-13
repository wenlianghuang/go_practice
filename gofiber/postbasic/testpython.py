import requests

# 定義要傳遞的 JSON 數據
data = {"name": "John"}

# 發送 POST 請求到 Go Fiber 伺服器
response = requests.post("http://localhost:5050", json=data)

# 打印伺服器的回應
print(response.status_code)
print(response.json())