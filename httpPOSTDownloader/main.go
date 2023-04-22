package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	//http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/download", downloadHandlerV2)
	fmt.Println("Starting server...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
func downloadHandlerV2(w http.ResponseWriter, r *http.Request) {
	file := "D:\\論文\\Attention_is_all_You_Need.pdf"

	http.ServeFile(w, r, file)
}
func downloadHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/pdf" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filePath := "D:\\go_practice\\httpPOSTDownloadder\\"
	fileName := "D:\\go_practice\\httpPOSTDownloadder\\Attention_is_all_You_Need.pdf"

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error: failed to open the file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Conetnt-Type", "application/octet-stream")

	_, err = io.Copy(w, file)
	if err != nil {
		fmt.Println("Error: failed to write the file to the response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//20230416 https://mileslin.github.io/2020/03/Golang/%E5%BB%BA%E7%AB%8B%E4%B8%8B%E8%BC%89%E6%AA%94%E6%A1%88%E7%9A%84-Http-Response/
	/*
		file := "D:\\go_practice\\httpPOSTDownloader\\download\\Attention_is_all_You_Need.pdf"
		// 讀取檔案
		downloadBytes, err := ioutil.ReadFile(file)

		if err != nil {
			fmt.Println(err)
		}

		// 取得檔案的 MIME type
		mime := http.DetectContentType(downloadBytes)

		fileSize := len(string(downloadBytes))

		w.Header().Set("Content-Type", mime)
		w.Header().Set("Content-Disposition", "attachment; filename="+file)
		w.Header().Set("Content-Length", strconv.Itoa(fileSize))

		http.ServeContent(w, r, file, time.Now(), bytes.NewReader(downloadBytes))
	*/
}

/*
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello!")
}

func main() {
	http.HandleFunc("/hello", helloHandler) // Update this line of code

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
*/
