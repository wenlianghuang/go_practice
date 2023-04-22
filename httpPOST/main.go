package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	fileUrl := "https://proceedings.neurips.cc/paper_files/paper/2017/file/3f5ee243547dee91fbd053c1c4a845aa-Paper.pdf"
	filePath := "D:\\論文\\Attention_is_all_You_Need.pdf"

	client := &http.Client{}
	//req, err := http.NewRequest("POST", fileUrl, strings.NewReader(""))
	req, err := http.NewRequest("GET", fileUrl, nil)
	if err != nil {
		fmt.Println("Error: failed to create the request")
		os.Exit(1)
	}

	req.Header.Set("Accept", "application/pdf")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error: failed to request the file")
		os.Exit(1)

	}

	defer resp.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error: failed to create the file")
		os.Exit(1)
	}

	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Error: failed to save the file")
		os.Exit(1)
	}

	fmt.Println("Successfully downloaded the file!")
}
