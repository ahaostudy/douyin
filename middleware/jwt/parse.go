package jwt

import (
	"github.com/gin-gonic/gin"
	"main/utils"
)

func Parse() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()

		// 获取token
		token := c.Query("token")
		if len(token) == 0 {
			token = c.PostForm("token")
		}
		if len(token) == 0 {
			c.Set("user_id", 0)
			return
		}

		// 解析token
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.Set("user_id", 0)
			return
		}

		c.Set("user_id", claims.ID)
		c.Set("username", claims.Username)
	}
}
