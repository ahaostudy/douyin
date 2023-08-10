package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"main/service"
	"path/filepath"
)

func PubAction(c *gin.Context) {
	uId := c.GetUint("user_id")
	title := c.PostForm("title")
	file, _ := c.FormFile("data")
	ext := filepath.Ext(file.Filename)
	finalName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	saveFile := filepath.Join("./public/", finalName)
	err := c.SaveUploadedFile(file, saveFile)
	if err != nil {
		c.JSON(200, gin.H{
			"status_code": "1",
			"status_msg":  "Submission failed",
		})
		return
	}
	err = service.SavaFile(uId, finalName, title)
	if err != nil {
		c.JSON(200, gin.H{
			"status_code": "1",
			"status_msg":  "Submission failed",
		})
		return
	}
	c.JSON(200, gin.H{
		"status_code": "1",
		"status_msg":  "Submission success",
	})
}
