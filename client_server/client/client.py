import requests
import json 

#url = "http://localhost:8080/"
url = "http://192.168.100.9:80/"
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