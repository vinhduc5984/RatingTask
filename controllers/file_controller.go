package controllers

import (
	"fmt"
	"net/http"
	common "todo-list/common"

	"github.com/gin-gonic/gin"
)

func UploadFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy metadata từ form-data
		title := c.PostForm("title")             // Lấy giá trị metadata "title"
		description := c.PostForm("description") // Lấy giá trị metadata "description"

		// Lấy file từ form-data
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File not found"})
			return
		}
		fmt.Println("data file file: ", common.JsonPrettyAny(file))
		// Lưu file vào thư mục "uploads/"
		savePath := "./uploads/" + file.Filename
		fmt.Println("file path data: ", savePath)
		err = c.SaveUploadedFile(file, savePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		// Trả về JSON chứa thông tin file và metadata
		c.JSON(http.StatusOK, gin.H{
			"message":     "File uploaded successfully",
			"file_name":   file.Filename,
			"file_size":   file.Size,
			"title":       title,
			"description": description,
		})
	}
}
