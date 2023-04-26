package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 建立 Gin 引擎
	r := gin.Default()

	// 設置路由
	r.POST("/upload", func(c *gin.Context) {
		// 從請求中讀取檔案
		file, err := c.FormFile("localfile")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// 將檔案保存到本地
		if err := c.SaveUploadedFile(file, file.Filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// 回應成功訊息
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("檔案 %s 上傳成功", file.Filename),
		})
	})

	// 啟動 HTTP 伺服器
	if err := r.Run(":9090"); err != nil {
		//if err := r.Run("asgard-uat.acer.com/Asgard/index.php:443"); err != nil {
		panic(err)
	}
}
