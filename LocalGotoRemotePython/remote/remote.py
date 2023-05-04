import json
from http.server import BaseHTTPRequestHandler, HTTPServer
from io import BytesIO

class RequestHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        # 讀取 POST 請求中的內容
        content_length = int(self.headers['Content-Length'])
        post_data = self.rfile.read(content_length)
        data = json.loads(post_data.decode())

        # 輸出接收到的內容
        print(data)

        # 回傳成功訊息
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()
        response = {
            "status":"success",
            "message":"Successfully received the data.",
            "data":data
        }
        
        response_json = json.dumps(response)
        self.wfile.write(response_json.encode())

def run(server_class=HTTPServer, handler_class=RequestHandler):
    # 設定 HTTP 伺服器
    server_address = ('', 8000)
    httpd = server_class(server_address, handler_class)
    print('Starting server...')

    # 啟動 HTTP 伺服器
    httpd.serve_forever()

if __name__ == '__main__':
    run()
