import json

import requests

#url = "http://localhost:8080/"
url = "http://192.168.100.9:8080/testpy"
#url = "http://10.36.172.78:8080/testpy"
data = {
    "message":"Hello from client!",
    "number":42
}

jsonData = json.dumps(data)

response = requests.post(url,data=jsonData)

responseData = json.loads(response.text)
print("Response Message:",responseData['message'])
print("Response Modified:",responseData['modified'])
print("Response Number:",responseData['newNumber'])