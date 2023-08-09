package jwt

import (
	"github.com/gin-gonic/gin"
	"main/utils"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取token
		token := c.Query("token")
		if len(token) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"status_code": 1,
				"status_msg":  "User authentication failed",
			})
			c.Abort()
		}

		// 解析token
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status_code": 1,
				"status_msg":  "User authentication failed",
			})
			c.Abort()
		}

		c.Set("user_id", claims.ID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
